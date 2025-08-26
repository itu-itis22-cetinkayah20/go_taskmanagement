package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"go_taskmanagement/internal/app"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// isDatabaseAvailable checks if PostgreSQL is available and configured
func isDatabaseAvailable() bool {
	// Load environment variables
	godotenv.Load()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("TEST_DB_HOST", "localhost"),
		getEnv("TEST_DB_USER", "postgres"),
		getEnv("TEST_DB_PASSWORD", "1234"),
		getEnv("TEST_DB_NAME", "go_taskmanagement_test"),
		getEnv("TEST_DB_PORT", "5432"),
		getEnv("TEST_DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return false
	}

	// Test a simple query
	var result int
	if err := db.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return false
	}

	return true
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// TestRealAuthenticationFlow tests the complete auth flow with real user registration and login
func TestRealAuthenticationFlow(t *testing.T) {
	// Check if PostgreSQL is available
	if !isDatabaseAvailable() {
		t.Skip("Skipping real auth flow test - PostgreSQL database not available or not configured")
	}

	// Use test database
	f := app.NewTestApp()

	// Generate unique test user to avoid conflicts
	timestamp := time.Now().UnixNano()
	username := fmt.Sprintf("testuser_%d", timestamp)
	email := fmt.Sprintf("testuser_%d@example.com", timestamp)

	// Step 1: Register a new user
	registerData := map[string]string{
		"username": username,
		"email":    email,
		"password": "password123",
	}
	registerBody, _ := json.Marshal(registerData)

	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(registerBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.Test(req, 10000)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Register failed with status %d, body: %s", resp.StatusCode, string(body))
	}
	t.Logf("✓ User registered successfully")

	// Step 2: Login with the registered user
	loginData := map[string]string{
		"email":    email,
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginData)

	req, _ = http.NewRequest("POST", "/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err = f.Test(req, 10000)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Login failed with status %d, body: %s", resp.StatusCode, string(body))
	} // Extract the token from login response
	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		} `json:"user"`
	}

	json.NewDecoder(resp.Body).Decode(&loginResp)

	if loginResp.Token == "" {
		t.Fatal("No token received from login")
	}
	t.Logf("✓ Login successful, received token")

	// Step 3: Test protected endpoints with real token
	req, _ = http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	resp, err = f.Test(req, 10000)
	if err != nil {
		t.Fatalf("Get tasks failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Get tasks failed with status %d", resp.StatusCode)
	}
	t.Logf("✓ Protected endpoint accessible with valid token")

	// Step 4: Test creating a task
	taskData := map[string]interface{}{
		"title":       "Test Task",
		"description": "Test Description",
		"status":      "pending",
		"priority":    "medium",
	}
	taskBody, _ := json.Marshal(taskData)

	req, _ = http.NewRequest("POST", "/tasks", bytes.NewReader(taskBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	resp, err = f.Test(req, 10000)
	if err != nil {
		t.Fatalf("Create task failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Create task failed with status %d, body: %s", resp.StatusCode, string(body))
	}
	t.Logf("✓ Task created successfully")

	// Step 5: Test logout
	req, _ = http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", "Bearer "+loginResp.Token)

	resp, err = f.Test(req, 10000)
	if err != nil {
		t.Fatalf("Logout failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Logout failed with status %d", resp.StatusCode)
	}
	t.Logf("✓ Logout successful")
}
