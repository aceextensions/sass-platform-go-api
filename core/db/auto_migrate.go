package db

import (
	"context"
	"fmt"
	"time"
)

// AutoMigrate is a convenience function to run all pending migrations on startup.
// It uses the central MigrationRegistry to discover all module-level migrations.
func AutoMigrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if MainPool == nil {
		return fmt.Errorf("database pool not initialized")
	}

	// 1. Get all migrations across all modules (already sorted by dependency)
	migrations, err := GetAllMigrations()
	if err != nil {
		return fmt.Errorf("failed to gather migrations: %w", err)
	}

	// 2. Run the migrator
	migrator := NewMigrator(MainPool)
	if err := migrator.ApplyMigrations(ctx, migrations); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	return nil
}
