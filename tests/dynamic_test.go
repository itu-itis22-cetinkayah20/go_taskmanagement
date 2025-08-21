package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	// Load OpenAPI schema for response validation
	// Start test server: register routes by operationId from schema
	router := mux.NewRouter()
	// handler registry maps operationId to functions
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
	for path, ops := range swagger.Paths {
		for mRaw, opObj := range ops {
			method := strings.ToUpper(mRaw)
			// parse operation object into map
			var opRawMap map[string]interface{}
			json.Unmarshal(opObj, &opRawMap)
			// fetch operationId and lookup handler
			opID, _ := opRawMap["operationId"].(string)
			handler, exists := handlerRegistry[opID]
			if !exists {
				continue
			}
			// wrap auth if security defined
			if sec, ok := opRawMap["security"].([]interface{}); ok && len(sec) > 0 {
				handler = middleware.AuthMiddleware(handler)
			}
			router.HandleFunc(path, handler).Methods(method)
		}
	}
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
			// parse operation object into map for parameters
			var opRawMap map[string]interface{}
			if err := json.Unmarshal(rawOp, &op); err != nil {
				t.Fatalf("invalid operation for %s %s: %v", method, path, err)
			}
			if err := json.Unmarshal(rawOp, &opRawMap); err != nil {
				t.Fatalf("invalid operation map for %s %s: %v", method, path, err)
			}
			t.Run(fmt.Sprintf("%s %s", method, path), func(t *testing.T) {
				// convert method to uppercase for HTTP request
				m := strings.ToUpper(method)
				// replace path parameters with dummy values
				urlPath := strings.ReplaceAll(path, "{id}", "1")
				var req *http.Request
				url := srv.URL + urlPath
				var err error
				// prepare request body from schema parameters
				var bodyBytes []byte
				if params, ok := opRawMap["parameters"].([]interface{}); ok {
					// find body parameter
					for _, pi := range params {
						if p, ok := pi.(map[string]interface{}); ok && p["in"] == "body" {
							if schema, ok := p["schema"].(map[string]interface{}); ok {
								if props, ok := schema["properties"].(map[string]interface{}); ok {
									payload := map[string]interface{}{}
									for name, def := range props {
										if d, ok := def.(map[string]interface{}); ok {
											switch d["type"] {
											case "string":
												payload[name] = fmt.Sprintf("test_%s", name)
											case "integer":
												payload[name] = 1
											case "array":
												payload[name] = []interface{}{}
											default:
												payload[name] = nil
											}
										}
									}
									if b, err2 := json.Marshal(payload); err2 == nil {
										bodyBytes = b
									}
								}
							}
							break
						}
					}
				}
				if m == "POST" || m == "PUT" {
					if len(bodyBytes) == 0 {
						bodyBytes = []byte("{}")
					}
					req, err = http.NewRequest(m, url, bytes.NewReader(bodyBytes))
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
				// read full body for status and JSON validation
				bodyBytes, err = io.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
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
				// Basic JSON validity check for application/json responses
				if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
					var js interface{}
					if err := json.Unmarshal(bodyBytes, &js); err != nil {
						t.Errorf("%s %s: invalid JSON response: %v", method, path, err)
					}
				}
			})
		}
	}
}
