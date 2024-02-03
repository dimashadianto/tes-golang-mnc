package controllers

import (
	"net/http"
	"tes-mnc/initializers"
	"tes-mnc/models"

	"github.com/gin-gonic/gin"
)

func RegisterStore(c *gin.Context) {
	var body struct {
		NoSiup      string
		Name        string
		Address     string
		PhoneNumber string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	store := models.Store{
		NoSiup:      body.NoSiup,
		Name:        body.Name,
		Address:     body.Address,
		PhoneNumber: body.PhoneNumber,
	}

	result := initializers.DB.Create(&store)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create store",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Succesfully created store",
	})

}

func GetAllStore(c *gin.Context) {
	var stores []models.Store

	result := initializers.DB.Find(&stores)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch store list",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No store registered",
		})
		return
	}

	var simplifiedStores []map[string]interface{}

	for _, store := range stores {
		simplifiedStore := map[string]interface{}{
			"id":           store.ID,
			"store_name":   store.Name,
			"no_siup":      store.NoSiup,
			"phone_number": store.PhoneNumber,
			"address":      store.Address,
		}
		simplifiedStores = append(simplifiedStores, simplifiedStore)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully fetch all store data",
		"data":    simplifiedStores,
	})
}