package controllers

import (
	"net/http"
	"tes-mnc/initializers"
	"tes-mnc/models"

	"github.com/gin-gonic/gin"
)

func RegisterSeller(c *gin.Context) {
	var body struct {
		FullName    string
		PhoneNumber string
		Email       string
		Address     string
		Store       models.Store
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	seller := models.Seller{
		FullName:    body.FullName,
		PhoneNumber: body.PhoneNumber,
		Address:     body.Address,
		Email:       body.Email,
		Store: models.Store{
			NoSiup:      body.Store.NoSiup,
			Name:        body.Store.Name,
			Address:     body.Store.Address,
			PhoneNumber: body.Store.PhoneNumber,
		},
	}

	result := initializers.DB.Create(&seller)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Failed to create seller",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully created seller",
	})
}

func GetAllSeller(c *gin.Context) {
	var sellers []models.Seller

	result := initializers.DB.Preload("Store").Find(&sellers)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch seller list",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"error":   "No seller registered",
			"message": "Successfully fetch",
		})
		return
	}

	var simplifiedSellers []map[string]interface{}

	for _, seller := range sellers {
		simplifiedSeller := map[string]interface{}{
			"id":                 seller.ID,
			"full_name":          seller.FullName,
			"phone_number":       seller.PhoneNumber,
			"address":            seller.Address,
			"email":              seller.Email,
			"store_id":           seller.Store.ID,
			"store_no_siup":      seller.Store.NoSiup,
			"store_name":         seller.Store.Name,
			"store_address":      seller.Store.Address,
			"store_phone_number": seller.Store.PhoneNumber,
		}
		simplifiedSellers = append(simplifiedSellers, simplifiedSeller)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    simplifiedSellers,
		"message": "Successfully fetch all seller data",
	})
}
