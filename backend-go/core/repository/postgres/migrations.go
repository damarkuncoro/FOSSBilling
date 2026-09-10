package postgres

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations checks for a schema_migrations table and executes missing .sql files
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	log.Println("🔄 Checking database migrations...")

	// 1. Create migrations table if not exists
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// 2. Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}
	sort.Strings(migrationFiles)

	// 3. Apply missing migrations
	for _, filename := range migrationFiles {
		var exists bool
		err := pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", filename).Scan(&exists)
		if err != nil {
			return err
		}

		if !exists {
			log.Printf("🚀 Applying migration: %s", filename)
			content, err := os.ReadFile(filepath.Join(migrationsDir, filename))
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", filename, err)
			}

			tx, err := pool.Begin(ctx)
			if err != nil {
				return err
			}

			if _, err := tx.Exec(ctx, string(content)); err != nil {
				tx.Rollback(ctx)
				return fmt.Errorf("failed to execute migration %s: %w", filename, err)
			}

			if _, err := tx.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", filename); err != nil {
				tx.Rollback(ctx)
				return err
			}

			if err := tx.Commit(ctx); err != nil {
				return err
			}
			log.Printf("✅ Migration successful: %s", filename)
		}
	}

	log.Println("   Database is up to date.")
	return nil
}
