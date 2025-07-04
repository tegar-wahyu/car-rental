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
