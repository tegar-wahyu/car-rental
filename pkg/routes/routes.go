package routes

import (
	"car-rental/pkg/handlers"
	
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Customer routes
	customers := r.Group("/customers")
	{
		customers.GET("", handlers.GetCustomers)
		customers.GET("/:id", handlers.GetCustomer)
		customers.POST("", handlers.CreateCustomer)
		customers.PUT("/:id", handlers.UpdateCustomer)
		customers.DELETE("/:id", handlers.DeleteCustomer)
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Car Rental API",
		})
	})
}
