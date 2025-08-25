package contract

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"go_taskmanagement/internal/app"
)

// TestRealAuthenticationFlow tests the complete auth flow with real user registration and login
func TestRealAuthenticationFlow(t *testing.T) {
	t.Skip("Skipping real auth flow test - requires persistent database")

	f := app.NewApp()

	// Step 1: Register a new user
	registerData := map[string]string{
		"username": "testuser",
		"email":    "testuser@example.com",
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
		"email":    "testuser@example.com",
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
