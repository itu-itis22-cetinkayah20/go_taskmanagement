package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

var baseURL = "http://localhost:8080"
var testToken string
var testTaskID int

func TestAPI(t *testing.T) {
	// Benzersiz username ve task title üret
	uniqueUsername := fmt.Sprintf("testuser_%d", time.Now().UnixNano())
	uniqueTaskTitle := fmt.Sprintf("Test Görev %d", time.Now().UnixNano())

	// Register
	resp := post(t, "/register", map[string]string{"username": uniqueUsername, "password": "123456"}, "")
	if resp.StatusCode != 201 {
		t.Fatalf("Register failed: %s", respBody(resp))
	}
	fmt.Println("Register test passed.")

	// Login
	resp = post(t, "/login", map[string]string{"username": uniqueUsername, "password": "123456"}, "")
	if resp.StatusCode != 200 {
		t.Fatalf("Login failed: %s", respBody(resp))
	}
	fmt.Println("Login test passed.")
	var loginResp map[string]string
	json.Unmarshal([]byte(respBody(resp)), &loginResp)
	testToken = loginResp["token"]

	// Public tasks
	resp = get(t, "/tasks/public", "")
	if resp.StatusCode != 200 {
		t.Fatalf("Public tasks failed: %s", respBody(resp))
	}
	fmt.Println("Public tasks test passed.")

	// Create task
	task := map[string]string{"title": uniqueTaskTitle, "details": "Detay"}
	resp = post(t, "/tasks", task, testToken)
	if resp.StatusCode != 201 {
		t.Fatalf("Task create failed: %s", respBody(resp))
	}
	fmt.Println("Task create test passed.")
	var createdTask map[string]interface{}
	json.Unmarshal([]byte(respBody(resp)), &createdTask)
	testTaskID = int(createdTask["id"].(float64))

	// List tasks
	resp = get(t, "/tasks", testToken)
	if resp.StatusCode != 200 {
		t.Fatalf("Task list failed: %s", respBody(resp))
	}
	fmt.Println("Task list test passed.")

	// Task detail
	resp = get(t, fmt.Sprintf("/tasks/%d", testTaskID), testToken)
	if resp.StatusCode != 200 {
		t.Fatalf("Task detail failed: %s", respBody(resp))
	}
	fmt.Println("Task detail test passed.")

	// Update task
	update := map[string]string{"title": "Güncellendi", "details": "Yeni detay"}
	resp = put(t, fmt.Sprintf("/tasks/%d", testTaskID), update, testToken)
	if resp.StatusCode != 200 {
		t.Fatalf("Task update failed: %s", respBody(resp))
	}
	fmt.Println("Task update test passed.")

	// Delete task
	resp = deleteReq(t, fmt.Sprintf("/tasks/%d", testTaskID), testToken)
	if resp.StatusCode != 204 {
		t.Fatalf("Task delete failed: %s", respBody(resp))
	}
	fmt.Println("Task delete test passed.")

	// Logout
	resp = post(t, "/logout", nil, testToken)
	if resp.StatusCode != 200 {
		t.Fatalf("Logout failed: %s", respBody(resp))
	}
	fmt.Println("Logout test passed.")

	fmt.Println("Tüm API endpoint testleri başarıyla geçti!")
}

func get(t *testing.T, path, token string) *http.Response {
	req, _ := http.NewRequest("GET", baseURL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s error: %v", path, err)
	}
	return resp
}

func post(t *testing.T, path string, body interface{}, token string) *http.Response {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest("POST", baseURL+path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s error: %v", path, err)
	}
	return resp
}

func put(t *testing.T, path string, body interface{}, token string) *http.Response {
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(body)
	req, _ := http.NewRequest("PUT", baseURL+path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT %s error: %v", path, err)
	}
	return resp
}

func deleteReq(t *testing.T, path, token string) *http.Response {
	req, _ := http.NewRequest("DELETE", baseURL+path, nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE %s error: %v", path, err)
	}
	return resp
}

func respBody(resp *http.Response) string {
	b, _ := io.ReadAll(resp.Body)
	return string(b)
}
