package repository

import (
	"context"
	"jobsity-backend/pkg/domain"
)

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	// Create creates a new message
	Create(ctx context.Context, message *domain.Message) error

	// FindByID finds a message by ID
	FindByID(ctx context.Context, id string) (*domain.Message, error)

	// FindByChannelID finds all messages for a specific channel
	FindByChannelID(ctx context.Context, channelID string, limit int) ([]*domain.Message, error)

	// Update updates an existing message
	Update(ctx context.Context, message *domain.Message) error

	// Delete deletes a message by ID
	Delete(ctx context.Context, id string) error
}
