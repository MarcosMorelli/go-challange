package service

import (
	"context"
	"jobsity-backend/internal/websocket"
	"jobsity-backend/pkg/domain"
	"time"
)

// WebSocketMessageService wraps the message service with WebSocket broadcasting
type WebSocketMessageService struct {
	messageService MessageService
	wsHandler      *websocket.Handler
}

// NewWebSocketMessageService creates a new WebSocket-enabled message service
func NewWebSocketMessageService(messageService MessageService, wsHandler *websocket.Handler) MessageService {
	return &WebSocketMessageService{
		messageService: messageService,
		wsHandler:      wsHandler,
	}
}

// CreateMessage creates a new message and broadcasts it via WebSocket
func (s *WebSocketMessageService) CreateMessage(ctx context.Context, req *domain.CreateMessageRequest, userEmail string) (*domain.Message, error) {
	// Create the message using the underlying service
	message, err := s.messageService.CreateMessage(ctx, req, userEmail)
	if err != nil {
		return nil, err
	}

	// Broadcast the new message to all clients in the channel
	s.wsHandler.BroadcastMessage(req.ChannelID, "new_message", map[string]interface{}{
		"id":         message.ID,
		"channel_id": message.ChannelID,
		"user_email": message.UserEmail,
		"content":    message.Content,
		"created_at": message.CreatedAt.Format(time.RFC3339),
	})

	return message, nil
}

// GetMessage gets a message by ID
func (s *WebSocketMessageService) GetMessage(ctx context.Context, id string) (*domain.Message, error) {
	return s.messageService.GetMessage(ctx, id)
}

// GetMessagesByChannel gets messages for a specific channel
func (s *WebSocketMessageService) GetMessagesByChannel(ctx context.Context, channelID string, limit int) ([]*domain.Message, error) {
	return s.messageService.GetMessagesByChannel(ctx, channelID, limit)
}

// UpdateMessage updates an existing message and broadcasts the update
func (s *WebSocketMessageService) UpdateMessage(ctx context.Context, id string, content string, userEmail string) (*domain.Message, error) {
	// Update the message using the underlying service
	message, err := s.messageService.UpdateMessage(ctx, id, content, userEmail)
	if err != nil {
		return nil, err
	}

	// Broadcast the message update to all clients in the channel
	s.wsHandler.BroadcastMessage(message.ChannelID, "message_updated", map[string]interface{}{
		"id":         message.ID,
		"channel_id": message.ChannelID,
		"user_email": message.UserEmail,
		"content":    message.Content,
		"created_at": message.CreatedAt.Format(time.RFC3339),
	})

	return message, nil
}

// DeleteMessage deletes a message and broadcasts the deletion
func (s *WebSocketMessageService) DeleteMessage(ctx context.Context, id string, userEmail string) error {
	// Get the message first to know which channel to broadcast to
	message, err := s.messageService.GetMessage(ctx, id)
	if err != nil {
		return err
	}

	// Delete the message using the underlying service
	err = s.messageService.DeleteMessage(ctx, id, userEmail)
	if err != nil {
		return err
	}

	// Broadcast the message deletion to all clients in the channel
	s.wsHandler.BroadcastMessage(message.ChannelID, "message_deleted", map[string]interface{}{
		"id":         message.ID,
		"channel_id": message.ChannelID,
		"user_email": message.UserEmail,
	})

	return nil
}
