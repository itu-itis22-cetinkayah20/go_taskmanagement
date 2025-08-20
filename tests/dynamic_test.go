package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"

	"go_taskmanagement/models"

	"github.com/gorilla/mux"
)

func TestDynamicEndpoints(t *testing.T) {
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

	// Reset data models
	models.Users = []models.User{}
	models.Tasks = []models.Task{}
	// Start test server
	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/tasks/public", handlers.PublicTasksHandler).Methods("GET")
	router.HandleFunc("/tasks", middleware.AuthMiddleware(handlers.TasksListHandler)).Methods("GET")
	router.HandleFunc("/tasks", middleware.AuthMiddleware(handlers.TaskCreateHandler)).Methods("POST")
	router.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(handlers.TaskDetailHandler)).Methods("GET")
	router.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(handlers.TaskUpdateHandler)).Methods("PUT")
	router.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(handlers.TaskDeleteHandler)).Methods("DELETE")
	router.HandleFunc("/logout", middleware.AuthMiddleware(handlers.LogoutHandler)).Methods("POST")
	srv := httptest.NewServer(router)
	defer srv.Close()

	for path, ops := range swagger.Paths {
		for method, rawOp := range ops {
			// parse operation object
			var op struct {
				Security    []map[string][]string      `json:"security"`
				RequestBody json.RawMessage            `json:"requestBody"`
				Responses   map[string]json.RawMessage `json:"responses"`
			}
			if err := json.Unmarshal(rawOp, &op); err != nil {
				t.Fatalf("invalid operation for %s %s: %v", method, path, err)
			}
			t.Run(fmt.Sprintf("%s %s", method, path), func(t *testing.T) {
				// convert method to uppercase for HTTP request
				m := strings.ToUpper(method)
				// replace path parameters with dummy values
				urlPath := strings.ReplaceAll(path, "{id}", "1")
				var req *http.Request
				url := srv.URL + urlPath
				var err error
				// send dummy JSON body for POST and PUT
				if m == "POST" || m == "PUT" {
					req, err = http.NewRequest(m, url, bytes.NewBuffer([]byte("{}")))
					req.Header.Set("Content-Type", "application/json")
				} else {
					req, err = http.NewRequest(m, url, nil)
				}
				if err != nil {
					t.Fatalf("failed to create request: %v", err)
				}
				if len(op.Security) > 0 {
					req.Header.Set("Authorization", "Bearer invalid")
				}
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatalf("request failed: %v", err)
				}
				defer resp.Body.Close()
				expected := []int{}
				for codeStr := range op.Responses {
					if code, err := strconv.Atoi(codeStr); err == nil {
						expected = append(expected, code)
					}
				}
				// secured endpoints should allow 401
				if len(op.Security) > 0 {
					expected = append(expected, http.StatusUnauthorized)
				}
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
			})
		}
	}
}
