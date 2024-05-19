package controllers

import (
	"github.com/AlphaaaDev/GoGonicEcommerceApi/dtos"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/middlewares"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/services"
	"github.com/gin-gonic/gin"

	"net/http"
	"strconv"
)

func RegisterAddressesRoutes(router *gin.RouterGroup) {
	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("", GetAddresses)
		router.POST("", CreateAddress)
		router.DELETE("/:id", DeleteAddress)
	}
}

func GetAddresses(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)
	addresses, err := services.FetchAddresses(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, "Unable to fetch addresses")
		return
	}
	c.JSON(http.StatusOK, dtos.GetAddressesDto(addresses))
}

func CreateAddress(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)

	var json dtos.CreateAddress
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	firstName := json.FirstName
	lastName := json.LastName
	if firstName == "" {
		firstName = user.FirstName
	}
	if lastName == "" {
		lastName = user.LastName
	}
	address := models.Address{
		FirstName:     firstName,
		LastName:      lastName,
		Country:       json.Country,
		City:          json.City,
		StreetAddress: json.StreetAddress,
		ZipCode:       json.ZipCode,
		User:          user,
		UserId:        user.ID,
	}

	if err := services.SaveOne(&address); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, dtos.GetAddressCreatedDto(address, false))
}

func DeleteAddress(c *gin.Context) {
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	id := uint(id64)

	err := services.DeleteAddress(id)
	if err != nil {
		c.JSON(http.StatusNotFound, "Invalid id")
		return
	}
	c.JSON(http.StatusOK, gin.H{"addresses": "Delete success"})
}
