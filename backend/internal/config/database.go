package config

import (
	"fmt"
	"os"

	"junk-journal-board/internal/migrations"
	"junk-journal-board/internal/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "junk_journal"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	// Configure GORM logger based on environment
	var gormLogger logger.Interface
	if os.Getenv("GO_ENV") == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// InitializeDatabase runs migrations and sets up the database
func InitializeDatabase(db *gorm.DB, zapLogger *zap.Logger) error {
	// Run SQL migrations first
	migrator := migrations.NewMigrator(db, zapLogger)
	if err := migrator.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Auto-migrate GORM models to ensure schema is up to date
	if err := db.AutoMigrate(
		&models.Board{},
		&models.Page{},
		&models.Element{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	zapLogger.Info("Database initialized successfully")
	return nil
}
