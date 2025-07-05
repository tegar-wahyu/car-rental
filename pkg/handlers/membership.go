package handlers

import (
	"car-rental/pkg/database"
	"car-rental/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMemberships(c *gin.Context) {
	var memberships []models.Membership
	result := database.DB.Find(&memberships)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve memberships"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": memberships})
}

func GetMembership(c *gin.Context) {
	id := c.Param("id")
	membershipID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid membership ID"})
		return
	}

	var membership models.Membership
	result := database.DB.First(&membership, membershipID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Membership not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": membership})
}
