package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"
	"go_taskmanagement/models"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Test PublicTasksHandler
func TestPublicTasksHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/tasks/public", nil)
	rr := httptest.NewRecorder()

	handlers.PublicTasksHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var got []models.Task
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}
	if len(got) != len(models.PublicTasks) {
		t.Errorf("expected %d tasks, got %d", len(models.PublicTasks), len(got))
	}
}

// Test TasksListHandler
func TestTasksListHandler(t *testing.T) {
	// reset tasks
	models.Tasks = []models.Task{
		{ID: 1, UserID: 1, Title: "A", Details: ""},
		{ID: 2, UserID: 2, Title: "B", Details: ""},
	}
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handlers.TasksListHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var got []models.Task
	json.Unmarshal(rr.Body.Bytes(), &got)
	if len(got) != 1 || got[0].UserID != 1 {
		t.Errorf("unexpected tasks: %v", got)
	}
}

// Test TaskCreateHandler
func TestTaskCreateHandler(t *testing.T) {
	models.Tasks = []models.Task{}
	newTask := models.Task{Title: "New", Details: "D"}
	body, _ := json.Marshal(newTask)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handlers.TaskCreateHandler(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	var got models.Task
	json.Unmarshal(rr.Body.Bytes(), &got)
	if got.ID != 1 || got.UserID != 1 || got.Title != newTask.Title {
		t.Errorf("unexpected task: %v", got)
	}
}

// Test TaskDetailHandler
func TestTaskDetailHandler(t *testing.T) {
	models.Tasks = []models.Task{
		{ID: 1, UserID: 1, Title: "A", Details: "1"},
	}
	req := httptest.NewRequest(http.MethodGet, "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handlers.TaskDetailHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var got models.Task
	json.Unmarshal(rr.Body.Bytes(), &got)
	if got.ID != 1 {
		t.Errorf("unexpected task: %v", got)
	}
}

// Test TaskUpdateHandler
func TestTaskUpdateHandler(t *testing.T) {
	models.Tasks = []models.Task{
		{ID: 1, UserID: 1, Title: "Old", Details: "OldD"},
	}
	updated := models.Task{Title: "New", Details: "NewD"}
	body, _ := json.Marshal(updated)
	req := httptest.NewRequest(http.MethodPut, "/tasks/1", bytes.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handlers.TaskUpdateHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var got models.Task
	json.Unmarshal(rr.Body.Bytes(), &got)
	if got.Title != "New" || got.Details != "NewD" {
		t.Errorf("unexpected task: %v", got)
	}
}

// Test TaskDeleteHandler
func TestTaskDeleteHandler(t *testing.T) {
	models.Tasks = []models.Task{
		{ID: 1, UserID: 1, Title: "A", Details: ""},
		{ID: 2, UserID: 2, Title: "B", Details: ""},
	}
	req := httptest.NewRequest(http.MethodDelete, "/tasks/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handlers.TaskDeleteHandler(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, rr.Code)
	}
	if len(models.Tasks) != 1 || models.Tasks[0].ID == 1 {
		t.Errorf("task not deleted: %v", models.Tasks)
	}
}

// Test RegisterHandler
func TestRegisterHandler(t *testing.T) {
	models.Users = []models.User{}
	user := models.User{Username: "u", Password: "p"}
	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handlers.RegisterHandler(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}
	var got models.User
	json.Unmarshal(rr.Body.Bytes(), &got)
	if got.ID != 1 || got.Username != "u" {
		t.Errorf("unexpected user: %v", got)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte("p")); err != nil {
		t.Errorf("password not hashed correctly")
	}
}

// Test LoginHandler
func TestLoginHandler(t *testing.T) {
	password := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	models.Users = []models.User{{ID: 1, Username: "user", Password: string(hash)}}

	creds := models.User{Username: "user", Password: password}
	body, _ := json.Marshal(creds)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handlers.LoginHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if _, ok := resp["token"]; !ok {
		t.Errorf("token not returned")
	}
}

// Test LogoutHandler
func TestLogoutHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/logout", nil)
	rr := httptest.NewRecorder()

	handlers.LogoutHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}
