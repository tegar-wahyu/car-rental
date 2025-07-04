package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBookingTypes(c *gin.Context) {
	var bookingTypes []models.BookingType
	result := database.DB.Find(&bookingTypes)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve booking types"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingTypes})
}

// Helper function to find a booking type by ID
func findBookingTypeByID(c *gin.Context) (*models.BookingType, int, error) {
	id := c.Param("id")
	bookingTypeID, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var bookingType models.BookingType
	result := database.DB.First(&bookingType, bookingTypeID)
	if result.Error != nil {
		return nil, http.StatusNotFound, result.Error
	}

	return &bookingType, http.StatusOK, nil
}

func GetBookingType(c *gin.Context) {
	bookingType, status, err := findBookingTypeByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid booking type ID"})
		} else {
			c.JSON(status, gin.H{"error": "Booking type not found"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingType})
}

func CreateBookingType(c *gin.Context) {
	var bookingType models.BookingType

	if err := c.ShouldBindJSON(&bookingType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&bookingType)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking type"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": bookingType})
}

func UpdateBookingType(c *gin.Context) {
	bookingType, status, err := findBookingTypeByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid booking type ID"})
		} else {
			c.JSON(status, gin.H{"error": "Booking type not found"})
		}
		return
	}

	if err := c.ShouldBindJSON(bookingType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Save(bookingType)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": bookingType})
}

func DeleteBookingType(c *gin.Context) {
	bookingType, status, err := findBookingTypeByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid booking type ID"})
		} else {
			c.JSON(status, gin.H{"error": "Booking type not found"})
		}
		return
	}

	// Check if booking type is being used
	var bookingCount int64
	database.DB.Model(&models.Booking{}).Where("booking_type_id = ?", bookingType.No).Count(&bookingCount)

	if bookingCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete booking type that is currently in use"})
		return
	}

	result := database.DB.Delete(bookingType)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking type deleted successfully"})
}
