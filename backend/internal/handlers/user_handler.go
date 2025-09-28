package handlers

import (
	"jobsity-backend/internal/service"
	"jobsity-backend/pkg/domain"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Login handles user login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.LoginResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	// Call service
	response, err := h.userService.Login(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.LoginResponse{
			Success: false,
			Message: "Internal server error",
		})
	}

	// Return appropriate status code based on success
	if response.Success {
		return c.JSON(response)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(response)
	}
}

// CreateUser handles user creation
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req domain.CreateUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid request body",
		})
	}

	// Call service
	user, err := h.userService.CreateUser(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "User created successfully",
		"user":    user,
	})
}

// GetUser handles getting user by email
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Email is required",
		})
	}

	user, err := h.userService.GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}
