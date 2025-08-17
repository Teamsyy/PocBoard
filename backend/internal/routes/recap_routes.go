package routes

import (
	"junk-journal-board/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRecapRoutes sets up all recap-related routes
func SetupRecapRoutes(api fiber.Router, db *gorm.DB) {
	recapHandler := handlers.NewRecapHandler(db)

	// Recap routes under /boards/:boardId/recap
	// Public route (no token required, but token can be provided for validation)
	api.Get("/boards/:boardId/recap", recapHandler.GetRecap) // GET /api/v1/boards/:boardId/recap
}
