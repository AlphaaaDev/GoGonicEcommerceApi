package dtos

import (
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
)

type CreateAddress struct {
	FirstName     string `form:"first_name" json:"first_name" xml:"first_name"`
	LastName      string `form:"last_name" json:"last_name" xml:"last_name"`
	Country       string `form:"country" json:"country" xml:"country" binding:"required"`
	City          string `form:"city" json:"city" xml:"city" binding:"required"`
	StreetAddress string `form:"address" json:"address" xml:"address" binding:"required"`
	ZipCode       string `form:"zip_code" json:"zip_code" xml:"zip_code" binding:"required"`
}

func GetAddressesDto(addresses []models.Address) []interface{} {
	var t = make([]interface{}, len(addresses))
	for index, address := range addresses {
		t[index] = GetAddressDto(address, false)
	}
	return t
}

func GetAddressDto(address models.Address, includeUser bool) map[string]interface{} {
	dto := map[string]interface{}{
		"id":             address.ID,
		"first_name":     address.FirstName,
		"last_name":      address.LastName,
		"zip_code":       address.ZipCode,
		"country":        address.Country,
		"city":           address.City,
		"street_address": address.StreetAddress,
	}

	if includeUser {
		dto["user"] = map[string]interface{}{
			"id":       address.UserId,
			"username": address.User.Username,
		}
	}
	return dto
}

func GetAddressCreatedDto(address models.Address, includeUser bool) map[string]interface{} {
	return GetAddressDto(address, includeUser)
}
