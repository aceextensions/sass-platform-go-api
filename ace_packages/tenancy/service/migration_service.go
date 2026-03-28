package service

import (
	"context"
	"fmt"
	"log"

	"github.com/aceextension/core/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MigrationService handles the automated provisioning of tenant schemas
type MigrationService struct{}

func NewMigrationService() *MigrationService {
	return &MigrationService{}
}

// RunTenantMigration applies all current platform migrations to a specific tenant database
func (s *MigrationService) RunTenantMigration(ctx context.Context, dbURL string) error {
	// 1. Establish a temporary connection pool to the tenant's dedicated database
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return fmt.Errorf("invalid tenant database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to connect to tenant database for migration: %w", err)
	}
	defer pool.Close()

	// 2. Instantiate the core Migrator
	migrator := db.NewMigrator(pool)

	// 3. Fetch all embedded migrations from across the monorepo
	allMigrations, err := db.GetAllMigrations()
	if err != nil {
		return fmt.Errorf("failed to load embedded migrations: %w", err)
	}

	log.Printf("🏗️ Starting automatic migration for tenant DB (Total scripts: %d)\n", len(allMigrations))

	// 4. Execute the migration logic
	if err := migrator.ApplyMigrations(ctx, allMigrations); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}

	log.Println("✅ Tenant database migration completed successfully")
	return nil
}
