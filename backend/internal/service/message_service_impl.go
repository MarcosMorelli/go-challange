package service

import (
	"context"
	"errors"
	"jobsity-backend/internal/repository"
	"jobsity-backend/pkg/domain"

	"go.mongodb.org/mongo-driver/mongo"
)

// MessageServiceImpl implements MessageService
type MessageServiceImpl struct {
	messageRepo repository.MessageRepository
	channelRepo repository.ChannelRepository
}

// NewMessageService creates a new message service
func NewMessageService(messageRepo repository.MessageRepository, channelRepo repository.ChannelRepository) MessageService {
	return &MessageServiceImpl{
		messageRepo: messageRepo,
		channelRepo: channelRepo,
	}
}

// CreateMessage creates a new message
func (s *MessageServiceImpl) CreateMessage(ctx context.Context, req *domain.CreateMessageRequest, userEmail string) (*domain.Message, error) {
	// Validate input
	if req.ChannelID == "" {
		return nil, errors.New("channel ID is required")
	}
	if req.Content == "" {
		return nil, errors.New("message content is required")
	}

	// Verify channel exists
	_, err := s.channelRepo.FindByID(ctx, req.ChannelID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("channel not found")
		}
		return nil, err
	}

	// Create new message
	newMessage := &domain.Message{
		ChannelID: req.ChannelID,
		UserEmail: userEmail,
		Content:   req.Content,
	}

	err = s.messageRepo.Create(ctx, newMessage)
	if err != nil {
		return nil, err
	}

	return newMessage, nil
}

// GetMessage gets a message by ID
func (s *MessageServiceImpl) GetMessage(ctx context.Context, id string) (*domain.Message, error) {
	return s.messageRepo.FindByID(ctx, id)
}

// GetMessagesByChannel gets messages for a specific channel
func (s *MessageServiceImpl) GetMessagesByChannel(ctx context.Context, channelID string, limit int) ([]*domain.Message, error) {
	// Verify channel exists
	_, err := s.channelRepo.FindByID(ctx, channelID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("channel not found")
		}
		return nil, err
	}

	return s.messageRepo.FindByChannelID(ctx, channelID, limit)
}

// UpdateMessage updates an existing message
func (s *MessageServiceImpl) UpdateMessage(ctx context.Context, id string, content string, userEmail string) (*domain.Message, error) {
	// Get existing message
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if user is the message author
	if message.UserEmail != userEmail {
		return nil, errors.New("only the message author can update it")
	}

	// Update message content
	message.Content = content

	// Save updated message
	err = s.messageRepo.Update(ctx, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// DeleteMessage deletes a message
func (s *MessageServiceImpl) DeleteMessage(ctx context.Context, id string, userEmail string) error {
	// Get existing message
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if user is the message author
	if message.UserEmail != userEmail {
		return errors.New("only the message author can delete it")
	}

	return s.messageRepo.Delete(ctx, id)
}
