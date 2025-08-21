package handlers

import (
	"encoding/json"
	"go_taskmanagement/middleware"
	"go_taskmanagement/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// PublicTasksHandler herkese açık görevleri listeler
// @ID PublicTasksHandler
// @Summary Public görevleri listele
// @Description Herkesin görebileceği görevleri döner
// @Tags Tasks
// @Produce json
// @Success 200 {array} models.Task
// @Router /tasks/public [get]
func PublicTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.PublicTasks)
}

// TasksListHandler kullanıcının kendi görevlerini listeler
// @ID TasksListHandler
// @Summary Kullanıcı görevlerini listele
// @Description Sadece giriş yapan kullanıcının görevlerini döner
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Task
// @Router /tasks [get]
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

// TaskCreateHandler yeni görev ekler
// @ID TaskCreateHandler
// @Summary Görev ekle
// @Description Yeni görev oluşturur
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body models.Task true "Görev"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
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

// TaskDetailHandler görev detayını döner
// @ID TaskDetailHandler
// @Summary Görev detayını görüntüle
// @Description Belirli bir görevin detayını döner
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Görev ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
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

// TaskUpdateHandler görevi günceller
// @ID TaskUpdateHandler
// @Summary Görev güncelle
// @Description Belirli bir görevi günceller
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Görev ID"
// @Param task body models.Task true "Görev"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [put]
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

// TaskDeleteHandler görevi siler
// @ID TaskDeleteHandler
// @Summary Görev sil
// @Description Belirli bir görevi siler
// @Tags Tasks
// @Security BearerAuth
// @Param id path int true "Görev ID"
// @Success 204 {string} string ""
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
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
