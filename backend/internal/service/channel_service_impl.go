package service

import (
	"context"
	"errors"
	"jobsity-backend/internal/repository"
	"jobsity-backend/pkg/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

// ChannelServiceImpl implements ChannelService
type ChannelServiceImpl struct {
	channelRepo repository.ChannelRepository
}

// NewChannelService creates a new channel service
func NewChannelService(channelRepo repository.ChannelRepository) ChannelService {
	return &ChannelServiceImpl{
		channelRepo: channelRepo,
	}
}

// CreateChannel creates a new channel
func (s *ChannelServiceImpl) CreateChannel(ctx context.Context, req *domain.CreateChannelRequest, userEmail string) (*domain.Channel, error) {
	// Validate input
	if req.Name == "" {
		return nil, errors.New("channel name is required")
	}

	// Check if channel with same name already exists
	existingChannel, err := s.channelRepo.FindByName(ctx, req.Name)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}
	if existingChannel != nil {
		return nil, errors.New("channel with this name already exists")
	}

	// Create new channel
	newChannel := &domain.Channel{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userEmail,
	}

	err = s.channelRepo.Create(ctx, newChannel)
	if err != nil {
		return nil, err
	}

	return newChannel, nil
}

// GetChannel gets a channel by ID
func (s *ChannelServiceImpl) GetChannel(ctx context.Context, id string) (*domain.Channel, error) {
	return s.channelRepo.FindByID(ctx, id)
}

// GetChannelByName gets a channel by name
func (s *ChannelServiceImpl) GetChannelByName(ctx context.Context, name string) (*domain.Channel, error) {
	return s.channelRepo.FindByName(ctx, name)
}

// GetAllChannels returns all channels
func (s *ChannelServiceImpl) GetAllChannels(ctx context.Context) ([]*domain.Channel, error) {
	return s.channelRepo.FindAll(ctx)
}

// UpdateChannel updates an existing channel
func (s *ChannelServiceImpl) UpdateChannel(ctx context.Context, id string, req *domain.UpdateChannelRequest, userEmail string) (*domain.Channel, error) {
	// Get existing channel
	channel, err := s.channelRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user is the creator
	if channel.CreatedBy != userEmail {
		return nil, errors.New("only the channel creator can update it")
	}

	// Update channel fields
	if req.Name != "" {
		// Check if new name conflicts with existing channel
		if req.Name != channel.Name {
			existingChannel, err := s.channelRepo.FindByName(ctx, req.Name)
			if err != nil && err != mongo.ErrNoDocuments {
				return nil, err
			}
			if existingChannel != nil {
				return nil, errors.New("channel with this name already exists")
			}
		}
		channel.Name = req.Name
	}
	if req.Description != "" {
		channel.Description = req.Description
	}

	// Save updated channel
	err = s.channelRepo.Update(ctx, channel)
	if err != nil {
		return nil, err
	}

	return channel, nil
}

// DeleteChannel deletes a channel
func (s *ChannelServiceImpl) DeleteChannel(ctx context.Context, id string, userEmail string) error {
	// Get existing channel
	channel, err := s.channelRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user is the creator
	if channel.CreatedBy != userEmail {
		return errors.New("only the channel creator can delete it")
	}

	return s.channelRepo.Delete(ctx, id)
}
