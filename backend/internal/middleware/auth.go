package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware is a simple middleware that extracts user email from headers
// In a real application, this would validate JWT tokens or session cookies
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// For testing purposes, we'll use the User-Email header
		// In production, you would validate JWT tokens here
		userEmail := c.Get("User-Email")
		if userEmail == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "User not authenticated",
			})
		}

		// Set user email in context for handlers to use
		c.Locals("userEmail", userEmail)
		return c.Next()
	}
}

// OptionalAuthMiddleware is a middleware that extracts user email if present
// but doesn't require authentication
func OptionalAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userEmail := c.Get("User-Email")
		if userEmail != "" {
			c.Locals("userEmail", userEmail)
		}
		return c.Next()
	}
}
