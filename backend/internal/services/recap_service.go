package services

import (
	"fmt"
	"time"

	"junk-journal-board/internal/dto"
	"junk-journal-board/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecapService struct {
	db *gorm.DB
}

func NewRecapService(db *gorm.DB) *RecapService {
	return &RecapService{db: db}
}

// GetRecapData retrieves recap data for a board with date filtering
func (s *RecapService) GetRecapData(boardID uuid.UUID, filter string, date *time.Time) (*dto.RecapResponse, error) {
	// Calculate date range based on filter
	startDate, endDate := s.calculateDateRange(filter, date)

	// Get pages within the date range with element counts
	var pages []models.Page
	query := s.db.Where("board_id = ? AND date >= ? AND date <= ?", boardID, startDate, endDate).
		Order("date DESC, order_idx ASC")

	if err := query.Find(&pages).Error; err != nil {
		return nil, fmt.Errorf("failed to get pages: %w", err)
	}

	// Get element counts for each page
	pageMetadata := make([]dto.RecapPageMetadata, len(pages))
	totalElementCount := 0

	for i, page := range pages {
		var elementCount int64
		if err := s.db.Model(&models.Element{}).Where("page_id = ?", page.ID).Count(&elementCount).Error; err != nil {
			return nil, fmt.Errorf("failed to count elements for page %s: %w", page.ID, err)
		}

		pageMetadata[i] = dto.RecapPageMetadata{
			ID:           page.ID,
			Title:        page.Title,
			Date:         page.Date,
			OrderIdx:     page.OrderIdx,
			ElementCount: int(elementCount),
			CreatedAt:    page.CreatedAt,
			UpdatedAt:    page.UpdatedAt,
		}

		totalElementCount += int(elementCount)
	}

	response := &dto.RecapResponse{
		Filter: filter,
		DateRange: dto.RecapDateRange{
			StartDate: startDate,
			EndDate:   endDate,
		},
		PageCount:    len(pages),
		ElementCount: totalElementCount,
		Pages:        pageMetadata,
	}

	return response, nil
}

// calculateDateRange calculates the start and end dates based on filter and reference date
func (s *RecapService) calculateDateRange(filter string, date *time.Time) (time.Time, time.Time) {
	var referenceDate time.Time
	if date != nil {
		referenceDate = *date
	} else {
		referenceDate = time.Now()
	}

	// Normalize to start of day
	referenceDate = time.Date(referenceDate.Year(), referenceDate.Month(), referenceDate.Day(), 0, 0, 0, 0, referenceDate.Location())

	switch filter {
	case "day":
		startDate := referenceDate
		endDate := referenceDate.Add(24*time.Hour - time.Nanosecond)
		return startDate, endDate

	case "week":
		// Calculate start of week (Monday)
		weekday := int(referenceDate.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		startDate := referenceDate.AddDate(0, 0, -(weekday - 1))
		endDate := startDate.AddDate(0, 0, 7).Add(-time.Nanosecond)
		return startDate, endDate

	case "month":
		// Start of month
		startDate := time.Date(referenceDate.Year(), referenceDate.Month(), 1, 0, 0, 0, 0, referenceDate.Location())
		// End of month
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)
		return startDate, endDate

	default:
		// Default to current day if filter is invalid or empty
		startDate := referenceDate
		endDate := referenceDate.Add(24*time.Hour - time.Nanosecond)
		return startDate, endDate
	}
}
