package services

import (
	"fmt"

	"junk-journal-board/internal/models"
	"junk-journal-board/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardService struct {
	db *gorm.DB
}

func NewBoardService(db *gorm.DB) *BoardService {
	return &BoardService{db: db}
}

// CreateBoard creates a new board with generated tokens
func (s *BoardService) CreateBoard(title, description, skin string) (*models.Board, error) {
	board := &models.Board{
		Title:       title,
		Description: description,
		Skin:        skin,
	}

	if err := s.db.Create(board).Error; err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	return board, nil
}

// GetBoardByEditToken retrieves a board by its edit token with pages
func (s *BoardService) GetBoardByEditToken(editToken uuid.UUID) (*models.Board, error) {
	var board models.Board
	err := s.db.Preload("Pages", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_idx ASC")
	}).Where("edit_token = ?", editToken).First(&board).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get board by edit token: %w", err)
	}

	return &board, nil
}

// GetBoardByPublicToken retrieves a board by its public token with pages (read-only)
func (s *BoardService) GetBoardByPublicToken(publicToken uuid.UUID) (*models.Board, error) {
	var board models.Board
	err := s.db.Preload("Pages", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_idx ASC")
	}).Where("public_token = ?", publicToken).First(&board).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get board by public token: %w", err)
	}

	return &board, nil
}

// GetBoardByID retrieves a board by its ID
func (s *BoardService) GetBoardByID(boardID uuid.UUID) (*models.Board, error) {
	var board models.Board
	err := s.db.Preload("Pages", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_idx ASC")
	}).Where("id = ?", boardID).First(&board).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get board by ID: %w", err)
	}

	return &board, nil
}

// UpdateBoard updates a board's properties
func (s *BoardService) UpdateBoard(boardID uuid.UUID, title, description, skin *string) (*models.Board, error) {
	var board models.Board
	if err := s.db.Where("id = ?", boardID).First(&board).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find board: %w", err)
	}

	updates := make(map[string]interface{})
	if title != nil {
		updates["title"] = *title
	}
	if description != nil {
		updates["description"] = *description
	}
	if skin != nil {
		updates["skin"] = *skin
	}

	if len(updates) > 0 {
		if err := s.db.Model(&board).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("failed to update board: %w", err)
		}
	}

	// Reload the board to get updated values
	if err := s.db.Preload("Pages", func(db *gorm.DB) *gorm.DB {
		return db.Order("order_idx ASC")
	}).Where("id = ?", boardID).First(&board).Error; err != nil {
		return nil, fmt.Errorf("failed to reload board: %w", err)
	}

	return &board, nil
}

// ValidateBoardEditAccess validates that the edit token is valid for the board
func (s *BoardService) ValidateBoardEditAccess(boardID, editToken uuid.UUID) error {
	var count int64
	err := s.db.Model(&models.Board{}).
		Where("id = ? AND edit_token = ?", boardID, editToken).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to validate board edit access: %w", err)
	}

	if count == 0 {
		// Check if board exists to differentiate between not found and unauthorized
		var boardCount int64
		if err := s.db.Model(&models.Board{}).Where("id = ?", boardID).Count(&boardCount).Error; err != nil {
			return fmt.Errorf("failed to check board existence: %w", err)
		}
		if boardCount == 0 {
			return utils.ErrNotFound
		}
		return utils.ErrUnauthorized
	}

	return nil
}

// ValidateBoardAccess validates that the token (edit or public) is valid for the board
func (s *BoardService) ValidateBoardAccess(boardID, token uuid.UUID) error {
	var count int64
	err := s.db.Model(&models.Board{}).
		Where("id = ? AND (edit_token = ? OR public_token = ?)", boardID, token, token).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to validate board access: %w", err)
	}

	if count == 0 {
		// Check if board exists to differentiate between not found and unauthorized
		var boardCount int64
		if err := s.db.Model(&models.Board{}).Where("id = ?", boardID).Count(&boardCount).Error; err != nil {
			return fmt.Errorf("failed to check board existence: %w", err)
		}
		if boardCount == 0 {
			return utils.ErrNotFound
		}
		return utils.ErrUnauthorized
	}

	return nil
}

// ValidateBoardExists checks if a board exists (for public access without token)
func (s *BoardService) ValidateBoardExists(boardID uuid.UUID) error {
	var count int64
	err := s.db.Model(&models.Board{}).
		Where("id = ?", boardID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to validate board existence: %w", err)
	}

	if count == 0 {
		return utils.ErrNotFound
	}

	return nil
}

// GetAllBoards retrieves all boards with page count
func (s *BoardService) GetAllBoards() ([]models.Board, error) {
	var boards []models.Board
	err := s.db.Preload("Pages").Find(&boards).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all boards: %w", err)
	}

	return boards, nil
}

// DeleteBoard deletes a board and all its associated data
func (s *BoardService) DeleteBoard(boardID uuid.UUID) error {
	// Start a transaction to ensure all related data is deleted
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First, delete all elements in all pages of this board
	if err := tx.Exec(`
		DELETE FROM elements 
		WHERE page_id IN (
			SELECT id FROM pages WHERE board_id = ?
		)
	`, boardID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete board elements: %w", err)
	}

	// Then delete all pages of this board
	if err := tx.Where("board_id = ?", boardID).Delete(&models.Page{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete board pages: %w", err)
	}

	// Finally delete the board itself
	if err := tx.Where("id = ?", boardID).Delete(&models.Board{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete board: %w", err)
	}

	return tx.Commit().Error
}
