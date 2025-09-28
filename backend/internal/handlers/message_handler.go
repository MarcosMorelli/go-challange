package handlers

import (
	"jobsity-backend/internal/service"
	"jobsity-backend/pkg/domain"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// MessageHandler handles HTTP requests for message operations
type MessageHandler struct {
	messageService service.MessageService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// CreateMessage handles message creation
func (h *MessageHandler) CreateMessage(c *fiber.Ctx) error {
	var req domain.CreateMessageRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Get user email from context
	userEmail := c.Locals("userEmail").(string)

	// Call service
	message, err := h.messageService.CreateMessage(c.Context(), &req, userEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.MessageResponse{
		Success: true,
		Message: "Message created successfully",
		Data:    message,
	})
}

// GetMessage handles getting a message by ID
func (h *MessageHandler) GetMessage(c *fiber.Ctx) error {
	messageID := c.Params("id")
	if messageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: "Message ID is required",
		})
	}

	message, err := h.messageService.GetMessage(c.Context(), messageID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.MessageResponse{
			Success: false,
			Message: "Message not found",
		})
	}

	return c.JSON(domain.MessageResponse{
		Success: true,
		Message: "Message retrieved successfully",
		Data:    message,
	})
}

// GetMessagesByChannel handles getting messages for a channel
func (h *MessageHandler) GetMessagesByChannel(c *fiber.Ctx) error {
	channelID := c.Params("channelId")
	if channelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessagesResponse{
			Success: false,
			Message: "Channel ID is required",
		})
	}

	// Parse limit parameter
	limitStr := c.Query("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	messages, err := h.messageService.GetMessagesByChannel(c.Context(), channelID, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessagesResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.MessagesResponse{
		Success:  true,
		Message:  "Messages retrieved successfully",
		Messages: messages,
	})
}

// GetMessagesByChannelAfter handles getting messages for a channel after a timestamp
func (h *MessageHandler) GetMessagesByChannelAfter(c *fiber.Ctx) error {
	channelID := c.Params("channelId")
	if channelID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessagesResponse{
			Success: false,
			Message: "Channel ID is required",
		})
	}

	// Parse after parameter
	afterStr := c.Query("after")
	if afterStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessagesResponse{
			Success: false,
			Message: "After timestamp is required",
		})
	}

	after, err := time.Parse(time.RFC3339, afterStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessagesResponse{
			Success: false,
			Message: "Invalid timestamp format. Use RFC3339 format",
		})
	}

	// Parse limit parameter
	limitStr := c.Query("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	messages, err := h.messageService.GetMessagesByChannelAfter(c.Context(), channelID, after, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessagesResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.MessagesResponse{
		Success:  true,
		Message:  "Messages retrieved successfully",
		Messages: messages,
	})
}

// UpdateMessage handles message updates
func (h *MessageHandler) UpdateMessage(c *fiber.Ctx) error {
	messageID := c.Params("id")
	if messageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: "Message ID is required",
		})
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if req.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: "Message content is required",
		})
	}

	// Get user email from context
	userEmail := c.Locals("userEmail").(string)

	message, err := h.messageService.UpdateMessage(c.Context(), messageID, req.Content, userEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.MessageResponse{
		Success: true,
		Message: "Message updated successfully",
		Data:    message,
	})
}

// DeleteMessage handles message deletion
func (h *MessageHandler) DeleteMessage(c *fiber.Ctx) error {
	messageID := c.Params("id")
	if messageID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: "Message ID is required",
		})
	}

	// Get user email from context
	userEmail := c.Locals("userEmail").(string)

	err := h.messageService.DeleteMessage(c.Context(), messageID, userEmail)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(domain.MessageResponse{
		Success: true,
		Message: "Message deleted successfully",
	})
}
