package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/dtos"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/infrastructure"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/middlewares"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/services"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func RegisterCategoryRoutes(router *gin.RouterGroup) {
	router.GET("", GetCategories)
	router.GET("/:slug", GetCategory)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("", CreateCategory)
		router.DELETE("/:id", DeleteCategory)
		router.PATCH("", PatchCategory)
	}
}

func GetCategories(c *gin.Context) {
	categories, err := services.FetchCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, "Unable to fetch categories")
		return
	}
	c.JSON(http.StatusOK, dtos.GetCategoriesDto(categories))
}

func GetCategory(c *gin.Context) {
	products, err := services.FetchProductsByCategory(c.Param("slug"))
	if err != nil {
		c.JSON(http.StatusNotFound, "Unable to fetch products")
		return
	}
	c.JSON(http.StatusOK, dtos.GetProductsDto(products))
}

func CreateCategory(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, "Permission denied, you must be admin")
		return
	}
	name := c.PostForm("name")
	description := c.PostForm("description")

	form, err := c.MultipartForm()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	files := form.File["images[]"]
	var categoryImages = make([]models.FileUpload, len(files))
	for index, file := range files {
		fileName := randomString(16) + ".png"

		dirPath := filepath.Join(".", "static", "images", "categories")
		filePath := filepath.Join(dirPath, fileName)
		// Create directory if does not exist
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, os.ModeDir)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}
		// Create file that will hold the image
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Open the temporary file that contains the uploaded image
		inputFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusOK, err.Error())
		}
		defer inputFile.Close()

		// Copy the temporary image to the permanent location outputFile
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		fileSize := (uint)(file.Size)
		categoryImages[index] = models.FileUpload{Filename: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	database := infrastructure.GetDb()
	category := models.Category{Name: name, Description: description, Images: categoryImages}

	// TODO: Why it is performing a SELECT SQL Query per image?
	// Even worse, it is selecting category_id, why??
	// SELECT "tag_id", "product_id" FROM "file_uploads"  WHERE (id = insertedFileUploadId)
	err = database.Create(&category).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, dtos.CreateCategoryCreatedDto(category))
}

func DeleteCategory(c *gin.Context) {
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	err := services.DeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, "Invalid id")
		return
	}
	c.JSON(http.StatusOK, gin.H{"categories": "Delete success"})
}

func PatchCategory(c *gin.Context) {
	type PatchCategoryRequest struct {
		Id          *uint   `form:"id"`
		Name        *string `form:"name"`
		Description *string `form:"description"`
	}

	var patchRequest PatchCategoryRequest
	if err := c.ShouldBind(&patchRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	category := services.FetchCategoryById(*patchRequest.Id)

	if patchRequest.Name != nil {
		category.Name = *patchRequest.Name
	}

	if patchRequest.Description != nil {
		category.Description = *patchRequest.Description
	}

	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	files := form.File["images[]"]
	var categoryImages = make([]models.FileUpload, len(files))
	if len(files) > 0 {
		for index, file := range files {
			fileName := randomString(16) + ".png"

			dirPath := filepath.Join(".", "static", "images", "categories")
			filePath := filepath.Join(dirPath, fileName)
			if _, err = os.Stat(dirPath); os.IsNotExist(err) {
				err = os.MkdirAll(dirPath, os.ModeDir)
				if err != nil {
					c.JSON(http.StatusInternalServerError, err.Error())
					return
				}
			}
			if err := c.SaveUploadedFile(file, filePath); err != nil {
				c.JSON(http.StatusBadRequest, err.Error())
				return
			}
			fileSize := (uint)(file.Size)
			categoryImages[index] = models.FileUpload{Filename: fileName, OriginalName: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
		}
	}

	category.Images = categoryImages

	category.Slug = slug.Make(category.Name)

	if err := services.UpdateCategory(category, category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"categories": "failed to update category"})
		return
	}

	c.JSON(http.StatusOK, patchRequest)
}
