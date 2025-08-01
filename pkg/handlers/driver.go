package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"car-rental/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDrivers(c *gin.Context) {
	var drivers []models.Driver
	result := database.DB.Where("deleted_at IS NULL").Find(&drivers)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve drivers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": drivers})
}

// Helper function to find a driver by ID
func findDriverByID(c *gin.Context) (*models.Driver, int, error) {
	id := c.Param("id")
	driverID, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var driver models.Driver
	result := database.DB.Where("deleted_at IS NULL").First(&driver, driverID)
	if result.Error != nil {
		return nil, http.StatusNotFound, result.Error
	}

	return &driver, http.StatusOK, nil
}

func GetDriver(c *gin.Context) {
	driver, status, err := findDriverByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid driver ID"})
		} else {
			c.JSON(status, gin.H{"error": "Driver not found"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": driver})
}

func CreateDriver(c *gin.Context) {
	var driver models.Driver

	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&driver)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create driver"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": driver})
}

func UpdateDriver(c *gin.Context) {
	driver, status, err := findDriverByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid driver ID"})
		} else {
			c.JSON(status, gin.H{"error": "Driver not found"})
		}
		return
	}

	var UpdateDriver models.Driver
	if err := c.ShouldBindJSON(&UpdateDriver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Model(driver).Updates(UpdateDriver)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update driver"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": driver})
}

func DeleteDriver(c *gin.Context) {
	driver, status, err := findDriverByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid driver ID"})
		} else {
			c.JSON(status, gin.H{"error": "Driver not found"})
		}
		return
	}

	constraints := utils.CheckDriverBookingConstraints(driver.No)

	if constraints.HasActive {
		details := map[string]interface{}{
			"active_bookings": constraints.ActiveBookings,
			"total_bookings":  constraints.TotalBookings,
		}
		utils.RespondWithConstraintError(c, "driver", driver.No, "active_bookings", details)
		return
	}

	// Use soft delete instead of hard delete
	err = utils.SoftDeleteDriver(driver.No)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete driver"})
		return
	}

	utils.RespondWithSoftDeleteSuccess(c, "driver", driver.No)
}

func GetDriverIncentives(c *gin.Context) {
	driver, status, err := findDriverByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid driver ID"})
		} else {
			c.JSON(status, gin.H{"error": "Driver not found"})
		}
		return
	}

	var incentives []models.DriverIncentive
	result := database.DB.Preload("Booking").Preload("Booking.Customer").Preload("Booking.Car").
		Joins("JOIN bookings ON driver_incentives.booking_id = bookings.no").
		Where("bookings.driver_id = ? AND driver_incentives.deleted_at IS NULL", driver.No).
		Find(&incentives)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve driver incentives"})
		return
	}

	// Calculate total incentives
	var totalIncentives float64
	for _, incentive := range incentives {
		totalIncentives += incentive.Incentive
	}

	c.JSON(http.StatusOK, gin.H{
		"data":             incentives,
		"total_incentives": totalIncentives,
	})
}
