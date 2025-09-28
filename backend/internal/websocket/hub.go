package websocket

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Channel-specific clients
	channelClients map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mutex sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		channelClients: make(map[string]map[*Client]bool),
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true

			// Add to channel-specific clients if channel is specified
			if client.ChannelID != "" {
				if h.channelClients[client.ChannelID] == nil {
					h.channelClients[client.ChannelID] = make(map[*Client]bool)
				}
				h.channelClients[client.ChannelID][client] = true
			}
			h.mutex.Unlock()

			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove from channel-specific clients
				if client.ChannelID != "" {
					if channelClients, exists := h.channelClients[client.ChannelID]; exists {
						delete(channelClients, client)
						if len(channelClients) == 0 {
							delete(h.channelClients, client.ChannelID)
						}
					}
				}
			}
			h.mutex.Unlock()

			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// BroadcastToChannel broadcasts a message to all clients in a specific channel
func (h *Hub) BroadcastToChannel(channelID string, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if channelClients, exists := h.channelClients[channelID]; exists {
		for client := range channelClients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(channelClients, client)
			}
		}
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetChannelClientCount returns the number of clients in a specific channel
func (h *Hub) GetChannelClientCount(channelID string) int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if channelClients, exists := h.channelClients[channelID]; exists {
		return len(channelClients)
	}
	return 0
}
