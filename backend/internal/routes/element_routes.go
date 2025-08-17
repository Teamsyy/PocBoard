package routes

import (
	"junk-journal-board/internal/handlers"
	"junk-journal-board/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupElementRoutes(api fiber.Router, db *gorm.DB) {
	elementHandler := handlers.NewElementHandler(db)

	// All element routes require authentication
	elements := api.Group("/boards/:boardId/pages/:pageId/elements")

	// Create element (requires edit token)
	elements.Post("/", middleware.TokenValidationMiddleware(), elementHandler.CreateElement)

	// Get elements for a page (allows both edit and public tokens, or no token for public access)
	elements.Get("/", middleware.OptionalTokenMiddleware(), elementHandler.GetElementsByPage)

	// Update element (requires edit token)
	elements.Put("/:elementId", middleware.TokenValidationMiddleware(), elementHandler.UpdateElement)

	// Delete element (requires edit token)
	elements.Delete("/:elementId", middleware.TokenValidationMiddleware(), elementHandler.DeleteElement)

	// Batch reorder elements (requires edit token)
	elements.Put("/reorder", middleware.TokenValidationMiddleware(), elementHandler.ReorderElements)
}
