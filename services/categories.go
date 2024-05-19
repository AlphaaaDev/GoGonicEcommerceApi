package services

import (
	"github.com/AlphaaaDev/GoGonicEcommerceApi/infrastructure"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
)

func FetchCategories() ([]models.Category, error) {
	database := infrastructure.GetDb()
	var categories []models.Category
	err := database.Preload("Images").Find(&categories).Error
	return categories, err
}

func FetchCategoryById(id uint) models.Category {
	db := infrastructure.GetDb()
	var category models.Category
	db.Model(models.Order{}).Preload("Images").First(&category, id)
	return category
}

func FetchCategoryBySlug(slug string) models.Category {
	db := infrastructure.GetDb()
	var category models.Category
	db.Model(models.Order{}).Preload("Images").First(&category, models.Category{Slug: slug})
	return category
}

func DeleteCategory(id uint) error {
	db := infrastructure.GetDb()
	err := db.Unscoped().Delete(models.Category{}, id).Error
	return err
}

func UpdateCategory(category models.Category, data interface{}) error {
	db := infrastructure.GetDb()
	err := db.Model(category).Update(data).Error
	return err
}
