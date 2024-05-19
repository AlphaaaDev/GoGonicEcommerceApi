package controllers

import (
	"net/http"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/dtos"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/services"
	"github.com/gin-gonic/gin"
)

func RegisterPageRoutes(router *gin.RouterGroup) {
	router.GET("", Home)
	router.GET("/home", Home)

}

func Home(c *gin.Context) {
	categories, err := services.FetchCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, "Unable to fetch categories")
		return
	}
	c.JSON(http.StatusOK, dtos.GetCategoriesDto(categories))
}
