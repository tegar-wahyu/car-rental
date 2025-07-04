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

	// Car routes
	cars := r.Group("/cars")
	{
		cars.GET("", handlers.GetCars)
		cars.GET("/:id", handlers.GetCar)
		cars.POST("", handlers.CreateCar)
		cars.PUT("/:id", handlers.UpdateCar)
		cars.DELETE("/:id", handlers.DeleteCar)
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Car Rental API",
		})
	})
}
