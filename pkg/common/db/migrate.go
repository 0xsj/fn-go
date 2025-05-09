// pkg/common/db/migrate.go
package db

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/0xsj/fn-go/pkg/common/log"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	SQL         string
	Timestamp   time.Time
}

// Migrator handles database migrations
type Migrator struct {
	db     DB
	logger log.Logger
}

// NewMigrator creates a new migrator
func NewMigrator(db DB, logger log.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// ensureMigrationsTable creates the migrations table if it doesn't exist
func (m *Migrator) ensureMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			version INT PRIMARY KEY,
			description VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
	
	_, err := m.db.Execute(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	
	return nil
}

// getAppliedMigrations retrieves the list of applied migrations
func (m *Migrator) getAppliedMigrations(ctx context.Context) (map[int]bool, error) {
	query := `SELECT version FROM migrations ORDER BY version ASC`
	
	rows, err := m.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()
	
	applied := make(map[int]bool)
	
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating migrations: %w", err)
	}
	
	return applied, nil
}

// LoadMigrationsFromDir loads migrations from a directory
func (m *Migrator) LoadMigrationsFromDir(dir string) ([]Migration, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	var migrations []Migration
	
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		
		// Parse filename: VXXX_description.sql
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) != 2 || !strings.HasPrefix(parts[0], "V") {
			m.logger.With("file", file.Name()).Warn("Skipping file with invalid migration format")
			continue
		}
		
		// Parse version
		var version int
		if _, err := fmt.Sscanf(parts[0], "V%d", &version); err != nil {
			m.logger.With("file", file.Name()).Warn("Skipping file with invalid version format")
			continue
		}
		
		// Parse description
		description := strings.TrimSuffix(parts[1], ".sql")
		
		// Read SQL content
		path := filepath.Join(dir, file.Name())
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}
		
		migrations = append(migrations, Migration{
			Version:     version,
			Description: description,
			SQL:         string(content),
			Timestamp:   file.ModTime(),
		})
	}
	
	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	
	return migrations, nil
}

// MigrateUp applies pending migrations
func (m *Migrator) MigrateUp(ctx context.Context, migrations []Migration) error {
	// Ensure migrations table exists
	if err := m.ensureMigrationsTable(ctx); err != nil {
		return err
	}
	
	// Get applied migrations
	applied, err := m.getAppliedMigrations(ctx)
	if err != nil {
		return err
	}
	
	m.logger.With("applied_count", len(applied)).Info("Found applied migrations")
	
	// Apply pending migrations
	for _, migration := range migrations {
		if applied[migration.Version] {
			m.logger.With("version", migration.Version).
				With("description", migration.Description).
				Debug("Skipping already applied migration")
			continue
		}
		
		m.logger.With("version", migration.Version).
			With("description", migration.Description).
			Info("Applying migration")
		
		// Run in a transaction
		err := WithTransaction(ctx, m.db, func(txCtx context.Context) error {
			// Apply migration
			_, err := m.db.Execute(txCtx, migration.SQL)
			if err != nil {
				return fmt.Errorf("failed to apply migration V%d: %w", migration.Version, err)
			}
			
			// Record migration
			_, err = m.db.Execute(txCtx, 
				"INSERT INTO migrations (version, description) VALUES (?, ?)",
				migration.Version, migration.Description)
			if err != nil {
				return fmt.Errorf("failed to record migration V%d: %w", migration.Version, err)
			}
			
			return nil
		})
		
		if err != nil {
			return err
		}
		
		m.logger.With("version", migration.Version).
			With("description", migration.Description).
			Info("Successfully applied migration")
	}
	
	return nil
}

// CreateMigration creates a new migration file
func (m *Migrator) CreateMigration(dir, description string) (string, error) {
	// Ensure migrations directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create migrations directory: %w", err)
	}
	
	// Get existing migrations to determine next version
	migrations, err := m.LoadMigrationsFromDir(dir)
	if err != nil {
		return "", err
	}
	
	nextVersion := 1
	if len(migrations) > 0 {
		nextVersion = migrations[len(migrations)-1].Version + 1
	}
	
	// Create filename
	safeName := strings.ReplaceAll(description, " ", "_")
	filename := fmt.Sprintf("V%03d_%s.sql", nextVersion, safeName)
	path := filepath.Join(dir, filename)
	
	// Create file with template
	template := fmt.Sprintf(`-- Migration: %s
-- Version: %d
-- Created: %s

`, description, nextVersion, time.Now().Format(time.RFC3339))
	
	if err := ioutil.WriteFile(path, []byte(template), 0644); err != nil {
		return "", fmt.Errorf("failed to create migration file: %w", err)
	}
	
	return path, nil
}