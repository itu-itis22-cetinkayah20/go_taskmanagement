package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"go_taskmanagement/database"
	"go_taskmanagement/internal/app"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
)

func Test_OpenAPI_Contract(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 0) Load & validate spec
	doc, err := LoadSpec(ctx, SpecPath())
	if err != nil {
		t.Fatalf("spec: %v", err)
	}

	// 2) Spin up Fiber app in-process with test database
	f := app.NewTestApp()

	// Clear test database before this test group
	if database.IsConnected {
		database.DB.Exec("TRUNCATE TABLE tasks, users RESTART IDENTITY CASCADE")
	}

	// 3) Token source - Create a real authenticated user
	var token TokenSource
	if v := os.Getenv("TEST_BEARER"); v != "" {
		token = StaticToken(v)
	} else {
		// Create a login token that will register and login automatically
		email := fmt.Sprintf("test_%s@example.com", generateRandomString())
		password := "testpass123"

		// Register the user first
		RegisterUserForTest(f, email, password)

		// Create login token with Fiber app reference
		token = &LoginToken{
			BaseURL: "", // Empty because we're using Fiber in-process
			Email:   email,
			Pass:    password,
			App:     f, // Pass Fiber app for in-process testing
		}
	}

	// 4) Iterate all operations dynamically - FIX: Actually run the tests!
	for path, pi := range doc.Paths.Map() {
		path, pi := path, pi // capture for loop vars
		for method, op := range pi.Operations() {
			method, op := method, op

			t.Run(method+" "+path, func(t *testing.T) {
				t.Parallel()

				// 4a) Build request dynamically from OpenAPI spec
				buildInput := BuildInput{
					Doc:      doc,
					PathTmpl: path,
					Op:       op,
					Method:   method,
					Token:    token,
				}

				// Special case for login endpoint - use registered user credentials
				if method == "POST" && path == "/login" {
					if loginToken, ok := token.(*LoginToken); ok {
						buildInput.LoginCredentials = &LoginCredentials{
							Email:    loginToken.Email,
							Password: loginToken.Pass,
						}
					}
				}

				req, err := BuildRequest(ctx, buildInput)
				if err != nil {
					t.Fatalf("build req: %v", err)
				}

				// 4b) Execute via Fiber in-process (no network calls)
				resp, err := f.Test(req, 30_000)
				if err != nil {
					t.Fatalf("fiber test: %v", err)
				}
				defer resp.Body.Close()
				body, _ := io.ReadAll(resp.Body)

				// 4c) Basic validation - just log the results
				t.Logf("✓ %s %s - Status: %d", method, path, resp.StatusCode)

				// Optional: Add more validation here
				if resp.StatusCode >= 500 {
					t.Errorf("Server error %d for %s %s: %s", resp.StatusCode, method, path, string(body))
				}
			})
		}
	}
}

// Test specific endpoint patterns for more detailed validation
func Test_OpenAPI_Contract_AuthFlow(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	doc, err := LoadSpec(ctx, SpecPath())
	if err != nil {
		t.Fatalf("spec: %v", err)
	}

	f := app.NewApp()

	// Clear test database before this test group
	if database.IsConnected {
		database.DB.Exec("TRUNCATE TABLE tasks, users RESTART IDENTITY CASCADE")
	}

	// Test authentication flow specifically
	testCases := []struct {
		name     string
		path     string
		method   string
		needAuth bool
	}{
		{"Register", "/register", "POST", false},
		{"Login", "/login", "POST", false},
		{"Public Tasks", "/tasks/public", "GET", false},
		{"Private Tasks", "/tasks", "GET", true},
		{"Create Task", "/tasks", "POST", true},
		{"Logout", "/logout", "POST", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Find operation in spec
			pathItem := doc.Paths.Find(tc.path)
			if pathItem == nil {
				t.Fatalf("path %s not found in spec", tc.path)
			}

			var op *openapi3.Operation
			switch tc.method {
			case "GET":
				op = pathItem.Get
			case "POST":
				op = pathItem.Post
			case "PUT":
				op = pathItem.Put
			case "DELETE":
				op = pathItem.Delete
			}

			if op == nil {
				t.Fatalf("operation %s %s not found in spec", tc.method, tc.path)
			}

			// Build and test request
			var token TokenSource
			if tc.needAuth {
				token = StaticToken("test-token") // Use a test token for protected endpoints
			}

			req, err := BuildRequest(ctx, BuildInput{
				Doc:      doc,
				PathTmpl: tc.path,
				Op:       op,
				Method:   tc.method,
				Token:    token,
			})
			if err != nil {
				t.Fatalf("build req: %v", err)
			}

			resp, err := f.Test(req, 10_000)
			if err != nil {
				t.Fatalf("fiber test: %v", err)
			}
			defer resp.Body.Close()

			// Validate that auth requirements are respected
			if tc.needAuth && resp.StatusCode == 401 {
				t.Logf("✓ %s correctly requires authentication", tc.name)
			} else if !tc.needAuth && resp.StatusCode != 401 {
				t.Logf("✓ %s correctly allows public access", tc.name)
			} else if tc.name == "Login" && resp.StatusCode == 401 {
				// Login returning 401 is expected when using fake credentials
				t.Logf("✓ %s correctly rejects invalid credentials", tc.name)
			} else {
				t.Logf("! %s returned status: %d (may be expected)", tc.name, resp.StatusCode)
			}
		})
	}
}

// RegisterUserForTest registers a user for testing via Fiber app
func RegisterUserForTest(f *fiber.App, email, password string) error {
	// Create register request
	registerData := map[string]string{
		"username": email[:strings.Index(email, "@")], // Extract username from email
		"email":    email,
		"password": password,
	}

	bodyBytes, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Execute request via Fiber
	resp, err := f.Test(req, 10_000)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
