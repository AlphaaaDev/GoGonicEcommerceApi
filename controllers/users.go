package controllers

import (
	"github.com/AlphaaaDev/GoGonicEcommerceApi/dtos"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/middlewares"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/services"
	"github.com/gin-gonic/gin"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"

	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("", UsersRegistration)
	router.POST("/login", UsersLogin)

	router.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		router.GET("", GetUser)
		router.PATCH("", PatchUser)
	}
}

func GetUser(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)

	c.JSON(http.StatusOK, dtos.CreateUserDto(user))
}

func PatchUser(c *gin.Context) {
	user := c.Keys["currentUser"].(models.User)

	type PatchUserRequest struct {
		Username             *string `json:"username"`
		FirstName            *string `json:"first_name"`
		LastName             *string `json:"last_name"`
		Email                *string `json:"email"`
		Password             *string `json:"password"`
		PasswordConfirmation *string `json:"password_confirmation"`
	}

	var patchRequest PatchUserRequest
	if err := c.ShouldBindJSON(&patchRequest); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if patchRequest.Username != nil {
		user.Username = *patchRequest.Username
	}

	if patchRequest.FirstName != nil {
		user.FirstName = *patchRequest.FirstName
	}

	if patchRequest.LastName != nil {
		user.LastName = *patchRequest.LastName
	}

	if patchRequest.Email != nil {
		user.Email = *patchRequest.Email
	}

	if patchRequest.Password != nil && patchRequest.PasswordConfirmation != nil {
		if *patchRequest.Password != *patchRequest.PasswordConfirmation {
			c.JSON(http.StatusUnprocessableEntity, "Passwords don't match")
			return
		}

		password, _ := bcrypt.GenerateFromPassword([]byte(*patchRequest.Password), bcrypt.DefaultCost)
		user.Password = string(password)
	}

	if err := services.UpdateUser(user, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, patchRequest)
}

func UsersRegistration(c *gin.Context) {
	var json dtos.RegisterRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if json.Password != json.PasswordConfirmation {
		c.JSON(http.StatusUnprocessableEntity, "Passwords don't match")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err := services.CreateOne(&models.User{
		Username:  json.Username,
		Password:  string(password),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Email:     json.Email,
	}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success":       true,
		"full_messages": []string{"User created successfully"}})
}

func UsersLogin(c *gin.Context) {

	var json dtos.LoginRequestDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	var err error
	_, emailError := mail.ParseAddress(json.Username)
	if emailError == nil {
		user, err = services.FindOneUser(&models.User{Email: json.Username})
	} else {
		user, err = services.FindOneUser(&models.User{Username: json.Username})
	}

	if err != nil {
		c.JSON(http.StatusForbidden, err.Error())
		return
	}

	if user.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, dtos.CreateLoginSuccessful(&user))

}
