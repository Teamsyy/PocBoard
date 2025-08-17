package dto

import (
	"github.com/google/uuid"
)

// CreateElementRequest represents the request payload for creating an element
type CreateElementRequest struct {
	Kind     string      `json:"kind" validate:"required,oneof=text image sticker shape"`
	X        float64     `json:"x" validate:"required"`
	Y        float64     `json:"y" validate:"required"`
	W        float64     `json:"w" validate:"required,gt=0"`
	H        float64     `json:"h" validate:"required,gt=0"`
	Rotation float64     `json:"rotation" validate:"omitempty"`
	Payload  interface{} `json:"payload" validate:"required"`
}

// UpdateElementRequest represents the request payload for updating an element
type UpdateElementRequest struct {
	X        *float64    `json:"x,omitempty"`
	Y        *float64    `json:"y,omitempty"`
	W        *float64    `json:"w,omitempty" validate:"omitempty,gt=0"`
	H        *float64    `json:"h,omitempty" validate:"omitempty,gt=0"`
	Rotation *float64    `json:"rotation,omitempty"`
	Payload  interface{} `json:"payload,omitempty"`
}

// ReorderElementsRequest represents the request payload for batch z-index updates
type ReorderElementsRequest struct {
	Elements []ElementZIndexUpdate `json:"elements" validate:"required,min=1,dive"`
}

// ElementZIndexUpdate represents a single element z-index update
type ElementZIndexUpdate struct {
	ID uuid.UUID `json:"id" validate:"required"`
	Z  int       `json:"z" validate:"required,min=0"`
}

// ElementResponse represents an element in the response
type ElementResponse struct {
	ID       uuid.UUID   `json:"id"`
	PageID   uuid.UUID   `json:"page_id"`
	Kind     string      `json:"kind"`
	X        float64     `json:"x"`
	Y        float64     `json:"y"`
	W        float64     `json:"w"`
	H        float64     `json:"h"`
	Rotation float64     `json:"rotation"`
	Z        int         `json:"z"`
	Payload  interface{} `json:"payload"`
}

// ElementsListResponse represents the response for listing elements
type ElementsListResponse struct {
	Elements []ElementResponse `json:"elements"`
	Total    int               `json:"total"`
}
