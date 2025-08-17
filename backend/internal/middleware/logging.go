package middleware

import (
	"time"

	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
)

// LoggingMiddleware creates a custom logging middleware using structured logging
func LoggingMiddleware(logger *utils.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Get request ID from context (set by requestid middleware)
		requestID := c.Locals("requestid").(string)

		// Create logger with request context
		reqLogger := logger.WithRequestID(requestID)

		// Store logger in context for use in handlers
		c.Locals("logger", reqLogger)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start).Milliseconds()

		// Log the request
		reqLogger.LogRequest(
			c.Method(),
			c.Path(),
			c.Get("User-Agent"),
			c.IP(),
			c.Response().StatusCode(),
			duration,
		)

		return err
	}
}

// GetLoggerFromContext retrieves the logger from fiber context
func GetLoggerFromContext(c *fiber.Ctx) *utils.Logger {
	logger := c.Locals("logger")
	if logger == nil {
		// Fallback to a basic logger if not found
		fallbackLogger, _ := utils.NewLogger()
		return fallbackLogger
	}

	return logger.(*utils.Logger)
}
