package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rapido-backend/models"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Connected to the database successfully")

	err = DB.AutoMigrate(&models.User{}, &models.Ride{}, &models.AdminAction{})

	if err != nil {
		log.Fatalf("Error migrating database: %v", err)
	} else {
		fmt.Println("Database migrated successfully")
	}
}
