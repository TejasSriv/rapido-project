package middlewares

import (
	"net/http"
	"strings"

	"rapido-backend/initializers"
	"rapido-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthRequired(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
		c.Abort()
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	userID, err := uuid.Parse(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		c.Abort()
		return
	}

	c.Set("user", user)

	c.Next()
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		if user.(models.User).Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admin role required."})
			c.Abort()
			return
		}

		c.Next()
	}
}
