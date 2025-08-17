package dto

import (
	"time"

	"github.com/google/uuid"
)

// RecapRequest represents the request payload for getting recap data
type RecapRequest struct {
	Filter string `query:"filter" validate:"omitempty,oneof=day week month"`
	Date   string `query:"date" validate:"omitempty"`
}

// RecapPageMetadata represents page metadata in recap response
type RecapPageMetadata struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Date         time.Time `json:"date"`
	OrderIdx     int       `json:"order_idx"`
	ElementCount int       `json:"element_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RecapResponse represents the response payload for recap data
type RecapResponse struct {
	Filter       string              `json:"filter"`
	DateRange    RecapDateRange      `json:"date_range"`
	PageCount    int                 `json:"page_count"`
	ElementCount int                 `json:"element_count"`
	Pages        []RecapPageMetadata `json:"pages"`
}

// RecapDateRange represents the date range for the recap
type RecapDateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
