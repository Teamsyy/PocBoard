package services

import (
	"fmt"
	"time"

	"junk-journal-board/internal/models"
	"junk-journal-board/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageService struct {
	db *gorm.DB
}

func NewPageService(db *gorm.DB) *PageService {
	return &PageService{db: db}
}

// CreatePage creates a new page with proper ordering
func (s *PageService) CreatePage(boardID uuid.UUID, title string, date time.Time) (*models.Page, error) {
	// Get the next order index for this board
	var maxOrder int
	err := s.db.Model(&models.Page{}).
		Where("board_id = ?", boardID).
		Select("COALESCE(MAX(order_idx), -1) + 1").
		Scan(&maxOrder).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get next order index: %w", err)
	}

	page := &models.Page{
		BoardID:  boardID,
		Title:    title,
		Date:     date,
		OrderIdx: maxOrder,
	}

	if err := s.db.Create(page).Error; err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	return page, nil
}

// GetPagesByBoard retrieves all pages for a board ordered by order_idx
func (s *PageService) GetPagesByBoard(boardID uuid.UUID) ([]models.Page, error) {
	var pages []models.Page
	err := s.db.Preload("Elements", func(db *gorm.DB) *gorm.DB {
		return db.Order("z ASC")
	}).Where("board_id = ?", boardID).
		Order("order_idx ASC").
		Find(&pages).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get pages: %w", err)
	}

	return pages, nil
}

// GetPageByID retrieves a single page with its elements
func (s *PageService) GetPageByID(pageID uuid.UUID) (*models.Page, error) {
	var page models.Page
	err := s.db.Preload("Elements", func(db *gorm.DB) *gorm.DB {
		return db.Order("z ASC")
	}).First(&page, "id = ?", pageID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get page: %w", err)
	}

	return &page, nil
}

// UpdatePage updates page metadata
func (s *PageService) UpdatePage(pageID uuid.UUID, title string, date time.Time, orderIdx *int) (*models.Page, error) {
	var page models.Page
	if err := s.db.First(&page, "id = ?", pageID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find page: %w", err)
	}

	// Update fields
	page.Title = title
	page.Date = date

	// Handle order index update if provided
	if orderIdx != nil {
		if err := s.updatePageOrder(&page, *orderIdx); err != nil {
			return nil, fmt.Errorf("failed to update page order: %w", err)
		}
	}

	if err := s.db.Save(&page).Error; err != nil {
		return nil, fmt.Errorf("failed to update page: %w", err)
	}

	return &page, nil
}

// DeletePage deletes a page and all its elements (cascade delete)
func (s *PageService) DeletePage(pageID uuid.UUID) error {
	// Start a transaction to ensure atomicity
	return s.db.Transaction(func(tx *gorm.DB) error {
		// First, delete all elements associated with this page
		if err := tx.Where("page_id = ?", pageID).Delete(&models.Element{}).Error; err != nil {
			return fmt.Errorf("failed to delete page elements: %w", err)
		}

		// Then delete the page itself
		result := tx.Delete(&models.Page{}, "id = ?", pageID)
		if result.Error != nil {
			return fmt.Errorf("failed to delete page: %w", result.Error)
		}

		if result.RowsAffected == 0 {
			return utils.ErrNotFound
		}

		return nil
	})
}

// updatePageOrder handles reordering pages when order_idx changes
func (s *PageService) updatePageOrder(page *models.Page, newOrderIdx int) error {
	if page.OrderIdx == newOrderIdx {
		return nil // No change needed
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		oldOrderIdx := page.OrderIdx

		if newOrderIdx > oldOrderIdx {
			// Moving down: shift pages between old and new position up
			err := tx.Model(&models.Page{}).
				Where("board_id = ? AND order_idx > ? AND order_idx <= ?",
					page.BoardID, oldOrderIdx, newOrderIdx).
				Update("order_idx", gorm.Expr("order_idx - 1")).Error
			if err != nil {
				return err
			}
		} else {
			// Moving up: shift pages between new and old position down
			err := tx.Model(&models.Page{}).
				Where("board_id = ? AND order_idx >= ? AND order_idx < ?",
					page.BoardID, newOrderIdx, oldOrderIdx).
				Update("order_idx", gorm.Expr("order_idx + 1")).Error
			if err != nil {
				return err
			}
		}

		// Update the page's order index
		page.OrderIdx = newOrderIdx
		return nil
	})
}

// ValidatePageBelongsToBoard checks if a page belongs to the specified board
func (s *PageService) ValidatePageBelongsToBoard(pageID, boardID uuid.UUID) error {
	var count int64
	err := s.db.Model(&models.Page{}).
		Where("id = ? AND board_id = ?", pageID, boardID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to validate page ownership: %w", err)
	}

	if count == 0 {
		return utils.ErrNotFound
	}

	return nil
}
