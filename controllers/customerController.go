package controllers

import (
	"net/http"
	"tes-mnc/initializers"
	"tes-mnc/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterCustomer(c *gin.Context) {
	var body struct {
		FullName    string
		Address     string
		PhoneNumber string
		User        models.User
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.User.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
	}

	customer := models.Customer{
		FullName:    body.FullName,
		Address:     body.Address,
		PhoneNumber: body.PhoneNumber,
		User:        models.User{Email: body.User.Email, Password: string(hash)},
	}

	result := initializers.DB.Create(&customer)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  "Failed to create customer",
		})

		return
	}

	fullNameCustomer := customer.FullName
	emailCustomer := customer.User.Email

	c.JSON(http.StatusOK, gin.H{
		"full_name": fullNameCustomer,
		"email":     emailCustomer,
		"status":    http.StatusOK,
		"message":   "Successfully created customer",
	})
}

func GetAllCustomer(c *gin.Context) {
	var customers []models.Customer

	result := initializers.DB.Preload("User").Find(&customers)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to fetch customer list",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"error":   "No customer registered",
			"message": "Successfully fetch",
		})
		return
	}

	// Buat slice baru untuk menyimpan data yang diperlukan
	var simplifiedCustomers []gin.H

	for _, customer := range customers {
		simplifiedCustomer := gin.H{
			"id":           customer.ID,
			"full_name":    customer.FullName,
			"address":      customer.Address,
			"phone_number": customer.PhoneNumber,
			"user_id":      customer.User.ID,
			"user_email":   customer.User.Email,
		}

		simplifiedCustomers = append(simplifiedCustomers, simplifiedCustomer)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"data":    simplifiedCustomers,
		"message": "Successfully fetch all customer data",
	})
}

func UpdateCustomer(c *gin.Context) {
	var body struct {
		ID          uint
		FullName    string
		Address     string
		PhoneNumber string
		User        models.User
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Unauthorized",
			"message": "User must be log in to update customer",
		})
		return
	}

	authenticatedUserID := user.(models.User).ID

	var existingCustomer models.Customer
	result := initializers.DB.Preload("User").First(&existingCustomer, body.ID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Customer not found",
			"message": "Failed to update customer",
		})
		return
	}

	if existingCustomer.User.ID != authenticatedUserID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Forbidden",
			"message": "You don't have permission to update customer",
		})
		return
	}

	if existingCustomer.User.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "User not found",
			"message": "Failed to update customer",
		})
		return
	}

	existingCustomer.User.Email = body.User.Email

	if body.User.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.User.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		existingCustomer.User.Password = string(hash)
	}

	result = initializers.DB.Save(&existingCustomer.User)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	existingCustomer.FullName = body.FullName
	existingCustomer.Address = body.Address
	existingCustomer.PhoneNumber = body.PhoneNumber

	result = initializers.DB.Save(&existingCustomer)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to update customer",
		})
		return
	}

	var updatedCustomer models.Customer
	result = initializers.DB.Preload("User").First(&updatedCustomer, existingCustomer.ID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch updated customer data",
			"message": "Successfully updated customer",
			"status":  http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Successfully updated customer",
		"data":    updatedCustomer,
	})
}
