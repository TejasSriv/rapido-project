package controllers

import (
	"log"
	"net/http"

	"rapido-backend/initializers"
	"rapido-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func UserSignup(c *gin.Context) {
	// data from request body

	var body struct {
		Username    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=6"`
		FullName    string `json:"fullName"`
		PhoneNumber string `json:"phoneNumber"`
		Role        string `json:"role"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//password hashing

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		log.Printf("Error hashing password: %v", err)
		return
	}

	//user creation

	user := models.User{
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: string(hash),
		FullName:     body.FullName,
		PhoneNumber:  body.PhoneNumber,
		Role:         "user",
	}

	if body.Role == "admin" {
		user.Role = "admin"
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		if result.Error.Error() == "ERROR: duplicate key value violates unique constraint \"idx_users_email\" (SQLSTATE 23505)" || result.Error.Error() == "ERROR: duplicate key value violates unique constraint \"idx_users_username\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": result.Error.Error()})
		log.Printf("Error creating user: %v", result.Error)
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})

}

//user login

func UserLogin(c *gin.Context) {
	var body struct {
		Email    string `json:"email" binding:"required, email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//search user by email

	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	//generating uuid token only for testing, JWT could be used for production

	token := uuid.New().String()

	c.JSON(http.StatusOK, gin.H{
		"message": "Login Successful",
		"token":   token,
		"userId":  user.ID,
		"role":    user.Role,
	})

}

//fetch user profile

func GetUserProfile(c *gin.Context) {
	userIDstr := c.Param("id")
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with the given ID"})
		return
	}

	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUserProfile(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	var body struct {
		Username    string `json:"username"`
		Email       string `json:"email" binding:"email"`
		FullName    string `json:"fullName"`
		PhoneNumber string `json:"phoneNumber"`
		Password    string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with the given ID"})
		return
	}

	updates := make(map[string]interface{})
	if body.Username != "" {
		updates["username"] = body.Username
	}
	if body.Email != "" {
		updates["email"] = body.Email
	}
	if body.FullName != "" {
		updates["full_name"] = body.FullName
	}
	if body.PhoneNumber != "" {
		updates["phone_number"] = body.PhoneNumber
	}
	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			log.Printf("Error hashing password: %v", err)
			return
		}
		updates["password_hash"] = string(hash)
	}

	updateResult := initializers.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates)

	if updateResult.Error != nil {
		if updateResult.Error.Error() == "ERROR: duplicate key value violates unique constraint \"idx_users_email\" (SQLSTATE 23505)" || updateResult.Error.Error() == "ERROR: duplicate key value violates unique constraint \"idx_users_username\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile", "details": updateResult.Error.Error()})
		log.Printf("Error updating user profile: %v", updateResult.Error)
		return
	}

	if updateResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No user found with the given ID"})
		return
	}

	initializers.DB.First(&user, userID)
	user.PasswordHash = ""
	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully", "user": user})

}
