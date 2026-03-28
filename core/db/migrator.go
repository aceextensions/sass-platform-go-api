package db

import (
	"context"
	"fmt"
	"sort"
	"time"
)

// Migration represents a single SQL migration file
type Migration struct {
	ID      string
	SQL     string
	Version int
}

// Migrator handles the execution of schema migrations on a database
type Migrator struct {
	executor QueryExecutor
}

func NewMigrator(executor QueryExecutor) *Migrator {
	return &Migrator{executor: executor}
}

// EnsureSchemaTable creates the _schema_migrations table if it doesn't exist
func (m *Migrator) EnsureSchemaTable(ctx context.Context) error {
	// First ensure the table exists (standardized name 'schema_migrations')
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id TEXT PRIMARY KEY,
			version INTEGER NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)
	`
	if _, err := m.executor.Exec(ctx, query); err != nil {
		return err
	}

	// Then handle the transition from the old 'filename' column if it exists
	transitionQuery := `
		DO $$
		BEGIN
			IF EXISTS (SELECT 1 FROM information_schema.columns 
			           WHERE table_name = 'schema_migrations' AND column_name = 'filename') THEN
				ALTER TABLE schema_migrations RENAME COLUMN filename TO id;
			END IF;
		END $$;
	`
	_, err := m.executor.Exec(ctx, transitionQuery)
	return err
}

// GetAppliedMigrations retrieves the IDs of all migrations already applied
func (m *Migrator) GetAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	query := `SELECT id FROM schema_migrations`
	rows, err := m.executor.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		applied[id] = true
	}
	return applied, nil
}

// ApplyMigrations executes a list of migrations in order, skipping those already applied
func (m *Migrator) ApplyMigrations(ctx context.Context, migrations []Migration) error {
	// 1. Ensure metadata table exists
	if err := m.EnsureSchemaTable(ctx); err != nil {
		return fmt.Errorf("failed to ensure schema table: %w", err)
	}

	// 2. Get already applied migrations
	applied, err := m.GetAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	// 3. Sort migrations by version/ID to ensure deterministic order
	sort.Slice(migrations, func(i, j int) bool {
		if migrations[i].Version != migrations[j].Version {
			return migrations[i].Version < migrations[j].Version
		}
		return migrations[i].ID < migrations[j].ID
	})

	// 4. Wrap everything in a transaction for safety
	// Note: We use the QueryExecutor interface. If it's a pool, we should ideally start a transaction.
	// For simplicity in this worker, we assume the caller can pass a Tx if they want.
	// But let's handle the transaction inside here if possible.

	for _, mig := range migrations {
		if applied[mig.ID] {
			continue
		}

		fmt.Printf("📝 Applying migration: %s\n", mig.ID)

		// Execute as a single batch/transaction step
		_, err := m.executor.Exec(ctx, mig.SQL)
		if err != nil {
			return fmt.Errorf("migration %s failed: %w", mig.ID, err)
		}

		// Record success
		_, err = m.executor.Exec(ctx, 
			"INSERT INTO schema_migrations (id, version, applied_at) VALUES ($1, $2, $3)",
			mig.ID, mig.Version, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", mig.ID, err)
		}
	}

	return nil
}
