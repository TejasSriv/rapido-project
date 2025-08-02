package controllers

import (
	"log"
	"net/http"
	"time"

	"rapido-backend/initializers"
	"rapido-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ViewAllRides(c *gin.Context) {
	var rides []models.Ride
	result := initializers.DB.Preload("User").Find(&rides)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch rides"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rides": rides})
}

func AdminUpdateRideStatus(c *gin.Context) {
	rideIDStr := c.Param("id")
	rideID, err := uuid.Parse(rideIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ride ID"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required,oneof=accepted rejected completed"`
		Notes  string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if ride.CurrentStatus == "cancelled" {
		c.JSON(http.StatusConflict, gin.H{"error": "Cannot update status of a cancelled ride"})
		return
	}

	updateMap := make(map[string]any)
	updateMap["current_status"] = body.Status
	updateMap["admin_notes"] = body.Notes

	switch body.Status {
	case "accepted":
		updateMap["accepted_at"] = time.Now()
		var driver models.User
		initializers.DB.Where("role = ?", "admin").First(&driver)
		if driver.ID != uuid.Nil {
			updateMap["driver_id"] = driver.ID
		}
	case "completed":
		updateMap["completed_at"] = time.Now()
	}

	updateResult := initializers.DB.Model(&ride).Updates(updateMap)
	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ride status"})
		log.Printf("Error updating ride status %s: %v", rideID, updateResult.Error)
		return
	}

	var admin models.User
	user, exists := c.Get("user")
	if exists {
		admin = user.(models.User)
	}

	adminAction := models.AdminAction{
		AdminID:       admin.ID,
		RideID:        ride.ID,
		ActionType:    body.Status,
		ActionDetails: body.Notes,
	}

	initializers.DB.Create(&adminAction)

	initializers.DB.First(&ride, rideID)
	c.JSON(http.StatusOK, gin.H{"message": "Ride status updated successfully", "ride": ride})
}

func GetRideAnalytics(c *gin.Context) {
	var results []struct {
		Date       time.Time `json:"date"`
		TotalRides int64     `json:"totalRides"`
	}

	initializers.DB.
		Model(&models.Ride{}).
		Select("DATE(requested_at) as date, count(*) as total_rides").
		Group("DATE(requested_at)").
		Order("date").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{"analytics": results})
}

func FilterRides(c *gin.Context) {
	var rides []models.Ride

	userIDStr := c.Query("userId")
	status := c.Query("status")
	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	db := initializers.DB.Model(&models.Ride{})

	if userIDStr != "" {
		db = db.Where("user_id = ?", userIDStr)
	}
	if status != "" {
		db = db.Where("current_status = ?", status)
	}
	if startDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err == nil {
			db = db.Where("requested_at >= ?", startDate)
		}
	}
	if endDateStr != "" {
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err == nil {
			db = db.Where("requested_at <= ?", endDate.Add(24*time.Hour))
		}
	}

	db = db.Preload("User")
	result := db.Find(&rides)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to filter rides"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"rides": rides})
}
