package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Common errors
var (
	ErrNotFound     = errors.New("resource not found")
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("forbidden access")
)

// Global validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeInternalError   ErrorCode = "INTERNAL_ERROR"
	ErrCodeValidationError ErrorCode = "VALIDATION_ERROR"
	ErrCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
)

// ErrorResponse represents the standardized error response format
type ErrorResponse struct {
	Error struct {
		Code    ErrorCode   `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(code ErrorCode, message string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Error: struct {
			Code    ErrorCode   `json:"code"`
			Message string      `json:"message"`
			Details interface{} `json:"details,omitempty"`
		}{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// SendError sends a standardized error response
func SendError(c *fiber.Ctx, statusCode int, code ErrorCode, message string, details interface{}) error {
	response := NewErrorResponse(code, message, details)
	return c.Status(statusCode).JSON(response)
}

// SendBadRequest sends a 400 Bad Request error
func SendBadRequest(c *fiber.Ctx, message string, details interface{}) error {
	return SendError(c, fiber.StatusBadRequest, ErrCodeBadRequest, message, details)
}

// SendUnauthorized sends a 401 Unauthorized error
func SendUnauthorized(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusUnauthorized, ErrCodeUnauthorized, message, nil)
}

// SendNotFound sends a 404 Not Found error
func SendNotFound(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusNotFound, ErrCodeNotFound, message, nil)
}

// SendForbidden sends a 403 Forbidden error
func SendForbidden(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusForbidden, ErrCodeForbidden, message, nil)
}

// SendInternalError sends a 500 Internal Server Error
func SendInternalError(c *fiber.Ctx, message string, details interface{}) error {
	return SendError(c, fiber.StatusInternalServerError, ErrCodeInternalError, message, details)
}

// SendValidationError sends a 422 Unprocessable Entity error for validation failures
func SendValidationError(c *fiber.Ctx, message string, details interface{}) error {
	return SendError(c, fiber.StatusUnprocessableEntity, ErrCodeValidationError, message, details)
}

// SendDatabaseError sends a 500 Internal Server Error for database issues
func SendDatabaseError(c *fiber.Ctx, message string) error {
	return SendError(c, fiber.StatusInternalServerError, ErrCodeDatabaseError, message, nil)
}

// ValidateStruct validates a struct using the validator package
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, formatValidationError(err))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}
	return nil
}

// formatValidationError formats a validation error into a human-readable message
func formatValidationError(err validator.FieldError) string {
	field := strings.ToLower(err.Field())

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// Convenience functions for common error responses
func SendNotFoundError(c *fiber.Ctx, message string) error {
	return SendNotFound(c, message)
}

func SendUnauthorizedError(c *fiber.Ctx, message string) error {
	return SendUnauthorized(c, message)
}

func SendValidationErrorSimple(c *fiber.Ctx, message string) error {
	return SendValidationError(c, message, nil)
}

func SendInternalErrorSimple(c *fiber.Ctx, message string) error {
	return SendInternalError(c, message, nil)
}

// SendBadRequestError sends a 400 Bad Request error with simple message
func SendBadRequestError(c *fiber.Ctx, message string) error {
	return SendBadRequest(c, message, nil)
}

// NewValidationError creates a validation error
func NewValidationError(message string) error {
	return fmt.Errorf("validation error: %s", message)
}
