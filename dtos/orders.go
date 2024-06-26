package dtos

import (
	"strings"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/services"
)

type CreateOrderRequestDto struct {
	FirstName     string `form:"first_name" json:"first_name" xml:"first_name"`
	LastName      string `form:"last_name" json:"last_name" xml:"last_name"`
	Country       string `form:"country" json:"country" xml:"country"`
	City          string `form:"city" json:"city" xml:"city"`
	StreetAddress string `form:"street_address" json:"street_address" xml:"street_address" `
	ZipCode       string `form:"zip_code" json:"zip_code" xml:"zip_code" `
	AddressId     uint   `form:"address_id" json:"address_id" xml:"address_id" `
	CartItems     []struct {
		Id       uint `form:"id" json:"id" binding:"required"`
		Quantity int  `form:"quantity" json:"quantity" binding:"required"`
	} `json:"cart_items"`
}

func GetOrdersDto(orders []models.Order) []interface{} {
	var t = make([]interface{}, len(orders))
	for i := 0; i < len(orders); i++ {
		t[i] = CreateOrderDto(orders[i])
	}
	return t
}

func CreateOrderDto(order models.Order, includes ...bool) map[string]interface{} {

	includeAddress, includeOrderItems, includeUser := getIncludeFlags(includes...)

	result := map[string]interface{}{
		"id":              order.ID,
		"tracking_number": order.TrackingNumber,
		"order_status":    order.GetOrderStatusAsString(),
		"created_at":      order.CreatedAt.UTC().Format("2006-01-02"),
	}

	if includeAddress {
		result["address"] = map[string]interface{}{
			"first_name":     order.Address.FirstName,
			"last_name":      order.Address.LastName,
			"street_address": order.Address.StreetAddress,
			"city":           order.Address.City,
			"country":        order.Address.Country,
			"zip_code":       order.Address.ZipCode,
		}
	}

	if includeOrderItems {
		orderItems := make([]map[string]interface{}, len(order.OrderItems))
		for i := 0; i < len(order.OrderItems); i++ {
			oi := order.OrderItems[i]
			p := services.FetchProductById(oi.ProductId)
			var images = make([]string, len(p.Images))

			replaceAllFlag := -1
			for index, image := range p.Images {
				images[index] = strings.Replace(image.FilePath, "\\", "/", replaceAllFlag)
			}
			orderItems[i] = map[string]interface{}{
				"name":       oi.ProductName,
				"slug":       oi.Slug,
				"price":      oi.Price,
				"quantity":   oi.Quantity,
				"image_urls": images,
			}
		}
		result["order_items"] = orderItems
	} else {
		result["order_items_count"] = order.OrderItemsCount
	}

	if includeUser {
		result["user"] = map[string]interface{}{
			"id":       order.UserId,
			"username": order.User.Username,
		}
	}

	return result
}

func CreateOrderDetailsDto(order models.Order) map[string]interface{} {
	// includeUser -> false
	// includeOrderItems -> true
	// includeUser -> false
	return CreateOrderDto(order, true, true, false)
}

func getIncludeFlags(includes ...bool) (includeAddress, includeOrderItems, includeUser bool) {

	if len(includes) > 0 {
		includeAddress = includes[0]
	}

	if len(includes) > 1 {
		includeOrderItems = includes[1]
	}

	if len(includes) > 2 {
		includeUser = includes[2]
	}
	return
}

func CreateOrderCreatedDto(order models.Order) map[string]interface{} {
	return CreateOrderDetailsDto(order)
}
