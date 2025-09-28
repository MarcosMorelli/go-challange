package service

import (
	"context"
	"jobsity-backend/pkg/domain"
)

// ChannelService defines the interface for channel business logic
type ChannelService interface {
	// CreateChannel creates a new channel
	CreateChannel(ctx context.Context, req *domain.CreateChannelRequest, userEmail string) (*domain.Channel, error)

	// GetChannel gets a channel by ID
	GetChannel(ctx context.Context, id string) (*domain.Channel, error)

	// GetChannelByName gets a channel by name
	GetChannelByName(ctx context.Context, name string) (*domain.Channel, error)

	// GetAllChannels returns all channels
	GetAllChannels(ctx context.Context) ([]*domain.Channel, error)

	// UpdateChannel updates an existing channel
	UpdateChannel(ctx context.Context, id string, req *domain.UpdateChannelRequest, userEmail string) (*domain.Channel, error)

	// DeleteChannel deletes a channel
	DeleteChannel(ctx context.Context, id string, userEmail string) error
}
