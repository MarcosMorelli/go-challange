package main

import (
	"log"

	"jobsity-backend/internal/config"
	"jobsity-backend/internal/database"
	"jobsity-backend/internal/handlers"
	"jobsity-backend/internal/repository"
	"jobsity-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	client, collection, err := database.ConnectMongo(database.Config{
		URI:      cfg.Database.URI,
		Database: cfg.Database.Database,
	})
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer database.DisconnectMongo(client)

	// Initialize repository
	userRepo := repository.NewMongoUserRepository(collection)

	// Initialize service
	userService := service.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Jobsity Backend API v1.0.0",
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := app.Group("/api/v1")
	api.Post("/login", userHandler.Login)
	api.Post("/users", userHandler.CreateUser)
	api.Get("/users/:email", userHandler.GetUser)

	// Start server
	log.Printf("Server starting on :%s", cfg.Server.Port)
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
