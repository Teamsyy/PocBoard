package migrations

import (
	"embed"
	"fmt"
	"sort"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

//go:embed *.sql
var migrationFiles embed.FS

// Migration represents a database migration
type Migration struct {
	Version string
	Name    string
	SQL     string
}

// Migrator handles database migrations
type Migrator struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *gorm.DB, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// RunMigrations executes all pending migrations
func (m *Migrator) RunMigrations() error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Load all migration files
	migrations, err := m.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range migrations {
		if _, applied := appliedMigrations[migration.Version]; !applied {
			if err := m.runMigration(migration); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.Version, err)
			}
			m.logger.Info("Applied migration", zap.String("version", migration.Version), zap.String("name", migration.Name))
		}
	}

	m.logger.Info("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`
	return m.db.Exec(query).Error
}

// loadMigrations loads all migration files from the embedded filesystem
func (m *Migrator) loadMigrations() ([]Migration, error) {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := migrationFiles.ReadFile(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		// Extract version and name from filename (e.g., "001_create_boards_table.sql")
		name := entry.Name()
		version := strings.Split(name, "_")[0]
		displayName := strings.TrimSuffix(strings.TrimPrefix(name, version+"_"), ".sql")

		migrations = append(migrations, Migration{
			Version: version,
			Name:    displayName,
			SQL:     string(content),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// getAppliedMigrations returns a map of applied migration versions
func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var versions []string
	err := m.db.Raw("SELECT version FROM schema_migrations").Scan(&versions).Error
	if err != nil {
		return nil, err
	}

	applied := make(map[string]bool)
	for _, version := range versions {
		applied[version] = true
	}

	return applied, nil
}

// runMigration executes a single migration
func (m *Migrator) runMigration(migration Migration) error {
	// Start transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Execute migration SQL
	if err := tx.Exec(migration.SQL).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Record migration as applied
	if err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration.Version).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}
