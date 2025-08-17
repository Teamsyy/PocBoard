package main

import (
	"log"
	"os"

	"junk-journal-board/internal/config"
	"junk-journal-board/internal/middleware"
	"junk-journal-board/internal/routes"
	"junk-journal-board/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize structured logger
	logger, err := utils.NewLogger()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	// Connect to database
	db, err := config.ConnectDatabase()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Initialize database (run migrations)
	if err := config.InitializeDatabase(db, logger.Logger); err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Create Fiber app with custom error handler
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(logger),
	})

	// Core middleware
	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			return utils.GenerateUUID().String()
		},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Request-ID",
	}))

	app.Use(middleware.LoggingMiddleware(logger))

	// Health check endpoints
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/health/db", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB()
		if err != nil {
			return utils.SendDatabaseError(c, "Failed to get database instance")
		}

		if err := sqlDB.Ping(); err != nil {
			return utils.SendDatabaseError(c, "Database connection failed")
		}

		return c.JSON(fiber.Map{"status": "ok", "database": "connected"})
	})

	// API routes
	api := app.Group("/api/v1")

	// Add optional token middleware to API routes for token extraction
	api.Use(middleware.OptionalTokenMiddleware())

	// Setup board routes
	routes.SetupBoardRoutes(api, db)

	// Setup page routes
	routes.SetupPageRoutes(api, db)

	// Setup element routes
	routes.SetupElementRoutes(api, db)

	// Setup upload routes
	routes.SetupUploadRoutes(api, db)

	// Setup recap routes
	routes.SetupRecapRoutes(api, db)

	// Static file serving for uploads
	app.Static("/uploads", "./uploads")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Server starting", zap.String("port", port))
	log.Fatal(app.Listen(":" + port))
}
