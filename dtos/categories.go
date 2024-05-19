package dtos

import (
	"strings"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
)

func GetCategoriesDto(categories []models.Category) []interface{} {
	var t = make([]interface{}, len(categories))
	for i := 0; i < len(categories); i++ {
		t[i] = CreateCategoryDto(categories[i])
	}
	return t
}

func CreateCategoryDto(category models.Category) map[string]interface{} {
	var imageUrls = make([]string, len(category.Images))
	replaceAllFlag := -1
	for i := 0; i < len(category.Images); i++ {
		imageUrls[i] = strings.Replace(category.Images[i].FilePath, "\\", "/", replaceAllFlag)
	}
	return map[string]interface{}{
		"id":          category.ID,
		"name":        category.Name,
		"description": category.Description,
		"image_urls":  imageUrls,
		"slug":        category.Slug,
	}
}

func CreateCategoryCreatedDto(category models.Category) map[string]interface{} {
	return CreateCategoryDto(category)
}
