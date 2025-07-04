package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCars(c *gin.Context) {
	var cars []models.Car
	result := database.DB.Order("no").Find(&cars)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cars"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cars})
}

// Helper function to find a car by ID
func findCarByID(c *gin.Context) (*models.Car, int, error) {
	id := c.Param("id")
	carID, err := strconv.Atoi(id)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var car models.Car
	result := database.DB.First(&car, carID)
	if result.Error != nil {
		return nil, http.StatusNotFound, result.Error
	}

	return &car, http.StatusOK, nil
}

func GetCar(c *gin.Context) {
	car, status, err := findCarByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid car ID"})
		} else {
			c.JSON(status, gin.H{"error": "Car not found"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": car})
}

func CreateCar(c *gin.Context) {
	var car models.Car

	if err := c.ShouldBindJSON(&car); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Create(&car)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create car"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": car})
}

func UpdateCar(c *gin.Context) {
	car, status, err := findCarByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid car ID"})
		} else {
			c.JSON(status, gin.H{"error": "Car not found"})
		}
		return
	}

	var updateData models.Car
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	result := database.DB.Model(car).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": car})
}

func DeleteCar(c *gin.Context) {
	car, status, err := findCarByID(c)
	if err != nil {
		if status == http.StatusBadRequest {
			c.JSON(status, gin.H{"error": "Invalid car ID"})
		} else {
			c.JSON(status, gin.H{"error": "Car not found"})
		}
		return
	}

	result := database.DB.Delete(car)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete car"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Car deleted successfully"})
}
