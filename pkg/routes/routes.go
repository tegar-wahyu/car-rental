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

		// Membership subscription endpoints
		customers.PUT("/:id/subscribe/:membership_id", handlers.SubscribeToMembership)
		customers.DELETE("/:id/unsubscribe", handlers.UnsubscribeFromMembership)
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

	// Booking routes
	bookings := r.Group("/bookings")
	{
		bookings.GET("", handlers.GetBookings)
		bookings.GET("/:id", handlers.GetBooking)
		bookings.POST("", handlers.CreateBooking)
		bookings.PUT("/:id", handlers.UpdateBooking)
		bookings.DELETE("/:id", handlers.DeleteBooking)
		bookings.PUT("/:id/finish", handlers.FinishBooking)

		// Booking Type sub-routes (read-only)
		bookings.GET("/types", handlers.GetBookingTypes)
		bookings.GET("/types/:id", handlers.GetBookingType)
	}

	// Membership routes (read-only)
	memberships := r.Group("/memberships")
	{
		memberships.GET("", handlers.GetMemberships)
		memberships.GET("/:id", handlers.GetMembership)
	}

	// Driver routes
	drivers := r.Group("/drivers")
	{
		drivers.GET("", handlers.GetDrivers)
		drivers.GET("/:id", handlers.GetDriver)
		drivers.POST("", handlers.CreateDriver)
		drivers.PUT("/:id", handlers.UpdateDriver)
		drivers.DELETE("/:id", handlers.DeleteDriver)
		drivers.GET("/:id/incentives", handlers.GetDriverIncentives)
	}

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Car Rental API",
		})
	})
}
