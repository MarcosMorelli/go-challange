package main

import (
	"log"

	"jobsity-backend/internal/config"
	"jobsity-backend/internal/database"
	"jobsity-backend/internal/handlers"
	"jobsity-backend/internal/middleware"
	"jobsity-backend/internal/repository"
	"jobsity-backend/internal/service"
	"jobsity-backend/internal/websocket"

	fiberws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to MongoDB
	client, err := database.ConnectMongo(database.Config{
		URI:      cfg.Database.URI,
		Database: cfg.Database.Database,
	})
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer database.DisconnectMongo(client)

	// Get database instance
	db := client.Database(cfg.Database.Database)

	// Initialize repositories
	userRepo := repository.NewMongoUserRepository(db.Collection("users"))
	channelRepo := repository.NewMongoChannelRepository(db.Collection("channels"))
	messageRepo := repository.NewMongoMessageRepository(db.Collection("messages"))

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Initialize WebSocket handler
	wsHandler := websocket.NewHandler(wsHub)

	// Initialize services
	userService := service.NewUserService(userRepo)
	channelService := service.NewChannelService(channelRepo)
	baseMessageService := service.NewMessageService(messageRepo, channelRepo)
	messageService := service.NewWebSocketMessageService(baseMessageService, wsHandler)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	channelHandler := handlers.NewChannelHandler(channelService)
	messageHandler := handlers.NewMessageHandler(messageService)

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

	// User routes
	api.Post("/login", userHandler.Login)
	api.Post("/users", userHandler.CreateUser)
	api.Get("/users/:email", userHandler.GetUser)

	// Channel routes
	api.Post("/channels", middleware.AuthMiddleware(), channelHandler.CreateChannel)
	api.Get("/channels", channelHandler.GetAllChannels)
	api.Get("/channels/:id", channelHandler.GetChannel)
	api.Get("/channels/name/:name", channelHandler.GetChannelByName)
	api.Put("/channels/:id", middleware.AuthMiddleware(), channelHandler.UpdateChannel)
	api.Delete("/channels/:id", middleware.AuthMiddleware(), channelHandler.DeleteChannel)

	// Message routes
	api.Post("/messages", middleware.AuthMiddleware(), messageHandler.CreateMessage)
	api.Get("/messages/:id", messageHandler.GetMessage)
	api.Get("/channels/:channelId/messages", messageHandler.GetMessagesByChannel)
	api.Get("/channels/:channelId/messages/after", messageHandler.GetMessagesByChannelAfter)
	api.Put("/messages/:id", middleware.AuthMiddleware(), messageHandler.UpdateMessage)
	api.Delete("/messages/:id", middleware.AuthMiddleware(), messageHandler.DeleteMessage)

	// WebSocket routes
	api.Get("/ws", middleware.AuthMiddleware(), fiberws.New(wsHandler.HandleWebSocket))
	api.Get("/ws/stats", wsHandler.GetStats())

	// Start server
	log.Printf("Server starting on :%s", cfg.Server.Port)
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
