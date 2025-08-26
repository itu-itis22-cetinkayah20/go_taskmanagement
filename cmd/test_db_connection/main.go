package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Test main database connection
	mainDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "1234"),
		getEnv("DB_NAME", "go_taskmanagement"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_SSLMODE", "disable"),
	)

	fmt.Println("Testing main database connection...")
	fmt.Printf("DSN: %s\n", mainDSN)

	mainDB, err := gorm.Open(postgres.Open(mainDSN), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to main database: %v", err)
	} else {
		fmt.Println("✅ Main database connection successful!")

		// Test query
		var result int
		mainDB.Raw("SELECT 1").Scan(&result)
		fmt.Printf("✅ Main database query successful! Result: %d\n", result)
	}

	// Test test database connection
	testDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("TEST_DB_HOST", "localhost"),
		getEnv("TEST_DB_USER", "postgres"),
		getEnv("TEST_DB_PASSWORD", "1234"),
		getEnv("TEST_DB_NAME", "go_taskmanagement_test"),
		getEnv("TEST_DB_PORT", "5432"),
		getEnv("TEST_DB_SSLMODE", "disable"),
	)

	fmt.Println("\nTesting test database connection...")
	fmt.Printf("DSN: %s\n", testDSN)

	testDB, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to test database: %v", err)
	} else {
		fmt.Println("✅ Test database connection successful!")

		// Test query
		var result int
		testDB.Raw("SELECT 1").Scan(&result)
		fmt.Printf("✅ Test database query successful! Result: %d\n", result)
	}

	fmt.Println("\n🎯 Database connectivity test completed!")
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
