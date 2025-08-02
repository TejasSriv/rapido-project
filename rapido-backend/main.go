package main

import (
	"rapido-backend/controllers"
	"rapido-backend/initializers"
	"rapido-backend/middlewares"

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

	authenticated := r.Group("/")
	authenticated.Use(middlewares.AuthRequired)
	{
		authenticated.GET("/users/:id", controllers.GetUserProfile)
		authenticated.PUT("/users/:id", controllers.UpdateUserProfile)

		authenticated.POST("/rides", controllers.CreateRide)
		authenticated.GET("/rides/:id", controllers.GetRideDetails)
		authenticated.GET("/users/:id/rides", controllers.GetUserRides)
		authenticated.PUT("/rides/:id/cancel", controllers.CancelRide)
	}

	admin := r.Group("/admin")
	admin.Use(middlewares.AuthRequired, middlewares.AdminRequired())
	{
		// Admin APIs
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Rapido Backend API",
		})
	})

	r.Run(":8080")
}
