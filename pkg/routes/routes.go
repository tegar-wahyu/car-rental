package routes

import (
	"car-rental/pkg/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Create API v1 group
	v1 := r.Group("/api/v1")

	// Customer routes
	customers := v1.Group("/customers")
	{
		customers.GET("", handlers.GetCustomers)
		customers.GET("/:id", handlers.GetCustomer)
		customers.POST("", handlers.CreateCustomer)
		customers.PUT("/:id", handlers.UpdateCustomer)
		customers.DELETE("/:id", handlers.DeleteCustomer)
	}

	// Car routes
	cars := v1.Group("/cars")
	{
		cars.GET("", handlers.GetCars)
		cars.GET("/:id", handlers.GetCar)
		cars.POST("", handlers.CreateCar)
		cars.PUT("/:id", handlers.UpdateCar)
		cars.DELETE("/:id", handlers.DeleteCar)
	}

	// Booking routes
	bookings := v1.Group("/bookings")
	{
		bookings.GET("", handlers.GetBookings)
		bookings.GET("/:id", handlers.GetBooking)
		bookings.POST("", handlers.CreateBooking)
		bookings.PUT("/:id", handlers.UpdateBooking)
		bookings.DELETE("/:id", handlers.DeleteBooking)
		bookings.PUT("/:id/finish", handlers.FinishBooking)
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Car Rental API",
		})
	})
}
