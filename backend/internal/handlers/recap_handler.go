package handlers

import (
	"time"

	"junk-journal-board/internal/dto"
	"junk-journal-board/internal/services"
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecapHandler struct {
	recapService *services.RecapService
	boardService *services.BoardService
}

func NewRecapHandler(db *gorm.DB) *RecapHandler {
	return &RecapHandler{
		recapService: services.NewRecapService(db),
		boardService: services.NewBoardService(db),
	}
}

// GetRecap retrieves recap data for a board with optional date filtering
// GET /api/v1/boards/:boardId/recap?filter=day|week|month&date=2024-01-15
func (h *RecapHandler) GetRecap(c *fiber.Ctx) error {
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

	// Parse query parameters
	var req dto.RecapRequest
	if err := c.QueryParser(&req); err != nil {
		logger.Warnw("Failed to parse query parameters", "error", err)
		return utils.SendValidationError(c, "Invalid query parameters", nil)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Set default filter if not provided
	if req.Filter == "" {
		req.Filter = "day"
	}

	// Parse date if provided
	var date *time.Time
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			logger.Warnw("Invalid date format", "date", req.Date)
			return utils.SendValidationError(c, "Invalid date format. Use YYYY-MM-DD", nil)
		}
		date = &parsedDate
	}

	// Get recap data
	recap, err := h.recapService.GetRecapData(boardID, req.Filter, date)
	if err != nil {
		logger.Errorw("Failed to get recap data", "error", err)
		return utils.SendInternalError(c, "Failed to get recap data", nil)
	}

	logger.Infow("Recap data retrieved successfully",
		"boardId", boardID,
		"filter", req.Filter,
		"pageCount", recap.PageCount,
		"elementCount", recap.ElementCount)

	return c.JSON(recap)
}
