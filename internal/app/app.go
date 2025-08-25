package app

import (
	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// NewApp creates and configures a new Fiber application
func NewApp() *fiber.App {
	// Create Fiber app with custom config
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	// Public routes
	app.Post("/register", handlers.RegisterHandler)
	app.Post("/login", handlers.LoginHandler)
	app.Get("/tasks/public", handlers.PublicTasksHandler)

	// Protected routes with JWT middleware
	protected := app.Group("/", middleware.AuthMiddleware)
	protected.Get("/tasks", handlers.TasksListHandler)
	protected.Post("/tasks", handlers.TaskCreateHandler)
	protected.Get("/tasks/:id", handlers.TaskDetailHandler)
	protected.Put("/tasks/:id", handlers.TaskUpdateHandler)
	protected.Delete("/tasks/:id", handlers.TaskDeleteHandler)
	protected.Post("/logout", handlers.LogoutHandler)

	return app
}
