package domain

import "time"

// Channel represents a chat channel
type Channel struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	CreatedBy   string    `bson:"created_by" json:"created_by"` // User email
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at" json:"updated_at"`
}

// CreateChannelRequest represents the create channel request structure
type CreateChannelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateChannelRequest represents the update channel request structure
type UpdateChannelRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ChannelResponse represents the channel response structure
type ChannelResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Channel *Channel `json:"channel,omitempty"`
}

// ChannelsResponse represents the channels list response structure
type ChannelsResponse struct {
	Success  bool       `json:"success"`
	Message  string     `json:"message"`
	Channels []*Channel `json:"channels,omitempty"`
}
