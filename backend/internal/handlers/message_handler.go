package handlers

import (
	"jobsity-backend/internal/service"
	"jobsity-backend/pkg/domain"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/streadway/amqp"
)

// MessageHandler handles HTTP requests for message operations
type MessageHandler struct {
	messageService service.MessageService
	rabbitMQ       *amqp.Channel
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageService service.MessageService, rabbitMQ *amqp.Channel) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
		rabbitMQ:       rabbitMQ,
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

	// Check if it's a stock command
	if strings.HasPrefix(req.Content, "/stock=") {
		stockCode := strings.TrimPrefix(req.Content, "/stock=")
		if stockCode == "" {
			return c.Status(fiber.StatusBadRequest).JSON(domain.MessageResponse{
				Success: false,
				Message: "Stock code is required",
			})
		}

		// Publish to RabbitMQ for stock bot processing
		command := req.ChannelID + "|" + userEmail + "|" + stockCode
		err := h.rabbitMQ.Publish(
			"",               // exchange
			"stock_commands", // routing key
			false,            // mandatory
			false,            // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(command),
			},
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(domain.MessageResponse{
				Success: false,
				Message: "Failed to process stock command",
			})
		}

		// Return success but don't save the command to database
		return c.Status(fiber.StatusOK).JSON(domain.MessageResponse{
			Success: true,
			Message: "Stock command processed",
		})
	}

	// Call service for regular messages
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
