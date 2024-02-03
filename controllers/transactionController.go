package controllers

import (
	"net/http"
	"tes-mnc/initializers"
	"tes-mnc/models"

	"github.com/gin-gonic/gin"
)

func CreateTransaction(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "Customer must be log in to create a transaction",
		})
		return
	}

	customerID := user.(models.User).ID

	var body struct {
		TransactionAmount      uint64
		TransactionDescription string
		SellerID               uint32
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	history := models.Transaction{
		CustomerID:             uint32(customerID),
		TransactionAmount:      body.TransactionAmount,
		TransactionDescription: body.TransactionDescription,
		SellerID:               body.SellerID,
	}

	result := initializers.DB.Create(&history)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create transaction",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully created transaction",
	})

}

func GetTransaction(c *gin.Context) {
	var transactions []models.Transaction

	result := initializers.DB.Preload("Customer.User").Preload("Seller.Store").Find(&transactions)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch transaction list",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"error":   "No transaction found",
			"message": "Successfully fetch",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    transactions,
		"status":  http.StatusOK,
		"message": "Successfully fetch all transaction data",
	})
}