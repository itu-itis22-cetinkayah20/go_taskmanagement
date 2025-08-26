package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go_taskmanagement/models"
)

var DB *gorm.DB
var IsConnected bool

// Connect initializes the database connection
func Connect() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "1234"),
		getEnv("DB_NAME", "go_taskmanagement"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Running in in-memory mode")
		IsConnected = false
		return
	}

	IsConnected = true
	log.Println("Database connected successfully")
}

// ConnectTest initializes the test database connection
func ConnectTest() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("TEST_DB_HOST", "localhost"),
		getEnv("TEST_DB_USER", "postgres"),
		getEnv("TEST_DB_PASSWORD", "1234"),
		getEnv("TEST_DB_NAME", "go_taskmanagement_test"),
		getEnv("TEST_DB_PORT", "5432"),
		getEnv("TEST_DB_SSLMODE", "disable"),
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Quiet during tests
	})

	if err != nil {
		log.Printf("Failed to connect to test database: %v", err)
		log.Println("Running tests in in-memory mode")
		IsConnected = false
		return
	}

	IsConnected = true
	log.Println("Test database connected successfully")
}

// Migrate runs the database migrations
func Migrate() {
	if !IsConnected {
		return
	}

	err := DB.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}

// CleanTestData cleans test data from database
func CleanTestData() {
	if !IsConnected {
		return
	}

	DB.Where("email LIKE ?", "%@example.com").Delete(&models.User{})
	DB.Where("title LIKE ?", "Test %").Delete(&models.Task{})
}

// SeedTestData seeds initial test data
func SeedTestData() {
	if !IsConnected {
		return
	}

	// Skip public tasks seeding to avoid foreign key constraint issues
	// Public tasks will be handled differently or users can create them manually
	log.Println("Test data seeding completed (public tasks skipped to avoid FK constraints)")
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// Close closes the database connection
func Close() {
	if !IsConnected {
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Error getting database instance: %v", err)
		return
	}
	sqlDB.Close()
}
