package routes

import (
	"junk-journal-board/internal/handlers"
	"junk-journal-board/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupPageRoutes sets up all page-related routes
func SetupPageRoutes(api fiber.Router, db *gorm.DB) {
	pageHandler := handlers.NewPageHandler(db)

	// Page routes under /boards/:boardId/pages
	pages := api.Group("/boards/:boardId/pages")

	// Public routes (no token required, but token can be provided for validation)
	pages.Get("/", pageHandler.GetPagesByBoard) // GET /api/v1/boards/:boardId/pages
	pages.Get("/:pageId", pageHandler.GetPage)  // GET /api/v1/boards/:boardId/pages/:pageId

	// Protected routes (require edit token)
	pages.Post("/", middleware.TokenValidationMiddleware(), pageHandler.CreatePage)          // POST /api/v1/boards/:boardId/pages
	pages.Put("/:pageId", middleware.TokenValidationMiddleware(), pageHandler.UpdatePage)    // PUT /api/v1/boards/:boardId/pages/:pageId
	pages.Delete("/:pageId", middleware.TokenValidationMiddleware(), pageHandler.DeletePage) // DELETE /api/v1/boards/:boardId/pages/:pageId
}
