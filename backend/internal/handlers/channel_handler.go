package handlers

import (
	"jobsity-backend/internal/service"
	"jobsity-backend/pkg/domain"

	"github.com/gofiber/fiber/v2"
)

// ChannelHandler handles HTTP requests for channel operations
type ChannelHandler struct {
	channelService service.ChannelService
}

// NewChannelHandler creates a new channel handler
func NewChannelHandler(channelService service.ChannelService) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
	}
}

// CreateChannel handles channel creation
func (h *ChannelHandler) CreateChannel(c *fiber.Ctx) error {
	var req domain.CreateChannelRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Get user email from context
	userEmail := c.Locals("userEmail").(string)

	// Call service
	channel, err := h.channelService.CreateChannel(c.Context(), &req, userEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.ChannelResponse{
		Success: true,
		Message: "Channel created successfully",
		Channel: channel,
	})
}

// GetChannel handles getting a channel by ID
func (h *ChannelHandler) GetChannel(c *fiber.Ctx) error {
	channelID := c.Params("id")
	if channelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Channel ID is required",
		})
	}

	channel, err := h.channelService.GetChannel(c.Context(), channelID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Channel not found",
		})
	}

	return c.JSON(domain.ChannelResponse{
		Success: true,
		Message: "Channel retrieved successfully",
		Channel: channel,
	})
}

// GetChannelByName handles getting a channel by name
func (h *ChannelHandler) GetChannelByName(c *fiber.Ctx) error {
	channelName := c.Params("name")
	if channelName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Channel name is required",
		})
	}

	channel, err := h.channelService.GetChannelByName(c.Context(), channelName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Channel not found",
		})
	}

	return c.JSON(domain.ChannelResponse{
		Success: true,
		Message: "Channel retrieved successfully",
		Channel: channel,
	})
}

// GetAllChannels handles getting all channels
func (h *ChannelHandler) GetAllChannels(c *fiber.Ctx) error {
	channels, err := h.channelService.GetAllChannels(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.ChannelsResponse{
			Success: false,
			Message: "Failed to retrieve channels",
		})
	}

	return c.JSON(domain.ChannelsResponse{
		Success:  true,
		Message:  "Channels retrieved successfully",
		Channels: channels,
	})
}

// UpdateChannel handles channel updates
func (h *ChannelHandler) UpdateChannel(c *fiber.Ctx) error {
	channelID := c.Params("id")
	if channelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Channel ID is required",
		})
	}

	var req domain.UpdateChannelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Get user email from context
	userEmail := c.Locals("userEmail").(string)

	channel, err := h.channelService.UpdateChannel(c.Context(), channelID, &req, userEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.ChannelResponse{
		Success: true,
		Message: "Channel updated successfully",
		Channel: channel,
	})
}

// DeleteChannel handles channel deletion
func (h *ChannelHandler) DeleteChannel(c *fiber.Ctx) error {
	channelID := c.Params("id")
	if channelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: "Channel ID is required",
		})
	}

	// Get user email from context
	userEmail := c.Locals("userEmail").(string)

	err := h.channelService.DeleteChannel(c.Context(), channelID, userEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.ChannelResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.ChannelResponse{
		Success: true,
		Message: "Channel deleted successfully",
	})
}
