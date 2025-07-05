package routes

import (
	"car-rental/pkg/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Create API v1 group
	v1 := r.Group("/api/v1")

	// Create API v2 group
	v2 := r.Group("/api/v2")

	// Customer routes for v1 - Basic CRUD only
	customersV1 := v1.Group("/customers")
	{
		customersV1.GET("", handlers.GetCustomers)
		customersV1.GET("/:id", handlers.GetCustomer)
		customersV1.POST("", handlers.CreateCustomer)
		customersV1.PUT("/:id", handlers.UpdateCustomer)
		customersV1.DELETE("/:id", handlers.DeleteCustomer)
	}

	// Customer routes for v2 - Full functionality
	customersV2 := v2.Group("/customers")
	{
		customersV2.GET("", handlers.GetCustomers)
		customersV2.GET("/:id", handlers.GetCustomer)
		customersV2.POST("", handlers.CreateCustomer)
		customersV2.PUT("/:id", handlers.UpdateCustomer)
		customersV2.DELETE("/:id", handlers.DeleteCustomer)

		// Membership subscription endpoints
		customersV2.PUT("/:id/subscribe/:membership_id", handlers.SubscribeToMembership)
		customersV2.DELETE("/:id/unsubscribe", handlers.UnsubscribeFromMembership)
	}

	// Car routes for v1 - Basic CRUD only
	carsV1 := v1.Group("/cars")
	{
		carsV1.GET("", handlers.GetCars)
		carsV1.GET("/:id", handlers.GetCar)
		carsV1.POST("", handlers.CreateCar)
		carsV1.PUT("/:id", handlers.UpdateCar)
		carsV1.DELETE("/:id", handlers.DeleteCar)
	}

	// Car routes for v2 - Full functionality
	carsV2 := v2.Group("/cars")
	{
		carsV2.GET("", handlers.GetCars)
		carsV2.GET("/:id", handlers.GetCar)
		carsV2.POST("", handlers.CreateCar)
		carsV2.PUT("/:id", handlers.UpdateCar)
		carsV2.DELETE("/:id", handlers.DeleteCar)
	}

	// Booking routes for v1 - Basic CRUD only
	bookingsV1 := v1.Group("/bookings")
	{
		bookingsV1.GET("", handlers.GetBookings)
		bookingsV1.GET("/:id", handlers.GetBooking)
		bookingsV1.POST("", handlers.CreateBooking)
		bookingsV1.PUT("/:id", handlers.UpdateBooking)
		bookingsV1.PUT("/:id/finish", handlers.FinishBooking)

		bookingsV1.DELETE("/:id", handlers.DeleteBooking)
	}

	// Booking routes for v2 - Full functionality
	bookingsV2 := v2.Group("/bookings")
	{
		bookingsV2.GET("", handlers.GetBookings)
		bookingsV2.GET("/:id", handlers.GetBooking)
		bookingsV2.POST("", handlers.CreateBooking)
		bookingsV2.PUT("/:id", handlers.UpdateBooking)
		bookingsV2.DELETE("/:id", handlers.DeleteBooking)
		bookingsV2.PUT("/:id/finish", handlers.FinishBooking)

		// Booking Type sub-routes (read-only)
		bookingsV2.GET("/types", handlers.GetBookingTypes)
		bookingsV2.GET("/types/:id", handlers.GetBookingType)
	}

	// Membership routes (v2 only)
	membershipsV2 := v2.Group("/memberships")
	{
		membershipsV2.GET("", handlers.GetMemberships)
		membershipsV2.GET("/:id", handlers.GetMembership)
	}

	// Driver routes (v2 only)
	driversV2 := v2.Group("/drivers")
	{
		driversV2.GET("", handlers.GetDrivers)
		driversV2.GET("/:id", handlers.GetDriver)
		driversV2.POST("", handlers.CreateDriver)
		driversV2.PUT("/:id", handlers.UpdateDriver)
		driversV2.DELETE("/:id", handlers.DeleteDriver)
		driversV2.GET("/:id/incentives", handlers.GetDriverIncentives)
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Car Rental API",
		})
	})
}
