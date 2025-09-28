package service

import (
	"context"
	"jobsity-backend/pkg/domain"
	"time"
)

// MessageService defines the interface for message business logic
type MessageService interface {
	// CreateMessage creates a new message
	CreateMessage(ctx context.Context, req *domain.CreateMessageRequest, userEmail string) (*domain.Message, error)

	// GetMessage gets a message by ID
	GetMessage(ctx context.Context, id string) (*domain.Message, error)

	// GetMessagesByChannel gets messages for a specific channel
	GetMessagesByChannel(ctx context.Context, channelID string, limit int) ([]*domain.Message, error)

	// GetMessagesByChannelAfter gets messages for a channel after a specific timestamp
	GetMessagesByChannelAfter(ctx context.Context, channelID string, after time.Time, limit int) ([]*domain.Message, error)

	// UpdateMessage updates an existing message
	UpdateMessage(ctx context.Context, id string, content string, userEmail string) (*domain.Message, error)

	// DeleteMessage deletes a message
	DeleteMessage(ctx context.Context, id string, userEmail string) error
}
