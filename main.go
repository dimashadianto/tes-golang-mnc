package main

import (
	"tes-mnc/controllers"
	"tes-mnc/initializers"
	"tes-mnc/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	r.POST("/api/auth/signup", controllers.Signup)
	r.POST("/api/auth/login", controllers.Login)
	r.GET("/api/auth/validate", middleware.RequireAuth, controllers.Validate)
	r.GET("/api/auth/logout", middleware.RequireAuth, controllers.Logout)

	r.POST("/api/auth/customer/register", controllers.RegisterCustomer)
	r.GET("/api/auth/customer/all", controllers.GetAllCustomer)
	r.PUT("/api/auth/customer/update", controllers.UpdateCustomer)

	r.POST("/api/auth/store/register", controllers.RegisterStore)
	r.GET("/api/auth/store/all", controllers.GetAllStore)

	r.POST("/api/auth/seller/register", controllers.RegisterSeller)
	r.GET("/api/auth/seller/all", controllers.GetAllSeller)

	r.POST("/api/auth/transaction/create", controllers.CreateTransaction)
	r.GET("/api/auth/transaction/all", controllers.GetTransaction)

	r.Run()
}
