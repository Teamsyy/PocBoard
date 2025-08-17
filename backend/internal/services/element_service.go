package services

import (
	"encoding/json"
	"fmt"

	"junk-journal-board/internal/models"
	"junk-journal-board/internal/utils"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ElementService struct {
	db *gorm.DB
}

func NewElementService(db *gorm.DB) *ElementService {
	return &ElementService{db: db}
}

// CreateElement creates a new element on a page
func (s *ElementService) CreateElement(pageID uuid.UUID, kind string, x, y, w, h, rotation float64, payload interface{}) (*models.Element, error) {
	// Get the next z-index for this page
	var maxZ int
	err := s.db.Model(&models.Element{}).
		Where("page_id = ?", pageID).
		Select("COALESCE(MAX(z), -1) + 1").
		Scan(&maxZ).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get next z-index: %w", err)
	}

	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	element := &models.Element{
		PageID:   pageID,
		Kind:     kind,
		X:        x,
		Y:        y,
		W:        w,
		H:        h,
		Rotation: rotation,
		Z:        maxZ,
		Payload:  datatypes.JSON(payloadJSON),
	}

	if err := s.db.Create(element).Error; err != nil {
		return nil, fmt.Errorf("failed to create element: %w", err)
	}

	return element, nil
}

// GetElementsByPage retrieves all elements for a page ordered by z-index
func (s *ElementService) GetElementsByPage(pageID uuid.UUID) ([]models.Element, error) {
	var elements []models.Element
	err := s.db.Where("page_id = ?", pageID).
		Order("z ASC").
		Find(&elements).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get elements: %w", err)
	}

	return elements, nil
}

// GetElementByID retrieves a single element by ID
func (s *ElementService) GetElementByID(elementID uuid.UUID) (*models.Element, error) {
	var element models.Element
	err := s.db.First(&element, "id = ?", elementID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get element: %w", err)
	}

	return &element, nil
}

// UpdateElement updates element properties
func (s *ElementService) UpdateElement(elementID uuid.UUID, updates map[string]interface{}) (*models.Element, error) {
	var element models.Element
	if err := s.db.First(&element, "id = ?", elementID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find element: %w", err)
	}

	// Handle payload separately if it exists
	if payload, exists := updates["payload"]; exists {
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		updates["payload"] = datatypes.JSON(payloadJSON)
	}

	// Update the element
	if err := s.db.Model(&element).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update element: %w", err)
	}

	// Reload the element to get updated values
	if err := s.db.First(&element, "id = ?", elementID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload element: %w", err)
	}

	return &element, nil
}

// DeleteElement deletes an element
func (s *ElementService) DeleteElement(elementID uuid.UUID) error {
	result := s.db.Delete(&models.Element{}, "id = ?", elementID)
	if result.Error != nil {
		return fmt.Errorf("failed to delete element: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return utils.ErrNotFound
	}

	return nil
}

// BatchUpdateZIndex updates z-index for multiple elements in a single transaction
func (s *ElementService) BatchUpdateZIndex(pageID uuid.UUID, updates []struct {
	ID uuid.UUID
	Z  int
}) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			// Verify element belongs to the page
			var count int64
			err := tx.Model(&models.Element{}).
				Where("id = ? AND page_id = ?", update.ID, pageID).
				Count(&count).Error
			if err != nil {
				return fmt.Errorf("failed to verify element ownership: %w", err)
			}
			if count == 0 {
				return fmt.Errorf("element %s not found on page %s", update.ID, pageID)
			}

			// Update z-index
			err = tx.Model(&models.Element{}).
				Where("id = ?", update.ID).
				Update("z", update.Z).Error
			if err != nil {
				return fmt.Errorf("failed to update z-index for element %s: %w", update.ID, err)
			}
		}
		return nil
	})
}

// ValidateElementBelongsToPage checks if an element belongs to the specified page
func (s *ElementService) ValidateElementBelongsToPage(elementID, pageID uuid.UUID) error {
	var count int64
	err := s.db.Model(&models.Element{}).
		Where("id = ? AND page_id = ?", elementID, pageID).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to validate element ownership: %w", err)
	}

	if count == 0 {
		return utils.ErrNotFound
	}

	return nil
}

// GetNextZIndex returns the next available z-index for a page
func (s *ElementService) GetNextZIndex(pageID uuid.UUID) (int, error) {
	var maxZ int
	err := s.db.Model(&models.Element{}).
		Where("page_id = ?", pageID).
		Select("COALESCE(MAX(z), -1) + 1").
		Scan(&maxZ).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get next z-index: %w", err)
	}

	return maxZ, nil
}
