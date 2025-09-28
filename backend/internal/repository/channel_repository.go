package repository

import (
	"context"
	"jobsity-backend/pkg/domain"
)

// ChannelRepository defines the interface for channel data operations
type ChannelRepository interface {
	// Create creates a new channel
	Create(ctx context.Context, channel *domain.Channel) error

	// FindByID finds a channel by ID
	FindByID(ctx context.Context, id string) (*domain.Channel, error)

	// FindByName finds a channel by name
	FindByName(ctx context.Context, name string) (*domain.Channel, error)

	// FindAll returns all channels
	FindAll(ctx context.Context) ([]*domain.Channel, error)

	// Update updates an existing channel
	Update(ctx context.Context, channel *domain.Channel) error

	// Delete deletes a channel by ID
	Delete(ctx context.Context, id string) error
}
