package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, db *pgxpool.Pool, dir string, logger *slog.Logger) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return fmt.Errorf("find migration files: %w", err)
	}
	sort.Strings(files)

	if len(files) == 0 {
		logger.Warn("no migration files found", "dir", dir)
		return nil
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migrations transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	for _, file := range files {
		version := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`, version).Scan(&exists); err != nil {
			return fmt.Errorf("check migration %s: %w", version, err)
		}
		if exists {
			logger.Debug("migration already applied", "version", version)
			continue
		}

		sql, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", version, err)
		}

		if _, err := tx.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("apply migration %s: %w", version, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		logger.Info("migration applied", "version", version)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migrations: %w", err)
	}

	return nil
}
