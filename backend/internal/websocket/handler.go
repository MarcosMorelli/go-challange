package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

// Handler handles WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// HandleWebSocket handles WebSocket connections
func (h *Handler) HandleWebSocket(c *websocket.Conn) {
	// Get user email from context (set by auth middleware)
	userEmail := c.Locals("userEmail")
	if userEmail == nil {
		log.Printf("WebSocket connection rejected: user not authenticated")
		return
	}

	// Get channel ID from query parameter (optional)
	channelID := c.Query("channel_id", "")

	// Create new client
	client := &Client{
		conn:      c,
		send:      make(chan []byte, 256),
		hub:       h.hub,
		UserEmail: userEmail.(string),
		ChannelID: channelID,
	}

	// Register client with hub
	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// BroadcastMessage broadcasts a message to all clients in a channel
func (h *Handler) BroadcastMessage(channelID string, messageType string, data interface{}) {
	message := Message{
		Type:      messageType,
		ChannelID: channelID,
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Data:      data,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	h.hub.BroadcastToChannel(channelID, messageBytes)
}

// BroadcastToAll broadcasts a message to all connected clients
func (h *Handler) BroadcastToAll(messageType string, data interface{}) {
	message := Message{
		Type:      messageType,
		Timestamp: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Data:      data,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message: %v", err)
		return
	}

	select {
	case h.hub.broadcast <- messageBytes:
	default:
		log.Printf("Broadcast channel is full, dropping message")
	}
}

// GetStats returns WebSocket connection statistics
func (h *Handler) GetStats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		stats := fiber.Map{
			"total_clients": h.hub.GetClientCount(),
			"channels":      make(map[string]int),
		}

		// Get client count per channel
		h.hub.mutex.RLock()
		for channelID := range h.hub.channelClients {
			stats["channels"].(map[string]int)[channelID] = h.hub.GetChannelClientCount(channelID)
		}
		h.hub.mutex.RUnlock()

		return c.JSON(fiber.Map{
			"success": true,
			"data":    stats,
		})
	}
}
