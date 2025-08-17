package handlers

import (
	"fmt"
	"os"

	"junk-journal-board/internal/dto"
	"junk-journal-board/internal/models"
	"junk-journal-board/internal/services"
	"junk-journal-board/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BoardHandler struct {
	boardService *services.BoardService
	validator    *validator.Validate
}

func NewBoardHandler(db *gorm.DB) *BoardHandler {
	return &BoardHandler{
		boardService: services.NewBoardService(db),
		validator:    validator.New(),
	}
}

// CreateBoard creates a new board
// POST /api/v1/boards
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
	var req dto.CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendValidationError(c, "Invalid request body", nil)
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Set default skin if not provided
	if req.Skin == "" {
		req.Skin = "default"
	}

	// Create board
	board, err := h.boardService.CreateBoard(req.Title, req.Skin)
	if err != nil {
		return utils.SendDatabaseError(c, "Failed to create board")
	}

	// Build URLs - Point to frontend, not backend
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // Default frontend URL
	}
	editURL := fmt.Sprintf("%s/board/%s/edit?edit_token=%s", frontendURL, board.ID, board.EditToken)
	publicURL := fmt.Sprintf("%s/board/%s/public?public_token=%s", frontendURL, board.ID, board.PublicToken)

	// Convert to response DTO (include edit token for board creation)
	response := dto.CreateBoardResponse{
		Board: dto.BoardWithTokensResponse{
			ID:          board.ID,
			Title:       board.Title,
			Skin:        board.Skin,
			EditToken:   board.EditToken,
			PublicToken: board.PublicToken,
			CreatedAt:   board.CreatedAt,
			UpdatedAt:   board.UpdatedAt,
		},
		EditURL:   editURL,
		PublicURL: publicURL,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": response})
}

// GetBoardByEditToken retrieves a board by edit token (with edit permissions)
// GET /api/v1/boards/edit/:editToken
func (h *BoardHandler) GetBoardByEditToken(c *fiber.Ctx) error {
	editTokenStr := c.Params("editToken")
	editToken, err := uuid.Parse(editTokenStr)
	if err != nil {
		return utils.SendValidationError(c, "Invalid edit token format", nil)
	}

	board, err := h.boardService.GetBoardByEditToken(editToken)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Board not found")
		}
		return utils.SendDatabaseError(c, "Failed to retrieve board")
	}

	// Convert to response DTO with tokens (for edit access)
	response := convertToBoardWithTokensResponse(board)
	return c.JSON(fiber.Map{"data": response})
}

// GetBoardByPublicToken retrieves a board by public token (read-only)
// GET /api/v1/boards/public/:publicToken
func (h *BoardHandler) GetBoardByPublicToken(c *fiber.Ctx) error {
	publicTokenStr := c.Params("publicToken")
	publicToken, err := uuid.Parse(publicTokenStr)
	if err != nil {
		return utils.SendValidationError(c, "Invalid public token format", nil)
	}

	board, err := h.boardService.GetBoardByPublicToken(publicToken)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Board not found")
		}
		return utils.SendDatabaseError(c, "Failed to retrieve board")
	}

	// Convert to response DTO (without edit token)
	response := convertToBoardResponse(board)
	return c.JSON(fiber.Map{"data": response})
}

// UpdateBoard updates a board's properties
// PUT /api/v1/boards/:boardId
func (h *BoardHandler) UpdateBoard(c *fiber.Ctx) error {
	boardIDStr := c.Params("boardId")
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		return utils.SendValidationError(c, "Invalid board ID format", nil)
	}

	var req dto.UpdateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendValidationError(c, "Invalid request body", nil)
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Update board
	board, err := h.boardService.UpdateBoard(boardID, req.Title, req.Skin)
	if err != nil {
		if err == utils.ErrNotFound {
			return utils.SendNotFoundError(c, "Board not found")
		}
		return utils.SendDatabaseError(c, "Failed to update board")
	}

	// Convert to response DTO (without edit token for security)
	response := convertToBoardResponse(board)
	return c.JSON(fiber.Map{"data": response})
}

// Helper functions

func convertToBoardResponse(board *models.Board) dto.BoardResponse {
	response := dto.BoardResponse{
		ID:          board.ID,
		Title:       board.Title,
		Skin:        board.Skin,
		PublicToken: board.PublicToken,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}

	// Convert pages if present
	if len(board.Pages) > 0 {
		response.Pages = make([]dto.PageResponse, len(board.Pages))
		for i, page := range board.Pages {
			response.Pages[i] = dto.PageResponse{
				ID:        page.ID,
				BoardID:   page.BoardID,
				Title:     page.Title,
				Date:      page.Date,
				OrderIdx:  page.OrderIdx,
				CreatedAt: page.CreatedAt,
				UpdatedAt: page.UpdatedAt,
			}
		}
	}

	return response
}

func convertToBoardWithTokensResponse(board *models.Board) dto.BoardWithTokensResponse {
	response := dto.BoardWithTokensResponse{
		ID:          board.ID,
		Title:       board.Title,
		Skin:        board.Skin,
		EditToken:   board.EditToken,
		PublicToken: board.PublicToken,
		CreatedAt:   board.CreatedAt,
		UpdatedAt:   board.UpdatedAt,
	}

	// Convert pages if present
	if len(board.Pages) > 0 {
		response.Pages = make([]dto.PageResponse, len(board.Pages))
		for i, page := range board.Pages {
			response.Pages[i] = dto.PageResponse{
				ID:        page.ID,
				BoardID:   page.BoardID,
				Title:     page.Title,
				Date:      page.Date,
				OrderIdx:  page.OrderIdx,
				CreatedAt: page.CreatedAt,
				UpdatedAt: page.UpdatedAt,
			}
		}
	}

	return response
}
