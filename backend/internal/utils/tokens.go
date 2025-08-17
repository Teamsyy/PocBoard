package utils

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// ParseUUID parses a string into a UUID, returns zero UUID if invalid
func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

// ValidateToken validates that a token string is a valid UUID
func ValidateToken(token string) (uuid.UUID, bool) {
	if token == "" {
		return uuid.Nil, false
	}

	parsedUUID, err := uuid.Parse(token)
	if err != nil {
		return uuid.Nil, false
	}

	return parsedUUID, true
}
