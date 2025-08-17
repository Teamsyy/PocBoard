package middleware

import (
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandler creates a custom error handler that returns standardized error responses
func ErrorHandler(logger *utils.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Get request logger if available
		reqLogger := logger
		if ctxLogger := c.Locals("logger"); ctxLogger != nil {
			reqLogger = ctxLogger.(*utils.Logger)
		}

		// Default error response
		code := fiber.StatusInternalServerError
		errorCode := utils.ErrCodeInternalError
		message := "Internal server error"

		// Handle Fiber errors
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
			message = e.Message

			// Map Fiber status codes to our error codes
			switch code {
			case fiber.StatusBadRequest:
				errorCode = utils.ErrCodeBadRequest
			case fiber.StatusUnauthorized:
				errorCode = utils.ErrCodeUnauthorized
			case fiber.StatusForbidden:
				errorCode = utils.ErrCodeForbidden
			case fiber.StatusNotFound:
				errorCode = utils.ErrCodeNotFound
			case fiber.StatusUnprocessableEntity:
				errorCode = utils.ErrCodeValidationError
			default:
				errorCode = utils.ErrCodeInternalError
			}
		}

		// Log the error
		reqLogger.LogError(err, "Request error occurred",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status_code", code),
			zap.String("error_code", string(errorCode)),
		)

		// Return standardized error response
		return utils.SendError(c, code, errorCode, message, nil)
	}
}
