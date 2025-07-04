package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"fmt"
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

	// Validate customer data
	if err := validateCustomer(customer); err != nil {
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

	// Validate update data before database operation
	if err := validateCustomer(updateData); err != nil {
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

	result := database.DB.Delete(customer)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete customer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted successfully"})
}

// validateCustomer validates customer data before database operations
func validateCustomer(customer models.Customer) error {
	if len(customer.NIK) > 16 {
		return fmt.Errorf("NIK must be at most 16 characters, got %d characters", len(customer.NIK))
	}
	if len(customer.PhoneNumber) > 15 {
		return fmt.Errorf("phone number must be at most 15 characters, got %d characters", len(customer.PhoneNumber))
	}
	return nil
}
