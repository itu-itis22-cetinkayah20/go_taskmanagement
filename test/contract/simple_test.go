package contract

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"go_taskmanagement/internal/app"
)

// Simple test to verify endpoints work without OpenAPI validation
func TestBasicEndpoints(t *testing.T) {
	f := app.NewApp()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
		authRequired   bool
	}{
		{"Register", "POST", "/register", `{"username":"test","email":"test@example.com","password":"password123"}`, 201, false},
		{"Login", "POST", "/login", `{"email":"test@example.com","password":"password123"}`, 401, false}, // 401 because user doesn't exist yet
		{"Public Tasks", "GET", "/tasks/public", "", 200, false},
		{"Private Tasks", "GET", "/tasks", "", 401, true},
		{"Create Task", "POST", "/tasks", `{"title":"Test Task","description":"Test Description"}`, 401, true},
		{"Get Task", "GET", "/tasks/1", "", 401, true},
		{"Update Task", "PUT", "/tasks/1", `{"title":"Updated Task"}`, 401, true},
		{"Delete Task", "DELETE", "/tasks/1", "", 401, true},
		{"Logout", "POST", "/logout", "", 401, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req, _ = http.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
				req.Header.Set("Content-Type", "application/json")
				t.Logf("Request body: %s", tt.body)
			} else {
				req, _ = http.NewRequest(tt.method, tt.path, nil)
			}

			// Add auth header for protected endpoints
			if tt.authRequired {
				req.Header.Set("Authorization", "Bearer invalid-token")
				t.Logf("Added auth header: Bearer invalid-token")
			}

			resp, err := f.Test(req, 10000)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			t.Logf("✓ %s %s -> %d", tt.method, tt.path, resp.StatusCode)
		})
	}
}

// Test OpenAPI spec loading separately
func TestOpenAPISpecLoading(t *testing.T) {
	// Test that spec file can be found and loaded
	err := ValidateSpecFile()
	if err != nil {
		t.Fatalf("Spec file validation failed: %v", err)
	}

	// Test that spec can be parsed
	ctx := context.Background()
	doc, err := LoadSpec(ctx, SpecPath())
	if err != nil {
		t.Fatalf("Spec loading failed: %v", err)
	}

	if doc == nil {
		t.Fatal("Loaded spec is nil")
	}

	// Check that we have the expected paths
	expectedPaths := []string{"/register", "/login", "/logout", "/tasks", "/tasks/public", "/tasks/{id}"}
	for _, path := range expectedPaths {
		if doc.Paths.Find(path) == nil {
			t.Errorf("Expected path %s not found in spec", path)
		}
	}

	t.Logf("✓ OpenAPI spec loaded successfully with %d paths", len(doc.Paths.Map()))
}
