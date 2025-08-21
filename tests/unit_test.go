package tests

import (
	"bytes"
	"encoding/json"
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

func TestHandlersSchemaDriven(t *testing.T) {
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

	// Initialize in-memory state
	models.Users = []models.User{}
	models.Tasks = []models.Task{}

	// Handler registry via operationId
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

	// Set up router and register static paths before parameterized to avoid conflicts
	router := mux.NewRouter()
	// collect static and param paths
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
	// register static routes first
	for _, path := range staticPaths {
		for mRaw, opRaw := range swagger.Paths[path] {
			method := strings.ToUpper(mRaw)
			var opMap map[string]interface{}
			json.Unmarshal(opRaw, &opMap)
			handler, ok := handlerRegistry[opMap["operationId"].(string)]
			if !ok {
				continue
			}
			if sec, ok := opMap["security"].([]interface{}); ok && len(sec) > 0 {
				handler = middleware.AuthMiddleware(handler)
			}
			router.HandleFunc(path, handler).Methods(method)
		}
	}
	// then parameterized routes
	for _, path := range paramPaths {
		for mRaw, opRaw := range swagger.Paths[path] {
			method := strings.ToUpper(mRaw)
			var opMap map[string]interface{}
			json.Unmarshal(opRaw, &opMap)
			handler, ok := handlerRegistry[opMap["operationId"].(string)]
			if !ok {
				continue
			}
			if sec, ok := opMap["security"].([]interface{}); ok && len(sec) > 0 {
				handler = middleware.AuthMiddleware(handler)
			}
			router.HandleFunc(path, handler).Methods(method)
		}
	}

	// Execute each operation against schema
	for path, ops := range swagger.Paths {
		for mRaw, opRaw := range ops {
			method := strings.ToUpper(mRaw)
			var opMap map[string]interface{}
			json.Unmarshal(opRaw, &opMap)
			opID, _ := opMap["operationId"].(string)
			// Build request URL, replace {id}
			urlPath := strings.ReplaceAll(path, "{id}", "1")
			var body bytes.Buffer
			if params, ok := opMap["parameters"].([]interface{}); ok {
				for _, pi := range params {
					p, _ := pi.(map[string]interface{})
					if p["in"] == "body" {
						// empty JSON object
						body.WriteString(`{}`)
						break
					}
				}
			}
			req := httptest.NewRequest(method, urlPath, &body)
			if body.Len() > 0 {
				req.Header.Set("Content-Type", "application/json")
			}
			// Route through mux to apply middleware
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			// Collect expected status codes
			expected := []int{}
			if responses, ok := opMap["responses"].(map[string]interface{}); ok {
				for codeStr := range responses {
					if code, err := strconv.Atoi(codeStr); err == nil {
						expected = append(expected, code)
					}
				}
			}
			// Allow unauthorized where security defined
			if sec, ok := opMap["security"].([]interface{}); ok && len(sec) > 0 {
				expected = append(expected, http.StatusUnauthorized)
			}

			// Assert
			status := rr.Code
			found := false
			for _, c := range expected {
				if status == c {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s %s -> %s: expected one of %v, got %d", mRaw, path, opID, expected, status)
			}
		}
	}
}
