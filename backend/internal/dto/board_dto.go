package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateBoardRequest represents the request to create a new board
type CreateBoardRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
	Skin        string `json:"skin,omitempty" validate:"omitempty,oneof=default wood notebook cork"`
}

// UpdateBoardRequest represents the request to update a board
type UpdateBoardRequest struct {
	Title       *string `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Skin        *string `json:"skin,omitempty" validate:"omitempty,oneof=default wood notebook cork"`
}

// BoardResponse represents a board in API responses
type BoardResponse struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Skin        string         `json:"skin"`
	PublicToken uuid.UUID      `json:"public_token"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	PageCount   int            `json:"pageCount"`
	Pages       []PageResponse `json:"pages,omitempty"`
}

// CreateBoardResponse represents the response when creating a board
type CreateBoardResponse struct {
	Board     BoardWithTokensResponse `json:"board"`
	EditURL   string                  `json:"edit_url"`
	PublicURL string                  `json:"public_url"`
}

// BoardWithTokensResponse represents a board with sensitive tokens (for edit access)
type BoardWithTokensResponse struct {
	ID          uuid.UUID      `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Skin        string         `json:"skin"`
	EditToken   uuid.UUID      `json:"edit_token"`
	PublicToken uuid.UUID      `json:"public_token"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	PageCount   int            `json:"pageCount"`
	Pages       []PageResponse `json:"pages,omitempty"`
}
