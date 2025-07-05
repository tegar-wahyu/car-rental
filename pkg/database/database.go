package database

import (
	"car-rental/pkg/models"
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Set defaults if not provided
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbName == "" {
		dbName = "car_rental"
	}
	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	// Build PostgreSQL DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	// Connect to database
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully")
}

func Migrate() {
	// Check if migrations should be run
	shouldMigrate := os.Getenv("AUTO_MIGRATE")
	if strings.ToLower(shouldMigrate) != "true" {
		fmt.Println("Skipping database migration based on AUTO_MIGRATE setting")
		return
	}

	err := DB.AutoMigrate(&models.Customer{}, &models.Car{}, &models.Booking{}, &models.Membership{},
		&models.BookingType{}, &models.DriverIncentive{}, &models.Driver{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migration completed")
}

// MigrateWithFeedback performs database migration and returns an error instead of fatal
func MigrateWithFeedback() error {
	err := DB.AutoMigrate(&models.Customer{}, &models.Car{}, &models.Booking{}, &models.Membership{},
		&models.BookingType{}, &models.DriverIncentive{}, &models.Driver{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	fmt.Println("Database migration completed")
	return nil
}
