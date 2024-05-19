package controllers

import (
	"net/http"
	"strconv"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/dtos"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/middlewares"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/services"
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(router *gin.RouterGroup) {
	router.POST("", CreateOrder)
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("", GetOrders)
		router.GET("/:id", GetOrder)
	}
}

func GetOrders(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)
	orders, err := services.FetchOrders(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Unable to fetch orders")
		return
	}
	c.JSON(http.StatusOK, dtos.GetOrdersDto(orders))
}

func GetOrder(c *gin.Context) {
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)
	user := c.MustGet("currentUser").(models.User)
	order, err := services.FetchOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if order.UserId == user.ID || user.IsAdmin() {
		c.JSON(http.StatusOK, dtos.CreateOrderDetailsDto(order))
	} else {
		c.JSON(http.StatusForbidden, "Permission denied, you can not view this order")
		return
	}
}

func CreateOrder(c *gin.Context) {
	var orderRequest dtos.CreateOrderRequestDto
	if err := c.ShouldBind(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	userObj, userLoggedIn := c.Get("currentUser")
	var user models.User
	if userLoggedIn {
		user = (userObj).(models.User)
	}

	var address models.Address
	// Reuse address can only be done by authenticated users
	if orderRequest.AddressId != 0 && userLoggedIn {
		address = services.FetchAddress(orderRequest.AddressId)
		if address.UserId != user.ID {
			c.JSON(http.StatusForbidden, "Permission denied")
			return
		}
	} else if orderRequest.AddressId == 0 {
		address = models.Address{
			FirstName:     orderRequest.FirstName,
			LastName:      orderRequest.LastName,
			City:          orderRequest.City,
			Country:       orderRequest.Country,
			StreetAddress: orderRequest.StreetAddress,
			ZipCode:       orderRequest.ZipCode,
		}
		if userLoggedIn {
			address.UserId = user.ID
		}
		err := services.CreateOne(&address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

	} else {
		c.JSON(http.StatusForbidden, "Operation not supported, what are you trying to do?")
		return
	}

	order := models.Order{
		TrackingNumber: randomString(16),
		OrderStatus:    0,
		Address:        address,
		AddressId:      address.ID,
	}

	if userLoggedIn {
		order.UserId = user.ID
		order.User = user
	}

	var productIds = make([]uint, len(orderRequest.CartItems))
	for i := 0; i < len(orderRequest.CartItems); i++ {
		productIds[i] = orderRequest.CartItems[i].Id
	}

	products, err := services.FetchProductsIdNameAndPrice(productIds)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if len(products) != len(orderRequest.CartItems) {
		c.JSON(http.StatusUnprocessableEntity, "Make sure all products are still available")
		return
	}
	orderItems := make([]models.OrderItem, len(products))

	for i := 0; i < len(products); i++ {
		// I am assuming product ids returned are in the same order as the cart_items, TODO: implement a more robust code to ensure
		orderItems[i] = models.OrderItem{
			ProductId:   products[i].ID,
			ProductName: products[i].Name,
			Slug:        products[i].Slug,
			Price:       products[i].Price,
			Quantity:    orderRequest.CartItems[i].Quantity,
		}
	}

	order.OrderItems = orderItems
	err = services.CreateOne(&order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, dtos.CreateOrderCreatedDto(order))

}
