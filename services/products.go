package services

import (
	"github.com/AlphaaaDev/GoGonicEcommerceApi/infrastructure"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
)

func FetchProducts() ([]models.Product, error) {
	database := infrastructure.GetDb()
	var products []models.Product
	err := database.Preload("Images").Preload("Categories").Find(&products).Error
	return products, err
}

func FetchProductById(id uint) models.Product {
	db := infrastructure.GetDb()
	var product models.Product
	db.Model(models.Order{}).Preload("Categories").Preload("Images").First(&product, id)
	return product
}

func FetchProductsByCategory(categorySlug string) ([]models.Product, error) {
	db := infrastructure.GetDb()
	category := FetchCategoryBySlug(categorySlug)

	var products []models.Product
	err := db.Preload("Images").Preload("Categories").
		Joins("left join products_categories on products.id = products_categories.product_id").
		Where("products_categories.category_id = ?", category.ID).
		Find(&products).Error
	return products, err
}

func FetchProduct(condition interface{}) models.Product {
	database := infrastructure.GetDb()
	var product models.Product

	query := database.Where(condition).Preload("Categories").Preload("Images")
	query.First(&product)

	return product
}

func FetchProductId(slug string) (uint, error) {
	productId := -1
	database := infrastructure.GetDb()
	err := database.Model(&models.Product{}).Where(&models.Product{Slug: slug}).Select("id").Row().Scan(&productId)
	return uint(productId), err
}

func DeleteProduct(id uint) error {
	db := infrastructure.GetDb()
	err := db.Unscoped().Delete(models.Product{}, id).Error
	return err
}

func FetchProductsIdNameAndPrice(productIds []uint) (products []models.Product, err error) {
	database := infrastructure.GetDb()
	err = database.Select([]string{"id", "name", "slug", "price"}).Find(&products, productIds).Error
	return products, err
}

func UpdateProduct(product models.Product, data interface{}) error {
	db := infrastructure.GetDb()
	if len(product.Categories) == 0 {
		db.Model(product).Association("Categories").Delete(FetchProductById(product.ID).Categories[0])
	}

	db.Model(product).Association("Images").Clear()

	err := db.Model(product).Update(data).Error
	return err
}
