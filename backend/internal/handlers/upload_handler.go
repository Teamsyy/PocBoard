package handlers

import (
	"junk-journal-board/internal/dto"
	"junk-journal-board/internal/services"
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UploadHandler struct {
	uploadService *services.UploadService
	logger        *utils.Logger
}

func NewUploadHandler(logger *utils.Logger) *UploadHandler {
	return &UploadHandler{
		uploadService: services.NewUploadService(),
		logger:        logger,
	}
}

// UploadFile handles file upload for a specific board
func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	// Get board ID from URL params
	boardIDStr := c.Params("boardId")
	if boardIDStr == "" {
		h.logger.Warn("Upload attempt without board ID")
		return utils.SendBadRequestError(c, "Board ID is required")
	}

	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		h.logger.Warn("Invalid board ID format", zap.String("boardId", boardIDStr))
		return utils.SendBadRequestError(c, "Invalid board ID format")
	}

	// Get the uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.logger.Warn("No file in upload request", zap.Error(err))
		return utils.SendBadRequestError(c, "No file provided")
	}

	// Validate the file
	if err := h.uploadService.ValidateFile(fileHeader); err != nil {
		h.logger.Warn("File validation failed",
			zap.String("filename", fileHeader.Filename),
			zap.Int64("size", fileHeader.Size),
			zap.Error(err))
		return utils.SendValidationError(c, err.Error(), nil)
	}

	// Save the file
	publicURL, err := h.uploadService.SaveFile(fileHeader, boardID)
	if err != nil {
		h.logger.Error("Failed to save uploaded file",
			zap.String("filename", fileHeader.Filename),
			zap.String("boardId", boardID.String()),
			zap.Error(err))
		return utils.SendInternalError(c, "Failed to save file", nil)
	}

	// Create response
	response := dto.UploadResponse{
		URL:      publicURL,
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		MimeType: h.uploadService.GetMimeType(fileHeader.Filename),
	}

	h.logger.Info("File uploaded successfully",
		zap.String("filename", fileHeader.Filename),
		zap.String("boardId", boardID.String()),
		zap.String("url", publicURL))

	return c.Status(fiber.StatusCreated).JSON(response)
}
