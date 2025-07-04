package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	result := database.DB.Preload("Customer").Preload("Car").Order("no").Find(&bookings)

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
	result := database.DB.Preload("Customer").Preload("Car").First(&booking, bookingID)
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
	if err := database.DB.First(&customer, booking.CustomerID).Error; err != nil {
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

	// Calculate total cost
	days := int(booking.EndRent.Sub(booking.StartRent).Hours()/24) + 1
	booking.TotalCost = float64(days) * car.DailyRent

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
	database.DB.Preload("Customer").Preload("Car").First(&booking, booking.No)

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

		// Recalculate total cost
		var car models.Car
		database.DB.First(&car, booking.CarsID)
		days := int(endRent.Sub(startRent).Hours()/24) + 1
		totalCost := float64(days) * car.DailyRent
		updateData.TotalCost = &totalCost
	}

	// Update only provided fields
	result := database.DB.Model(booking).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking"})
		return
	}

	// Reload booking with relationships
	database.DB.Preload("Customer").Preload("Car").First(booking, booking.No)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete a finished booking"})
		return
	}

	// Start a transaction
	tx := database.DB.Begin()

	// Delete booking
	result := tx.Delete(booking)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	// Restore car stock
	var car models.Car
	if err := tx.First(&car, booking.CarsID).Error; err == nil {
		tx.Model(&car).Update("stock", car.Stock+1)
	}

	// Commit transaction
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})
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

	// Commit transaction
	tx.Commit()

	// Reload booking with relationships
	database.DB.Preload("Customer").Preload("Car").First(booking, booking.No)

	c.JSON(http.StatusOK, gin.H{"data": booking, "message": "Booking finished successfully"})
}
