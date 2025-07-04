package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"car-rental/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Booking type constants
const (
	BOOKING_TYPE_CAR_ONLY   = "Car Only"
	BOOKING_TYPE_CAR_DRIVER = "Car & Driver"
)

func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	result := database.DB.Preload("Customer").Preload("Customer.Membership").Preload("Car").Preload("Driver").Preload("BookingType").Order("no").Find(&bookings)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

// Helper function to find a booking by ID
func findBookingByID(c *gin.Context) (*models.Booking, int, error) {
	id := c.Param("id")
	bookingID, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var booking models.Booking
	result := database.DB.Preload("Customer").Preload("Customer.Membership").Preload("Car").Preload("Driver").Preload("BookingType").First(&booking, bookingID)
	if result.Error != nil {
		return nil, http.StatusNotFound, result.Error
	}

	return &booking, http.StatusOK, nil
}

func GetBooking(c *gin.Context) {
	booking, status, err := findBookingByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid booking ID"})
		} else {
			c.JSON(status, gin.H{"error": "Booking not found"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func CreateBooking(c *gin.Context) {
	var booking models.Booking

	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that customer exists
	var customer models.Customer
	if err := database.DB.Preload("Membership").First(&customer, booking.CustomerID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer not found"})
		return
	}

	// Validate that car exists and is available
	var car models.Car
	if err := database.DB.First(&car, booking.CarsID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Car not found"})
		return
	}

	if car.Stock <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Car is not available for booking"})
		return
	}

	// Validate dates
	if booking.StartRent.After(booking.EndRent) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
		return
	}

	if booking.StartRent.Before(time.Now().Truncate(24 * time.Hour)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date cannot be in the past"})
		return
	}

	// Calculate total cost (days * daily_rent)
	days := int(booking.EndRent.Sub(booking.StartRent).Hours()/24) + 1
	booking.TotalCost = float64(days) * car.DailyRent

	// Calculate membership discount
	var discount float64 = 0
	if customer.Membership != nil {
		discount = booking.TotalCost * (customer.Membership.Discount / 100)
	}
	booking.Discount = discount
	// Validate booking type
	var bookingType models.BookingType
	if err := database.DB.First(&bookingType, booking.BookingTypeID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking type not found"})
		return
	}

	// Validate driver assignment based on booking type
	if booking.DriverID != nil && bookingType.BookingType != BOOKING_TYPE_CAR_DRIVER {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Driver can only be assigned for 'Car & Driver' booking type"})
		return
	}

	if booking.DriverID == nil && bookingType.BookingType == BOOKING_TYPE_CAR_DRIVER {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Driver must be assigned for 'Car & Driver' booking type"})
		return
	}

	// Calculate driver cost if driver is requested
	var totalDriverCost float64 = 0
	var driver models.Driver
	if booking.DriverID != nil {
		if err := database.DB.First(&driver, *booking.DriverID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Driver not found"})
			return
		}
		totalDriverCost = float64(days) * driver.DailyCost
	}
	booking.TotalDriverCost = totalDriverCost

	// Start a transaction
	tx := database.DB.Begin()

	// Create booking
	if err := tx.Create(&booking).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	// Update car stock
	if err := tx.Model(&car).Update("stock", car.Stock-1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update car stock"})
		return
	}

	// Commit transaction
	tx.Commit()

	// Reload booking with relationships
	database.DB.Preload("Customer").Preload("Customer.Membership").Preload("Car").Preload("Driver").Preload("BookingType").First(&booking, booking.No)

	c.JSON(http.StatusCreated, gin.H{"data": booking})
}

func UpdateBooking(c *gin.Context) {
	booking, status, err := findBookingByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid booking ID"})
		} else {
			c.JSON(status, gin.H{"error": "Booking not found"})
		}
		return
	}

	var updateData models.BookingUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Don't allow updating if booking is finished
	if booking.Finished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update a finished booking"})
		return
	}

	// If updating dates, validate and recalculate cost
	if updateData.StartRent != nil || updateData.EndRent != nil {
		startRent := booking.StartRent
		endRent := booking.EndRent

		if updateData.StartRent != nil {
			startRent = *updateData.StartRent
		}
		if updateData.EndRent != nil {
			endRent = *updateData.EndRent
		}

		if startRent.After(endRent) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
			return
		}

		if startRent.Before(time.Now().Truncate(24 * time.Hour)) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start date cannot be in the past"})
			return
		}

		// Recalculate total cost with membership and driver considerations
		var car models.Car
		var customer models.Customer
		database.DB.First(&car, booking.CarsID)
		database.DB.Preload("Membership").First(&customer, booking.CustomerID)

		days := int(endRent.Sub(startRent).Hours()/24) + 1

		// Total cost is always days * daily_rent
		totalCost := float64(days) * car.DailyRent
		updateData.TotalCost = &totalCost

		// Calculate membership discount
		var discount float64 = 0
		if customer.Membership != nil {
			discount = totalCost * (customer.Membership.Discount / 100)
		}
		updateData.Discount = &discount

		// Calculate driver cost if driver is assigned
		var totalDriverCost float64 = 0
		if booking.DriverID != nil {
			var driver models.Driver
			database.DB.First(&driver, *booking.DriverID)
			totalDriverCost = float64(days) * driver.DailyCost
		}
		updateData.TotalDriverCost = &totalDriverCost
	}

	// Update only provided fields
	result := database.DB.Model(booking).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	// Reload booking with relationships
	database.DB.Preload("Customer").Preload("Customer.Membership").Preload("Car").Preload("Driver").Preload("BookingType").First(booking, booking.No)

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func DeleteBooking(c *gin.Context) {
	booking, status, err := findBookingByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid booking ID"})
		} else {
			c.JSON(status, gin.H{"error": "Booking not found"})
		}
		return
	}

	// Don't allow deleting finished bookings
	if booking.Finished {
		details := map[string]interface{}{
			"booking_id": booking.No,
			"finished":   booking.Finished,
		}
		utils.RespondWithConstraintError(c, "booking", booking.No, "finished_booking", details)
		return
	}

	// Start a transaction
	tx := database.DB.Begin()

	// Delete booking
	result := tx.Delete(booking)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking from database"})
		return
	}

	// Restore car stock
	var car models.Car
	if err := tx.First(&car, booking.CarsID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find associated car for stock restoration"})
		return
	}

	if err := tx.Model(&car).Update("stock", car.Stock+1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore car stock"})
		return
	}

	// Commit transaction
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking deleted successfully and car stock restored",
		"details": map[string]interface{}{
			"deleted_booking_id": booking.No,
			"restored_car_stock": car.Stock + 1,
		},
	})
}

// FinishBooking marks a booking as finished and restores car stock
func FinishBooking(c *gin.Context) {
	booking, status, err := findBookingByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid booking ID"})
		} else {
			c.JSON(status, gin.H{"error": "Booking not found"})
		}
		return
	}

	if booking.Finished {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking is already finished"})
		return
	}

	// Start a transaction
	tx := database.DB.Begin()

	// Mark booking as finished
	if err := tx.Model(booking).Update("finished", true).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finish booking"})
		return
	}

	// Restore car stock
	var car models.Car
	if err := tx.First(&car, booking.CarsID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find car"})
		return
	}

	if err := tx.Model(&car).Update("stock", car.Stock+1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore car stock"})
		return
	}

	// Calculate and save driver incentive if driver was assigned
	if booking.DriverID != nil {
		days := int(booking.EndRent.Sub(booking.StartRent).Hours()/24) + 1
		baseCost := float64(days) * car.DailyRent
		incentive := baseCost * 0.05 // 5% of base cost (days * daily_rent)

		driverIncentive := models.DriverIncentive{
			BookingID: booking.No,
			Incentive: incentive,
		}

		if err := tx.Create(&driverIncentive).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver incentive"})
			return
		}
	}

	// Commit transaction
	tx.Commit()

	// Reload booking with relationships
	database.DB.Preload("Customer").Preload("Customer.Membership").Preload("Car").Preload("Driver").Preload("BookingType").First(booking, booking.No)

	c.JSON(http.StatusOK, gin.H{"data": booking, "message": "Booking finished successfully"})
}
