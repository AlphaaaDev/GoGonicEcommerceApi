package services

import (
	"github.com/AlphaaaDev/GoGonicEcommerceApi/infrastructure"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
)

func FetchOrders(userId uint) ([]models.Order, error) {
	database := infrastructure.GetDb()
	var orders []models.Order
	err := database.Where(&models.Order{UserId: userId}).Find(&orders).Error
	return orders, err
}

func FetchOrder(orderId uint) (order models.Order, err error) {
	database := infrastructure.GetDb()
	err = database.Model(models.Order{}).Preload("OrderItems").First(&order, orderId).Error
	var address models.Address
	database.Model(&order).Related(&address)
	order.Address = address
	return order, err
}
