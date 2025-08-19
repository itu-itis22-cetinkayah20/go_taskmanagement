package main

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	_ "go_taskmanagement/docs" // Swagger dokümantasyonu için
	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
)

func main() {
	r := mux.NewRouter()
	// Swagger UI endpoints
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	// Public endpoints
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/tasks/public", handlers.PublicTasksHandler).Methods("GET")
	// Private endpoints (JWT auth middleware eklenecek)
	r.HandleFunc("/tasks", middleware.AuthMiddleware(handlers.TasksListHandler)).Methods("GET")
	r.HandleFunc("/tasks", middleware.AuthMiddleware(handlers.TaskCreateHandler)).Methods("POST")
	r.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(handlers.TaskDetailHandler)).Methods("GET")
	r.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(handlers.TaskUpdateHandler)).Methods("PUT")
	r.HandleFunc("/tasks/{id}", middleware.AuthMiddleware(handlers.TaskDeleteHandler)).Methods("DELETE")
	r.HandleFunc("/logout", middleware.AuthMiddleware(handlers.LogoutHandler)).Methods("POST")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
