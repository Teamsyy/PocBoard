package handlers

import (
	"junk-journal-board/internal/dto"
	"junk-journal-board/internal/services"
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ElementHandler struct {
	elementService *services.ElementService
	pageService    *services.PageService
	boardService   *services.BoardService
}

func NewElementHandler(db *gorm.DB) *ElementHandler {
	return &ElementHandler{
		elementService: services.NewElementService(db),
		pageService:    services.NewPageService(db),
		boardService:   services.NewBoardService(db),
	}
}

// CreateElement creates a new element on a page
// POST /api/v1/boards/:boardId/pages/:pageId/elements
func (h *ElementHandler) CreateElement(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*utils.Logger)

	// Parse IDs from URL
	boardIDStr := c.Params("boardId")
	pageIDStr := c.Params("pageId")

	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		logger.Warnw("Invalid board ID", "boardId", boardIDStr)
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		logger.Warnw("Invalid page ID", "pageId", pageIDStr)
		return utils.SendValidationError(c, "Invalid page ID format", nil)
	}

	// Validate edit token and board access
	editToken := c.Locals("edit_token").(uuid.UUID)
	if err := h.boardService.ValidateBoardEditAccess(boardID, editToken); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Board not found")
		}
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorizedError(c, "Invalid edit token")
		}
		logger.Errorw("Failed to validate board access", "error", err)
		return utils.SendInternalError(c, "Failed to validate board access", nil)
	}

	// Validate page belongs to board
	if err := h.pageService.ValidatePageBelongsToBoard(pageID, boardID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to validate page ownership", "error", err)
		return utils.SendInternalError(c, "Failed to validate page", nil)
	}

	// Parse request body
	var req dto.CreateElementRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warnw("Failed to parse request body", "error", err)
		return utils.SendValidationError(c, "Invalid request body", nil)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Create element
	element, err := h.elementService.CreateElement(pageID, req.Kind, req.X, req.Y, req.W, req.H, req.Rotation, req.Payload)
	if err != nil {
		logger.Errorw("Failed to create element", "error", err)
		return utils.SendInternalError(c, "Failed to create element", nil)
	}

	// Convert to response DTO
	response := dto.ElementResponse{
		ID:       element.ID,
		PageID:   element.PageID,
		Kind:     element.Kind,
		X:        element.X,
		Y:        element.Y,
		W:        element.W,
		H:        element.H,
		Rotation: element.Rotation,
		Z:        element.Z,
		Payload:  element.Payload,
	}

	logger.Infow("Element created successfully", "elementId", element.ID, "pageId", pageID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": response})
}

// GetElementsByPage lists all elements for a page
// GET /api/v1/boards/:boardId/pages/:pageId/elements
func (h *ElementHandler) GetElementsByPage(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*utils.Logger)

	// Parse IDs from URL
	boardIDStr := c.Params("boardId")
	pageIDStr := c.Params("pageId")

	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		logger.Warnw("Invalid board ID", "boardId", boardIDStr)
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		logger.Warnw("Invalid page ID", "pageId", pageIDStr)
		return utils.SendValidationError(c, "Invalid page ID format", nil)
	}

	// Validate board access (both edit and public tokens are allowed for reading)
	token := c.Locals("token")
	if token != nil {
		if err := h.boardService.ValidateBoardAccess(boardID, token.(uuid.UUID)); err != nil {
			if err == utils.ErrNotFound {
				return utils.SendNotFoundError(c, "Board not found")
			}
			if err == utils.ErrUnauthorized {
				return utils.SendUnauthorizedError(c, "Invalid token")
			}
			logger.Errorw("Failed to validate board access", "error", err)
			return utils.SendInternalError(c, "Failed to validate board access", nil)
		}
	} else {
		if err := h.boardService.ValidateBoardExists(boardID); err != nil {
			if err == utils.ErrNotFound {
				return utils.SendNotFoundError(c, "Board not found")
			}
			logger.Errorw("Failed to validate board existence", "error", err)
			return utils.SendInternalError(c, "Failed to validate board", nil)
		}
	}

	// Validate page belongs to board
	if err := h.pageService.ValidatePageBelongsToBoard(pageID, boardID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to validate page ownership", "error", err)
		return utils.SendInternalError(c, "Failed to validate page", nil)
	}

	// Get elements
	elements, err := h.elementService.GetElementsByPage(pageID)
	if err != nil {
		logger.Errorw("Failed to get elements", "error", err)
		return utils.SendInternalError(c, "Failed to get elements", nil)
	}

	// Convert to response DTOs
	elementResponses := make([]dto.ElementResponse, len(elements))
	for i, element := range elements {
		elementResponses[i] = dto.ElementResponse{
			ID:       element.ID,
			PageID:   element.PageID,
			Kind:     element.Kind,
			X:        element.X,
			Y:        element.Y,
			W:        element.W,
			H:        element.H,
			Rotation: element.Rotation,
			Z:        element.Z,
			Payload:  element.Payload,
		}
	}

	response := dto.ElementsListResponse{
		Elements: elementResponses,
		Total:    len(elementResponses),
	}

	return c.JSON(fiber.Map{"data": response})
}

// UpdateElement updates an element's properties
// PUT /api/v1/boards/:boardId/pages/:pageId/elements/:elementId
func (h *ElementHandler) UpdateElement(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*utils.Logger)

	// Parse IDs from URL
	boardIDStr := c.Params("boardId")
	pageIDStr := c.Params("pageId")
	elementIDStr := c.Params("elementId")

	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		logger.Warnw("Invalid board ID", "boardId", boardIDStr)
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		logger.Warnw("Invalid page ID", "pageId", pageIDStr)
		return utils.SendValidationError(c, "Invalid page ID format", nil)
	}

	elementID, err := uuid.Parse(elementIDStr)
	if err != nil {
		logger.Warnw("Invalid element ID", "elementId", elementIDStr)
		return utils.SendValidationError(c, "Invalid element ID format", nil)
	}

	// Validate edit token and board access
	editToken := c.Locals("edit_token").(uuid.UUID)
	if err := h.boardService.ValidateBoardEditAccess(boardID, editToken); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Board not found")
		}
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorizedError(c, "Invalid edit token")
		}
		logger.Errorw("Failed to validate board access", "error", err)
		return utils.SendInternalError(c, "Failed to validate board access", nil)
	}

	// Validate page belongs to board
	if err := h.pageService.ValidatePageBelongsToBoard(pageID, boardID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to validate page ownership", "error", err)
		return utils.SendInternalError(c, "Failed to validate page", nil)
	}

	// Validate element belongs to page
	if err := h.elementService.ValidateElementBelongsToPage(elementID, pageID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Element not found")
		}
		logger.Errorw("Failed to validate element ownership", "error", err)
		return utils.SendInternalError(c, "Failed to validate element", nil)
	}

	// Parse request body
	var req dto.UpdateElementRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warnw("Failed to parse request body", "error", err)
		return utils.SendValidationError(c, "Invalid request body", nil)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Build updates map
	updates := make(map[string]interface{})
	if req.X != nil {
		updates["x"] = *req.X
	}
	if req.Y != nil {
		updates["y"] = *req.Y
	}
	if req.W != nil {
		updates["w"] = *req.W
	}
	if req.H != nil {
		updates["h"] = *req.H
	}
	if req.Rotation != nil {
		updates["rotation"] = *req.Rotation
	}
	if req.Payload != nil {
		updates["payload"] = req.Payload
	}

	// Update element
	element, err := h.elementService.UpdateElement(elementID, updates)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Element not found")
		}
		logger.Errorw("Failed to update element", "error", err)
		return utils.SendInternalError(c, "Failed to update element", nil)
	}

	// Convert to response DTO
	response := dto.ElementResponse{
		ID:       element.ID,
		PageID:   element.PageID,
		Kind:     element.Kind,
		X:        element.X,
		Y:        element.Y,
		W:        element.W,
		H:        element.H,
		Rotation: element.Rotation,
		Z:        element.Z,
		Payload:  element.Payload,
	}

	logger.Infow("Element updated successfully", "elementId", elementID)
	return c.JSON(fiber.Map{"data": response})
}

// DeleteElement deletes an element
// DELETE /api/v1/boards/:boardId/pages/:pageId/elements/:elementId
func (h *ElementHandler) DeleteElement(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*utils.Logger)

	// Parse IDs from URL
	boardIDStr := c.Params("boardId")
	pageIDStr := c.Params("pageId")
	elementIDStr := c.Params("elementId")

	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		logger.Warnw("Invalid board ID", "boardId", boardIDStr)
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		logger.Warnw("Invalid page ID", "pageId", pageIDStr)
		return utils.SendValidationError(c, "Invalid page ID format", nil)
	}

	elementID, err := uuid.Parse(elementIDStr)
	if err != nil {
		logger.Warnw("Invalid element ID", "elementId", elementIDStr)
		return utils.SendValidationError(c, "Invalid element ID format", nil)
	}

	// Validate edit token and board access
	editToken := c.Locals("edit_token").(uuid.UUID)
	if err := h.boardService.ValidateBoardEditAccess(boardID, editToken); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Board not found")
		}
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorizedError(c, "Invalid edit token")
		}
		logger.Errorw("Failed to validate board access", "error", err)
		return utils.SendInternalError(c, "Failed to validate board access", nil)
	}

	// Validate page belongs to board
	if err := h.pageService.ValidatePageBelongsToBoard(pageID, boardID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to validate page ownership", "error", err)
		return utils.SendInternalError(c, "Failed to validate page", nil)
	}

	// Validate element belongs to page
	if err := h.elementService.ValidateElementBelongsToPage(elementID, pageID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Element not found")
		}
		logger.Errorw("Failed to validate element ownership", "error", err)
		return utils.SendInternalError(c, "Failed to validate element", nil)
	}

	// Delete element
	if err := h.elementService.DeleteElement(elementID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Element not found")
		}
		logger.Errorw("Failed to delete element", "error", err)
		return utils.SendInternalError(c, "Failed to delete element", nil)
	}

	logger.Infow("Element deleted successfully", "elementId", elementID)
	return c.SendStatus(fiber.StatusNoContent)
}

// ReorderElements handles batch z-index updates for element reordering
// PUT /api/v1/boards/:boardId/pages/:pageId/elements/reorder
func (h *ElementHandler) ReorderElements(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*utils.Logger)

	// Parse IDs from URL
	boardIDStr := c.Params("boardId")
	pageIDStr := c.Params("pageId")

	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		logger.Warnw("Invalid board ID", "boardId", boardIDStr)
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	pageID, err := uuid.Parse(pageIDStr)
	if err != nil {
		logger.Warnw("Invalid page ID", "pageId", pageIDStr)
		return utils.SendValidationError(c, "Invalid page ID format", nil)
	}

	// Validate edit token and board access
	editToken := c.Locals("edit_token").(uuid.UUID)
	if err := h.boardService.ValidateBoardEditAccess(boardID, editToken); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Board not found")
		}
		if err == utils.ErrUnauthorized {
			return utils.SendUnauthorizedError(c, "Invalid edit token")
		}
		logger.Errorw("Failed to validate board access", "error", err)
		return utils.SendInternalError(c, "Failed to validate board access", nil)
	}

	// Validate page belongs to board
	if err := h.pageService.ValidatePageBelongsToBoard(pageID, boardID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to validate page ownership", "error", err)
		return utils.SendInternalError(c, "Failed to validate page", nil)
	}

	// Parse request body
	var req dto.ReorderElementsRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warnw("Failed to parse request body", "error", err)
		return utils.SendValidationError(c, "Invalid request body", nil)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Convert to service format
	updates := make([]struct {
		ID uuid.UUID
		Z  int
	}, len(req.Elements))

	for i, element := range req.Elements {
		updates[i] = struct {
			ID uuid.UUID
			Z  int
		}{
			ID: element.ID,
			Z:  element.Z,
		}
	}

	// Perform batch update
	if err := h.elementService.BatchUpdateZIndex(pageID, updates); err != nil {
		logger.Errorw("Failed to reorder elements", "error", err)
		return utils.SendInternalError(c, "Failed to reorder elements", nil)
	}

	logger.Infow("Elements reordered successfully", "pageId", pageID, "count", len(updates))
	return c.SendStatus(fiber.StatusNoContent)
}
