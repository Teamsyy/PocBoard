package routes

import (
	"junk-journal-board/internal/handlers"
	"junk-journal-board/internal/middleware"
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupUploadRoutes(api fiber.Router, db *gorm.DB) {
	logger, _ := utils.NewLogger()
	uploadHandler := handlers.NewUploadHandler(logger)

	// Upload routes - require edit token
	uploads := api.Group("/boards/:boardId/upload")

	// POST /api/v1/boards/:boardId/upload - Upload file (requires edit token)
	uploads.Post("/", middleware.TokenValidationMiddleware(), uploadHandler.UploadFile)
}
