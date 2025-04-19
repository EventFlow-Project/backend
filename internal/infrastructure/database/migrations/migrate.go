package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/EventFlow-Project/backend/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	migrationsDir, err := getMigrationsDir()
	if err != nil {
		return fmt.Errorf("failed to get migrations directory: %w", err)
	}

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsDir),
		dsn,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func getMigrationsDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	migrationsDir := filepath.Join(currentDir, "migrations")

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return "", fmt.Errorf("failed to read migrations directory: %w", err)
	}

	hasSQLFiles := false
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			hasSQLFiles = true
			break
		}
	}

	if !hasSQLFiles {
		return "", fmt.Errorf("no SQL files found in migrations directory")
	}

	return migrationsDir, nil
}
