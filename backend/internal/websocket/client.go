package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// WebSocket upgrader configuration is handled by Fiber's websocket package

// Client represents a websocket client
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Hub for managing clients
	hub *Hub

	// User email
	UserEmail string

	// Channel ID this client is connected to
	ChannelID string
}

// Message represents a websocket message
type Message struct {
	Type      string      `json:"type"`
	ChannelID string      `json:"channel_id,omitempty"`
	UserEmail string      `json:"user_email,omitempty"`
	Content   string      `json:"content,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		log.Printf("Client readPump exiting for user: %s", c.UserEmail)
		c.hub.unregister <- c
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	log.Printf("Starting readPump for user: %s", c.UserEmail)

	// Check if connection is valid
	if c.conn == nil {
		log.Printf("WebSocket connection is nil for user: %s", c.UserEmail)
		return
	}

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Send a welcome message to the client
	welcomeMessage := Message{
		Type:      "connected",
		ChannelID: c.ChannelID,
		UserEmail: c.UserEmail,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	if welcomeBytes, err := json.Marshal(welcomeMessage); err == nil {
		c.send <- welcomeBytes
		log.Printf("Sent welcome message to user %s", c.UserEmail)
	}

	for {
		// Check if connection is still valid
		if c.conn == nil {
			log.Printf("WebSocket connection became nil for user: %s", c.UserEmail)
			break
		}

		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", c.UserEmail, err)
			} else {
				log.Printf("WebSocket read error for user %s: %v", c.UserEmail, err)
			}
			break
		}

		log.Printf("Received message from user %s: %s", c.UserEmail, string(messageBytes))

		// Parse the incoming message
		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		// Handle different message types
		switch message.Type {
		case "join_channel":
			c.handleJoinChannel(message)
		case "leave_channel":
			c.handleLeaveChannel(message)
		case "ping":
			c.handlePing()
		default:
			log.Printf("Unknown message type: %s", message.Type)
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleJoinChannel handles when a client joins a channel
func (c *Client) handleJoinChannel(message Message) {
	if message.ChannelID == "" {
		return
	}

	// Remove from previous channel if any
	if c.ChannelID != "" {
		c.hub.mutex.Lock()
		if channelClients, exists := c.hub.channelClients[c.ChannelID]; exists {
			delete(channelClients, c)
			if len(channelClients) == 0 {
				delete(c.hub.channelClients, c.ChannelID)
			}
		}
		c.hub.mutex.Unlock()
	}

	// Add to new channel
	c.ChannelID = message.ChannelID
	c.hub.mutex.Lock()
	if c.hub.channelClients[c.ChannelID] == nil {
		c.hub.channelClients[c.ChannelID] = make(map[*Client]bool)
	}
	c.hub.channelClients[c.ChannelID][c] = true
	c.hub.mutex.Unlock()

	// Send confirmation
	response := Message{
		Type:      "channel_joined",
		ChannelID: c.ChannelID,
		UserEmail: c.UserEmail,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	responseBytes, _ := json.Marshal(response)
	c.send <- responseBytes

	log.Printf("User %s joined channel %s", c.UserEmail, c.ChannelID)
}

// handleLeaveChannel handles when a client leaves a channel
func (c *Client) handleLeaveChannel(_ Message) {
	if c.ChannelID == "" {
		return
	}

	// Remove from current channel
	c.hub.mutex.Lock()
	if channelClients, exists := c.hub.channelClients[c.ChannelID]; exists {
		delete(channelClients, c)
		if len(channelClients) == 0 {
			delete(c.hub.channelClients, c.ChannelID)
		}
	}
	c.hub.mutex.Unlock()

	// Send confirmation
	response := Message{
		Type:      "channel_left",
		ChannelID: c.ChannelID,
		UserEmail: c.UserEmail,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	responseBytes, _ := json.Marshal(response)
	c.send <- responseBytes

	log.Printf("User %s left channel %s", c.UserEmail, c.ChannelID)
	c.ChannelID = ""
}

// handlePing handles ping messages
func (c *Client) handlePing() {
	response := Message{
		Type:      "pong",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	responseBytes, _ := json.Marshal(response)
	c.send <- responseBytes
}
