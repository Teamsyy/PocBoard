package handlers

import (
	"junk-journal-board/internal/dto"
	"junk-journal-board/internal/services"
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageHandler struct {
	pageService  *services.PageService
	boardService *services.BoardService
}

func NewPageHandler(db *gorm.DB) *PageHandler {
	return &PageHandler{
		pageService:  services.NewPageService(db),
		boardService: services.NewBoardService(db),
	}
}

// CreatePage creates a new page for a board
// POST /api/v1/boards/:boardId/pages
func (h *PageHandler) CreatePage(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*utils.Logger)

	// Parse board ID from URL
	boardIDStr := c.Params("boardId")
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		logger.Warnw("Invalid board ID", "boardId", boardIDStr)
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	// Validate edit token and board existence
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

	// Parse request body
	var req dto.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warnw("Failed to parse request body", "error", err)
		return utils.SendValidationError(c, "Invalid request body", nil)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Create page
	page, err := h.pageService.CreatePage(boardID, req.Title, req.Date)
	if err != nil {
		logger.Errorw("Failed to create page", "error", err)
		return utils.SendInternalError(c, "Failed to create page", nil)
	}

	// Convert to response DTO
	response := dto.PageResponse{
		ID:        page.ID,
		BoardID:   page.BoardID,
		Title:     page.Title,
		Date:      page.Date,
		OrderIdx:  page.OrderIdx,
		CreatedAt: page.CreatedAt,
		UpdatedAt: page.UpdatedAt,
	}

	logger.Infow("Page created successfully", "pageId", page.ID, "boardId", boardID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": response})
}

// GetPagesByBoard lists all pages for a board
// GET /api/v1/boards/:boardId/pages
func (h *PageHandler) GetPagesByBoard(c *fiber.Ctx) error {
	logger := c.Locals("logger").(*utils.Logger)

	// Parse board ID from URL
	boardIDStr := c.Params("boardId")
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		logger.Warnw("Invalid board ID", "boardId", boardIDStr)
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	// Validate board access (both edit and public tokens are allowed for reading)
	token := c.Locals("token")
	if token != nil {
		// Token provided, validate it
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
		// No token provided, check if board exists
		if err := h.boardService.ValidateBoardExists(boardID); err != nil {
			if err == utils.ErrNotFound {
				return utils.SendNotFoundError(c, "Board not found")
			}
			logger.Errorw("Failed to validate board existence", "error", err)
			return utils.SendInternalError(c, "Failed to validate board", nil)
		}
	}

	// Get pages
	pages, err := h.pageService.GetPagesByBoard(boardID)
	if err != nil {
		logger.Errorw("Failed to get pages", "error", err)
		return utils.SendInternalError(c, "Failed to get pages", nil)
	}

	// Convert to response DTOs with elements
	pageResponses := make([]dto.PageWithElementsResponse, len(pages))
	for i, page := range pages {
		// Convert elements to response DTOs
		elementResponses := make([]dto.ElementResponse, len(page.Elements))
		for j, element := range page.Elements {
			elementResponses[j] = dto.ElementResponse{
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

		pageResponses[i] = dto.PageWithElementsResponse{
			ID:        page.ID,
			BoardID:   page.BoardID,
			Title:     page.Title,
			Date:      page.Date,
			OrderIdx:  page.OrderIdx,
			CreatedAt: page.CreatedAt,
			UpdatedAt: page.UpdatedAt,
			Elements:  elementResponses,
		}
	}

	response := dto.PagesWithElementsListResponse{
		Pages: pageResponses,
		Total: len(pageResponses),
	}

	return c.JSON(fiber.Map{"data": response})
}

// GetPage retrieves a single page with its elements
// GET /api/v1/boards/:boardId/pages/:pageId
func (h *PageHandler) GetPage(c *fiber.Ctx) error {
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

	// Validate board access
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

	// Get page with elements
	page, err := h.pageService.GetPageByID(pageID)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to get page", "error", err)
		return utils.SendInternalError(c, "Failed to get page", nil)
	}

	// Convert elements to response DTOs
	elementResponses := make([]dto.ElementResponse, len(page.Elements))
	for i, element := range page.Elements {
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

	response := dto.PageWithElementsResponse{
		ID:        page.ID,
		BoardID:   page.BoardID,
		Title:     page.Title,
		Date:      page.Date,
		OrderIdx:  page.OrderIdx,
		CreatedAt: page.CreatedAt,
		UpdatedAt: page.UpdatedAt,
		Elements:  elementResponses,
	}

	return c.JSON(fiber.Map{"data": response})
}

// UpdatePage updates a page's metadata
// PUT /api/v1/boards/:boardId/pages/:pageId
func (h *PageHandler) UpdatePage(c *fiber.Ctx) error {
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
	var req dto.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Warnw("Failed to parse request body", "error", err)
		return utils.SendValidationError(c, "Invalid request body", nil)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Update page
	page, err := h.pageService.UpdatePage(pageID, req.Title, req.Date, req.OrderIdx)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to update page", "error", err)
		return utils.SendInternalError(c, "Failed to update page", nil)
	}

	// Convert to response DTO
	response := dto.PageResponse{
		ID:        page.ID,
		BoardID:   page.BoardID,
		Title:     page.Title,
		Date:      page.Date,
		OrderIdx:  page.OrderIdx,
		CreatedAt: page.CreatedAt,
		UpdatedAt: page.UpdatedAt,
	}

	logger.Infow("Page updated successfully", "pageId", pageID)
	return c.JSON(fiber.Map{"data": response})
}

// DeletePage deletes a page and all its elements
// DELETE /api/v1/boards/:boardId/pages/:pageId
func (h *PageHandler) DeletePage(c *fiber.Ctx) error {
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

	// Delete page (cascade delete elements)
	if err := h.pageService.DeletePage(pageID); err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Page not found")
		}
		logger.Errorw("Failed to delete page", "error", err)
		return utils.SendInternalError(c, "Failed to delete page", nil)
	}

	logger.Infow("Page deleted successfully", "pageId", pageID)
	return c.SendStatus(fiber.StatusNoContent)
}
