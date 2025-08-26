package main

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	swagger "github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	_ "go_taskmanagement/docs"
	"log"
	"os"

	"go_taskmanagement/database"
	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Connect to database
	database.Connect()
	database.Migrate()
	database.SeedTestData()

	app := fiber.New()

	// CORS middleware
	app.Use(cors.New())

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
