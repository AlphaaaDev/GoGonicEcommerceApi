package controllers

// import "C"
import (
	"os"
	"path/filepath"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/dtos"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/middlewares"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/services"

	"net/http"
	"strconv"
)

func RegisterProductRoutes(router *gin.RouterGroup) {
	router.GET("", GetProducts)
	router.GET("/:slug", GetProduct)

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.POST("", CreateProduct)
		router.DELETE("/:id", DeleteProduct)
		router.PATCH("", PatchProduct)
	}
}

func GetProducts(c *gin.Context) {
	products, err := services.FetchProducts()
	if err != nil {
		c.JSON(http.StatusNotFound, "Unable to fetch products")
		return
	}
	c.JSON(http.StatusOK, dtos.GetProductsDto(products))
}

func GetProduct(c *gin.Context) {
	productSlug := c.Param("slug")

	product := services.FetchProduct(&models.Product{Slug: productSlug})
	if product.ID == 0 {
		c.JSON(http.StatusNotFound, "Invalid slug")
		return
	}
	c.JSON(http.StatusOK, dtos.CreateProductDto(product))
}

func CreateProduct(c *gin.Context) {
	// Only admin users can create products
	user := c.Keys["currentUser"].(models.User)
	if user.IsNotAdmin() {
		c.JSON(http.StatusForbidden, "Permission denied, you must be admin")
		return
	}

	var formDto dtos.CreateProduct
	if err := c.ShouldBind(&formDto); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	name := formDto.Name
	description := formDto.Description

	price := formDto.Price
	stock, err := strconv.ParseInt(c.PostForm("stock"), 10, 32)
	categoryId64, _ := strconv.ParseUint(c.PostForm("category"), 10, 32)
	categoryId := uint(categoryId64)
	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	files := form.File["images[]"]
	var productImages = make([]models.FileUpload, len(files))
	for index, file := range files {
		fileName := randomString(16) + ".png"

		dirPath := filepath.Join(".", "static", "images", "products")
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
		productImages[index] = models.FileUpload{Filename: fileName, OriginalName: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
	}

	var category []models.Category
	if categoryId == 0 {
		category = []models.Category{}
	} else {
		category = []models.Category{services.FetchCategoryById(categoryId)}
	}

	product := models.Product{
		Name:        name,
		Description: description,
		Categories:  category,
		Price:       (int)(price),
		Stock:       (int)(stock),
		Images:      productImages,
	}

	if err := services.CreateOne(&product); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, dtos.CreateProductCreatedDto(product))

}

func DeleteProduct(c *gin.Context) {
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)

	err := services.DeleteProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, "Invalid id")
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": "Delete success"})
}

func PatchProduct(c *gin.Context) {
	type PatchProductRequest struct {
		Id          *uint   `form:"id"`
		Name        *string `form:"name"`
		Description *string `form:"description"`
		Price       *int    `form:"price"`
		Stock       *int    `form:"stock"`
		Category    *uint   `form:"category"`
	}

	var patchRequest PatchProductRequest
	if err := c.ShouldBind(&patchRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	product := services.FetchProductById(*patchRequest.Id)

	if patchRequest.Name != nil {
		product.Name = *patchRequest.Name
	}

	if patchRequest.Description != nil {
		product.Description = *patchRequest.Description
	}

	if patchRequest.Price != nil {
		product.Price = *patchRequest.Price
	}

	if patchRequest.Stock != nil {
		product.Stock = *patchRequest.Stock
	}

	if patchRequest.Category != nil {
		if *patchRequest.Category == 0 {
			product.Categories = []models.Category{}
		} else {
			product.Categories = []models.Category{services.FetchCategoryById(*patchRequest.Category)}
		}
	}

	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	files := form.File["images[]"]
	var productImages = make([]models.FileUpload, len(files))
	if len(files) > 0 {
		for index, file := range files {
			fileName := randomString(16) + ".png"

			dirPath := filepath.Join(".", "static", "images", "products")
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
			productImages[index] = models.FileUpload{Filename: fileName, OriginalName: file.Filename, FilePath: string(filepath.Separator) + filePath, FileSize: fileSize}
		}
	}

	product.Images = productImages

	product.Slug = slug.Make(product.Name)

	if err := services.UpdateProduct(product, product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"product": "failed to update product"})
		return
	}

	c.JSON(http.StatusOK, patchRequest)
}
