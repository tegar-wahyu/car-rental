package utils

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingConstraintInfo struct {
	TotalBookings  int64 `json:"total_bookings"`
	ActiveBookings int64 `json:"active_bookings"`
	HasBookings    bool  `json:"has_bookings"`
	HasActive      bool  `json:"has_active"`
}

// Check if a customer has any bookings that would prevent deletion
func CheckCustomerBookingConstraints(customerID int) *BookingConstraintInfo {
	var totalBookings, activeBookings int64

	database.DB.Model(&models.Booking{}).Where("customer_id = ?", customerID).Count(&totalBookings)
	database.DB.Model(&models.Booking{}).Where("customer_id = ? AND finished = ?", customerID, false).Count(&activeBookings)

	return &BookingConstraintInfo{
		TotalBookings:  totalBookings,
		ActiveBookings: activeBookings,
		HasBookings:    totalBookings > 0,
		HasActive:      activeBookings > 0,
	}
}

// Check if a car has any bookings that would prevent deletion
func CheckCarBookingConstraints(carID int) *BookingConstraintInfo {
	var totalBookings, activeBookings int64

	database.DB.Model(&models.Booking{}).Where("cars_id = ?", carID).Count(&totalBookings)
	database.DB.Model(&models.Booking{}).Where("cars_id = ? AND finished = ?", carID, false).Count(&activeBookings)

	return &BookingConstraintInfo{
		TotalBookings:  totalBookings,
		ActiveBookings: activeBookings,
		HasBookings:    totalBookings > 0,
		HasActive:      activeBookings > 0,
	}
}

// Check if a driver has any bookings that would prevent deletion
func CheckDriverBookingConstraints(driverID int) *BookingConstraintInfo {
	var totalBookings, activeBookings int64

	database.DB.Model(&models.Booking{}).Where("driver_id = ?", driverID).Count(&totalBookings)
	database.DB.Model(&models.Booking{}).Where("driver_id = ? AND finished = ?", driverID, false).Count(&activeBookings)

	return &BookingConstraintInfo{
		TotalBookings:  totalBookings,
		ActiveBookings: activeBookings,
		HasBookings:    totalBookings > 0,
		HasActive:      activeBookings > 0,
	}
}

// Structured error response for FK constraint violations
type ReferentialIntegrityError struct {
	Message    string                 `json:"error"`
	EntityType string                 `json:"entity_type"`
	EntityID   int                    `json:"entity_id"`
	Constraint string                 `json:"constraint"`
	Details    map[string]interface{} `json:"details"`
}

func RespondWithConstraintError(c *gin.Context, entityType string, entityID int, constraint string, details map[string]interface{}) {
	var message string

	switch constraint {
	case "active_bookings":
		message = "Cannot delete " + entityType + " with active bookings. Please finish or cancel active bookings first."
	case "finished_booking":
		message = "Cannot delete finished booking. Finished bookings are kept for historical records."
	default:
		message = "Cannot delete " + entityType + " due to referential integrity constraints."
	}

	errorResponse := ReferentialIntegrityError{
		Message:    message,
		EntityType: entityType,
		EntityID:   entityID,
		Constraint: constraint,
		Details:    details,
	}

	c.JSON(http.StatusBadRequest, errorResponse)
}


func SoftDeleteCar(carID int) error {
	now := time.Now()
	return database.DB.Model(&models.Car{}).Where("no = ?", carID).Update("deleted_at", &now).Error
}

func SoftDeleteDriver(driverID int) error {
	now := time.Now()
	return database.DB.Model(&models.Driver{}).Where("no = ?", driverID).Update("deleted_at", &now).Error
}

func SoftDeleteDriverIncentive(incentiveID int) error {
	now := time.Now()
	return database.DB.Model(&models.DriverIncentive{}).Where("no = ?", incentiveID).Update("deleted_at", &now).Error
}

func SoftDeleteDriverIncentivesByBookingID(bookingID int) error {
	now := time.Now()
	return database.DB.Model(&models.DriverIncentive{}).Where("booking_id = ?", bookingID).Update("deleted_at", &now).Error
}

func SoftDeleteCustomer(customerID int) error {
	now := time.Now()
	return database.DB.Model(&models.Customer{}).Where("no = ?", customerID).Update("deleted_at", &now).Error
}

func RespondWithSoftDeleteSuccess(c *gin.Context, entityType string, entityID int) {
	c.JSON(http.StatusOK, gin.H{
		"message": entityType + " has been soft deleted successfully",
		"details": map[string]interface{}{
			"soft_deleted_entity_id": entityID,
			"entity_type":            entityType,
			"deleted_at":             time.Now(),
		},
	})
}
