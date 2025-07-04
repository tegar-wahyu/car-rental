package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"car-rental/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {
	var customers []models.Customer
	result := database.DB.Order("no").Find(&customers)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customers})
}

// Helper function to find a customer by ID
func findCustomerByID(c *gin.Context) (*models.Customer, int, error) {
	id := c.Param("id")
	customerID, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var customer models.Customer
	result := database.DB.First(&customer, customerID)
	if result.Error != nil {
		return nil, http.StatusNotFound, result.Error
	}

	return &customer, http.StatusOK, nil
}

func GetCustomer(c *gin.Context) {
	customer, status, err := findCustomerByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid customer ID"})
		} else {
			c.JSON(status, gin.H{"error": "Customer not found"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": customer})
}

func UpdateCustomer(c *gin.Context) {
	customer, status, err := findCustomerByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid customer ID"})
		} else {
			c.JSON(status, gin.H{"error": "Customer not found"})
		}
		return
	}

	var updateData models.Customer
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	result := database.DB.Model(customer).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": customer})
}

func DeleteCustomer(c *gin.Context) {
	customer, status, err := findCustomerByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid customer ID"})
		} else {
			c.JSON(status, gin.H{"error": "Customer not found"})
		}
		return
	}

	constraints := utils.CheckCustomerBookingConstraints(customer.No)

	if constraints.HasBookings {
		details := map[string]interface{}{
			"active_bookings": constraints.ActiveBookings,
			"total_bookings":  constraints.TotalBookings,
		}
		
		if constraints.HasActive {
			utils.RespondWithConstraintError(c, "customer", customer.No, "active_bookings", details)
		} else {
			utils.RespondWithConstraintError(c, "customer", customer.No, "booking_history", details)
		}
		return
	}

	result := database.DB.Delete(customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer from database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// SubscribeToMembership subscribes a customer to a membership plan
func SubscribeToMembership(c *gin.Context) {
	customer, status, err := findCustomerByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid customer ID"})
		} else {
			c.JSON(status, gin.H{"error": "Customer not found"})
		}
		return
	}

	membershipIDStr := c.Param("membership_id")
	membershipID, err := strconv.Atoi(membershipIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid membership ID"})
		return
	}

	// Validate membership exists
	var membership models.Membership
	if err := database.DB.First(&membership, membershipID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Membership not found"})
		return
	}

	// Update customer membership
	result := database.DB.Model(customer).Update("membership_id", membershipID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe to membership"})
		return
	}

	// Reload customer with membership
	database.DB.Preload("Membership").First(customer, customer.No)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully subscribed to membership",
		"data":    customer,
	})
}

// UnsubscribeFromMembership removes a customer's membership
func UnsubscribeFromMembership(c *gin.Context) {
	customer, status, err := findCustomerByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid customer ID"})
		} else {
			c.JSON(status, gin.H{"error": "Customer not found"})
		}
		return
	}

	if customer.MembershipID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer is not subscribed to any membership"})
		return
	}

	// Remove membership
	result := database.DB.Model(customer).Update("membership_id", nil)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe from membership"})
		return
	}

	// Reload customer
	database.DB.First(customer, customer.No)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully unsubscribed from membership",
		"data":    customer,
	})
}
