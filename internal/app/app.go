package app

import (
	"fmt"
	"log"
	"os"

	"go_taskmanagement/database"
	"go_taskmanagement/handlers"
	"go_taskmanagement/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewApp creates and configures a new Fiber application
func NewApp() *fiber.App {
	// Load environment variables for testing
	godotenv.Load()

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

// NewTestApp creates a new Fiber application with test database
func NewTestApp() *fiber.App {
	// Load environment variables
	godotenv.Load()

	// Try to connect to test database
	if isTestDatabaseAvailable() {
		// Connect to test database
		database.ConnectTest()
		database.Migrate()
		database.CleanTestData() // Clean before each test
		database.SeedTestData()
		log.Println("Test app initialized with PostgreSQL database")
	} else {
		log.Println("Test app initialized without database (in-memory mode)")
	}

	return NewApp()
}

// isTestDatabaseAvailable checks if test database is available
func isTestDatabaseAvailable() bool {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		getEnv("TEST_DB_HOST", "localhost"),
		getEnv("TEST_DB_USER", "postgres"),
		getEnv("TEST_DB_PASSWORD", "1234"),
		getEnv("TEST_DB_NAME", "go_taskmanagement_test"),
		getEnv("TEST_DB_PORT", "5432"),
		getEnv("TEST_DB_SSLMODE", "disable"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return false
	}

	var result int
	return db.Raw("SELECT 1").Scan(&result).Error == nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
