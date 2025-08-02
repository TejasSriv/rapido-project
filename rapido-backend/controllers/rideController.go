package controllers

import (
	"log"
	"net/http"

	"rapido-backend/initializers"
	"rapido-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateRide(c *gin.Context) {
	var body struct {
		UserId      string `json:"userId" binding:"required"`
		Origin      string `json:"origin" binding:"required"`
		Destination string `json:"destination" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := uuid.Parse(body.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ride := models.Ride{
		UserID:          userId,
		PickupLocation:  body.Origin,
		DropoffLocation: body.Destination,
		CurrentStatus:   "pending",
		Fare:            100,
	}

	result := initializers.DB.Create(&ride)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ride", "details": result.Error.Error()})
		log.Printf("Error creating ride: %v", result.Error)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Ride created successfully", "ride": ride})

}

func GetRideDetails(c *gin.Context) {
	rideIDstr := c.Param("id")
	rideID, err := uuid.Parse(rideIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID format"})
		return
	}
	var ride models.Ride
	result := initializers.DB.First(&ride, rideID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ride details"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"ride": ride})

}

func GetUserRides(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var rides []models.Ride
	result := initializers.DB.Where("user_id = ?", userID).Find(&rides)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user rides"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rides": rides})
}

func CancelRide(c *gin.Context) {
	rideIDStr := c.Param("id")
	rideID, err := uuid.Parse(rideIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	// checking if ride exists and its status
	var ride models.Ride
	result := initializers.DB.First(&ride, rideID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ride not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ride details"})
		}
		return
	}

	// only cancel if ride is pending ro accepted
	if ride.CurrentStatus != "pending" && ride.CurrentStatus != "accepted" {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot cancel a ride."})
		return
	}

	// update ride status
	updateResult := initializers.DB.Model(&ride).Update("current_status", "cancelled")
	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel ride"})
		log.Printf("Error cancelling ride %s: %v", rideID, updateResult.Error)
		return
	}

	initializers.DB.First(&ride, rideID)

	c.JSON(http.StatusOK, gin.H{"message": "Ride cancelled successfully", "ride": ride})
}
