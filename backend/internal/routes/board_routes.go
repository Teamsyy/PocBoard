package routes

import (
	"junk-journal-board/internal/handlers"
	"junk-journal-board/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupBoardRoutes(api fiber.Router, db *gorm.DB) {
	boardHandler := handlers.NewBoardHandler(db)

	// Board creation route (no token required)
	api.Post("/boards", boardHandler.CreateBoard) // POST /api/v1/boards

	// Board retrieval routes by token
	api.Get("/boards/edit/:editToken", boardHandler.GetBoardByEditToken)       // GET /api/v1/boards/edit/:editToken
	api.Get("/boards/public/:publicToken", boardHandler.GetBoardByPublicToken) // GET /api/v1/boards/public/:publicToken

	// Board update routes (require edit token)
	boards := api.Group("/boards/:boardId")
	boards.Put("/", middleware.TokenValidationMiddleware(), boardHandler.UpdateBoard) // PUT /api/v1/boards/:boardId
}
