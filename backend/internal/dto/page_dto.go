package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreatePageRequest represents the request payload for creating a page
type CreatePageRequest struct {
	Title string    `json:"title" validate:"required,min=1,max=255"`
	Date  time.Time `json:"date" validate:"required"`
}

// UpdatePageRequest represents the request payload for updating a page
type UpdatePageRequest struct {
	Title    string    `json:"title" validate:"required,min=1,max=255"`
	Date     time.Time `json:"date" validate:"required"`
	OrderIdx *int      `json:"order_idx,omitempty" validate:"omitempty,min=0"`
}

// PageResponse represents the response payload for a page
type PageResponse struct {
	ID        uuid.UUID `json:"id"`
	BoardID   uuid.UUID `json:"board_id"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	OrderIdx  int       `json:"order_idx"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PageWithElementsResponse represents a page response with its elements
type PageWithElementsResponse struct {
	ID        uuid.UUID         `json:"id"`
	BoardID   uuid.UUID         `json:"board_id"`
	Title     string            `json:"title"`
	Date      time.Time         `json:"date"`
	OrderIdx  int               `json:"order_idx"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Elements  []ElementResponse `json:"elements"`
}

// PagesListResponse represents the response for listing pages
type PagesListResponse struct {
	Pages []PageResponse `json:"pages"`
	Total int            `json:"total"`
}

// PagesWithElementsListResponse represents the response for listing pages with elements
type PagesWithElementsListResponse struct {
	Pages []PageWithElementsResponse `json:"pages"`
	Total int                        `json:"total"`
}
