package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"
	"go_taskmanagement/models"

	"github.com/gorilla/mux"
)

// TestSchemaDrivenHandlers dynamically tests all endpoints defined in swagger.json
func TestSchemaDrivenHandlers(t *testing.T) {
	// Load Swagger JSON
	data, err := os.ReadFile("../docs/swagger.json")
	if err != nil {
		t.Fatalf("failed to read swagger.json: %v", err)
	}
	var swagger struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(data, &swagger); err != nil {
		t.Fatalf("invalid swagger JSON: %v", err)
	}

	// Reset application state
	models.Users = []models.User{}
	models.Tasks = []models.Task{}

	// Register handlers by operationId
	handlerRegistry := map[string]http.HandlerFunc{
		"RegisterHandler":    handlers.RegisterHandler,
		"LoginHandler":       handlers.LoginHandler,
		"PublicTasksHandler": handlers.PublicTasksHandler,
		"TasksListHandler":   handlers.TasksListHandler,
		"TaskCreateHandler":  handlers.TaskCreateHandler,
		"TaskDetailHandler":  handlers.TaskDetailHandler,
		"TaskUpdateHandler":  handlers.TaskUpdateHandler,
		"TaskDeleteHandler":  handlers.TaskDeleteHandler,
		"LogoutHandler":      handlers.LogoutHandler,
	}

	// Build router and register routes: static first, then parameterized
	router := mux.NewRouter()
	staticPaths := []string{}
	paramPaths := []string{}
	for p := range swagger.Paths {
		if strings.Contains(p, "{") {
			paramPaths = append(paramPaths, p)
		} else {
			staticPaths = append(staticPaths, p)
		}
	}
	sort.Strings(staticPaths)
	sort.Strings(paramPaths)
	// register static routes
	for _, path := range staticPaths {
		for mRaw, rawOp := range swagger.Paths[path] {
			method := strings.ToUpper(mRaw)
			var opMap map[string]interface{}
			json.Unmarshal(rawOp, &opMap)
			opID, _ := opMap["operationId"].(string)
			handler, ok := handlerRegistry[opID]
			if !ok {
				continue
			}
			if sec, ok := opMap["security"].([]interface{}); ok && len(sec) > 0 {
				handler = middleware.AuthMiddleware(handler)
			}
			router.HandleFunc(path, handler).Methods(method)
		}
	}
	// register parameterized routes
	for _, path := range paramPaths {
		for mRaw, rawOp := range swagger.Paths[path] {
			method := strings.ToUpper(mRaw)
			var opMap map[string]interface{}
			json.Unmarshal(rawOp, &opMap)
			opID, _ := opMap["operationId"].(string)
			handler, ok := handlerRegistry[opID]
			if !ok {
				continue
			}
			if sec, ok := opMap["security"].([]interface{}); ok && len(sec) > 0 {
				handler = middleware.AuthMiddleware(handler)
			}
			router.HandleFunc(path, handler).Methods(method)
		}
	}

	// Iterate and test each endpoint
	for path, ops := range swagger.Paths {
		for mRaw, rawOp := range ops {
			method := strings.ToUpper(mRaw)
			var opMap map[string]interface{}
			json.Unmarshal(rawOp, &opMap)
			t.Run(fmt.Sprintf("%s %s", method, path), func(t *testing.T) {
				// Prepare request URL
				urlPath := strings.ReplaceAll(path, "{id}", "1")
				// Build request body if defined
				var bodyBytes []byte
				if params, ok := opMap["parameters"].([]interface{}); ok {
					for _, pi := range params {
						p, _ := pi.(map[string]interface{})
						if p["in"] == "body" {
							// default empty JSON object
							bodyBytes = []byte(`{}`)
							break
						}
					}
				}
				// Create request
				var req *http.Request
				if method == "POST" || method == "PUT" {
					req = httptest.NewRequest(method, urlPath, bytes.NewReader(bodyBytes))
					req.Header.Set("Content-Type", "application/json")
				} else {
					req = httptest.NewRequest(method, urlPath, nil)
				}
				// Add invalid auth for secured endpoints to test 401
				if sec, ok := opMap["security"].([]interface{}); ok && len(sec) > 0 {
					req.Header.Set("Authorization", "Bearer invalid")
				}
				// Execute
				rr := httptest.NewRecorder()
				router.ServeHTTP(rr, req)
				resp := rr.Result()
				defer resp.Body.Close()

				// Read body
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response: %v", err)
				}
				// Determine expected status codes
				expected := []int{}
				if responses, ok := opMap["responses"].(map[string]interface{}); ok {
					for codeStr := range responses {
						if code, err := strconv.Atoi(codeStr); err == nil {
							expected = append(expected, code)
						}
					}
				}
				// allow 401 for secured
				if sec, ok := opMap["security"].([]interface{}); ok && len(sec) > 0 {
					expected = append(expected, http.StatusUnauthorized)
				}
				// Assert status
				found := false
				for _, c := range expected {
					if resp.StatusCode == c {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s %s: expected one of %v, got %d", method, path, expected, resp.StatusCode)
				}
				// Basic JSON validation
				if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
					var js interface{}
					if err := json.Unmarshal(respBody, &js); err != nil {
						t.Errorf("%s %s: invalid JSON: %v", method, path, err)
					}
				}
			})
		}
	}
}
