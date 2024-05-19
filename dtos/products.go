package dtos

import (
	"strings"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
)

type ManagedModel models.Product

type CreateProduct struct {
	Name        string `form:"name" json:"name" xml:"name" binding:"required"`
	Description string `form:"description" json:"description" xml:"description" binding:"required"`
	Price       int    `form:"price" json:"price" xml:"price" binding:"required"`
	Stock       int    `form:"stock" json:"stock" xml:"stock" binding:"required"`
}

func GetProductsDto(products []models.Product) []interface{} {
	var t = make([]interface{}, len(products))
	for i := 0; i < len(products); i++ {
		t[i] = CreateProductDto(products[i])
	}
	return t
}

func CreateProductDto(product models.Product) map[string]interface{} {
	var images = make([]string, len(product.Images))

	replaceAllFlag := -1
	for index, image := range product.Images {
		images[index] = strings.Replace(image.FilePath, "\\", "/", replaceAllFlag)
	}

	var category map[string]interface{}
	if len(product.Categories) > 0 {
		category = map[string]interface{}{
			"id":   product.Categories[0].ID,
			"name": product.Categories[0].Name,
			"slug": product.Categories[0].Slug,
		}
	}

	result := map[string]interface{}{
		"id":          product.ID,
		"name":        product.Name,
		"description": product.Description,
		"slug":        product.Slug,
		"price":       product.Price,
		"stock":       product.Stock,
		"category":    category,
		"image_urls":  images,
		"created_at":  product.CreatedAt.UTC().Format("2006-01-02T15:04:05.999Z"),
	}

	return result
}

func CreateProductCreatedDto(product models.Product) map[string]interface{} {
	return CreateProductDto(product)
}
