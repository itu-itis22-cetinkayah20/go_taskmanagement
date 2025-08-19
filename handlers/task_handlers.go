package handlers

import (
	"encoding/json"
	"go_taskmanagement/models"
	"net/http"
	"strconv"

	"go_taskmanagement/middleware"

	"github.com/gorilla/mux"
)

func PublicTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.PublicTasks)
}

func TasksListHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Kullanıcı bilgisi alınamadı"}`))
		return
	}
	var userTasks []models.Task
	for _, t := range models.Tasks {
		if t.UserID == userID {
			userTasks = append(userTasks, t)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userTasks)
}

func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Geçersiz veri"}`))
		return
	}
	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Başlık zorunlu"}`))
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Kullanıcı bilgisi alınamadı"}`))
		return
	}
	task.ID = len(models.Tasks) + 1
	task.UserID = userID
	models.Tasks = append(models.Tasks, task)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func TaskDetailHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Kullanıcı bilgisi alınamadı"}`))
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Geçersiz görev ID"}`))
		return
	}
	for _, t := range models.Tasks {
		if t.ID == id && t.UserID == userID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Görev bulunamadı veya yetkiniz yok"}`))
}

func TaskUpdateHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Kullanıcı bilgisi alınamadı"}`))
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Geçersiz görev ID"}`))
		return
	}
	var updated models.Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Geçersiz veri"}`))
		return
	}
	for i, t := range models.Tasks {
		if t.ID == id && t.UserID == userID {
			models.Tasks[i].Title = updated.Title
			models.Tasks[i].Details = updated.Details
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.Tasks[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Görev bulunamadı veya yetkiniz yok"}`))
}

func TaskDeleteHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Kullanıcı bilgisi alınamadı"}`))
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Geçersiz görev ID"}`))
		return
	}
	for i, t := range models.Tasks {
		if t.ID == id && t.UserID == userID {
			models.Tasks = append(models.Tasks[:i], models.Tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error":"Görev bulunamadı veya yetkiniz yok"}`))
}
