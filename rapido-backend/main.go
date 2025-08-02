package main

import (
	"rapido-backend/controllers"
	"rapido-backend/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.UserSignup)
	r.POST("/login", controllers.UserLogin)

	r.GET("/users/:id", controllers.GetUserProfile)
	r.PUT("/users/:id", controllers.UpdateUserProfile)

	r.POST("/rides", controllers.CreateRide)
	r.GET("/rides/:id", controllers.GetRideDetails)
	r.GET("/users/:id/rides", controllers.GetUserRides)
	r.PUT("/rides/:id/cancel", controllers.CancelRide)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Rapido Backend API",
		})
	})

	r.Run(":8080")
}
