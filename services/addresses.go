package services

import (
	"github.com/AlphaaaDev/GoGonicEcommerceApi/infrastructure"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
)

func FetchAddresses(userId uint) ([]models.Address, error) {
	db := infrastructure.GetDb()
	var addresses []models.Address
	err := db.Where(&models.Address{UserId: userId}).Find(&addresses).Error
	return addresses, err
}

func FetchAddress(addressId uint) (address models.Address) {
	database := infrastructure.GetDb()
	database.First(&address, addressId)
	return address
}

func DeleteAddress(id uint) error {
	db := infrastructure.GetDb()
	err := db.Unscoped().Delete(models.Address{}, id).Error
	return err
}
