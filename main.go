package main

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	_ "go_taskmanagement/docs"
	"log"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"

	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"
)

func main() {
	app := fiber.New()

	// Swagger UI endpoints
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Public endpoints
	app.Post("/register", handlers.RegisterHandler)
	app.Post("/login", handlers.LoginHandler)
	app.Get("/tasks/public", handlers.PublicTasksHandler)

	// Private endpoints with JWT auth
	app.Get("/tasks", middleware.AuthMiddleware, handlers.TasksListHandler)
	app.Post("/tasks", middleware.AuthMiddleware, handlers.TaskCreateHandler)
	app.Get("/tasks/:id", middleware.AuthMiddleware, handlers.TaskDetailHandler)
	app.Put("/tasks/:id", middleware.AuthMiddleware, handlers.TaskUpdateHandler)
	app.Delete("/tasks/:id", middleware.AuthMiddleware, handlers.TaskDeleteHandler)
	app.Post("/logout", middleware.AuthMiddleware, handlers.LogoutHandler)

	log.Println("Server started on :8080")
	log.Fatal(app.Listen(":8080"))
}
