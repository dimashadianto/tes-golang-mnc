package initializers

import "tes-mnc/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Customer{})
	DB.AutoMigrate(&models.Store{})
	DB.AutoMigrate(&models.Seller{})
	DB.AutoMigrate(&models.Transaction{})
}