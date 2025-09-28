package domain

import "time"

// Message represents a chat message
type Message struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	ChannelID string    `bson:"channel_id" json:"channel_id"`
	UserEmail string    `bson:"user_email" json:"user_email"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// CreateMessageRequest represents the create message request structure
type CreateMessageRequest struct {
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

// MessageResponse represents the message response structure
type MessageResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    *Message `json:"data,omitempty"`
}

// MessagesResponse represents the messages list response structure
type MessagesResponse struct {
	Success  bool       `json:"success"`
	Message  string     `json:"message"`
	Messages []*Message `json:"messages,omitempty"`
}
