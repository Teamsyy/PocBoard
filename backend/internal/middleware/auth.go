package middleware

import (
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TokenValidationMiddleware validates edit tokens for mutation operations
func TokenValidationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get edit_token from query parameters
		editTokenStr := c.Query("edit_token")
		if editTokenStr == "" {
			return utils.SendUnauthorized(c, "Edit token is required for this operation")
		}

		// Validate token format
		editToken, valid := utils.ValidateToken(editTokenStr)
		if !valid {
			return utils.SendUnauthorized(c, "Invalid edit token format")
		}

		// Store the parsed token in context for use in handlers
		c.Locals("edit_token", editToken)

		return c.Next()
	}
}

// OptionalTokenMiddleware extracts tokens but doesn't require them
func OptionalTokenMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get edit_token from query parameters
		editTokenStr := c.Query("edit_token")
		if editTokenStr != "" {
			if editToken, valid := utils.ValidateToken(editTokenStr); valid {
				c.Locals("edit_token", editToken)
				c.Locals("token", editToken) // Store generic token for access validation
			}
		}

		// Get public_token from query parameters
		publicTokenStr := c.Query("public_token")
		if publicTokenStr != "" {
			if publicToken, valid := utils.ValidateToken(publicTokenStr); valid {
				c.Locals("public_token", publicToken)
				// Only set generic token if edit token wasn't already set
				if c.Locals("token") == nil {
					c.Locals("token", publicToken)
				}
			}
		}

		return c.Next()
	}
}

// GetEditTokenFromContext retrieves the edit token from fiber context
func GetEditTokenFromContext(c *fiber.Ctx) (uuid.UUID, bool) {
	token := c.Locals("edit_token")
	if token == nil {
		return uuid.Nil, false
	}

	editToken, ok := token.(uuid.UUID)
	return editToken, ok
}

// GetPublicTokenFromContext retrieves the public token from fiber context
func GetPublicTokenFromContext(c *fiber.Ctx) (uuid.UUID, bool) {
	token := c.Locals("public_token")
	if token == nil {
		return uuid.Nil, false
	}

	publicToken, ok := token.(uuid.UUID)
	return publicToken, ok
}
