package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RLSManager handles Row-Level Security operations
type RLSManager struct {
	pool *pgxpool.Pool
}

// NewRLSManager creates a new RLS manager
func NewRLSManager(pool *pgxpool.Pool) *RLSManager {
	return &RLSManager{pool: pool}
}

// SetTenantContext sets the current tenant ID in the database session
// This is used by RLS policies to filter data
func (m *RLSManager) SetTenantContext(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	_, err := tx.Exec(ctx, "SET LOCAL app.current_tenant_id = $1", tenantID.String())
	if err != nil {
		return fmt.Errorf("failed to set tenant context: %w", err)
	}
	return nil
}

// ClearTenantContext clears the tenant context (useful for super admin operations)
func (m *RLSManager) ClearTenantContext(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, "RESET app.current_tenant_id")
	if err != nil {
		return fmt.Errorf("failed to clear tenant context: %w", err)
	}
	return nil
}

// EnableRLS enables Row-Level Security on a table
func (m *RLSManager) EnableRLS(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("ALTER TABLE %s ENABLE ROW LEVEL SECURITY", tableName)
	_, err := m.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to enable RLS on %s: %w", tableName, err)
	}
	return nil
}

// DisableRLS disables Row-Level Security on a table
func (m *RLSManager) DisableRLS(ctx context.Context, tableName string) error {
	query := fmt.Sprintf("ALTER TABLE %s DISABLE ROW LEVEL SECURITY", tableName)
	_, err := m.pool.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to disable RLS on %s: %w", tableName, err)
	}
	return nil
}

// CreateTenantPolicy creates a standard tenant isolation policy
func (m *RLSManager) CreateTenantPolicy(ctx context.Context, tableName string) error {
	// Drop existing policy if it exists
	dropQuery := fmt.Sprintf("DROP POLICY IF EXISTS tenant_isolation ON %s", tableName)
	_, _ = m.pool.Exec(ctx, dropQuery)

	// Create new policy
	policyQuery := fmt.Sprintf(`
		CREATE POLICY tenant_isolation ON %s
		USING (tenant_id = current_setting('app.current_tenant_id', true)::uuid)
	`, tableName)

	_, err := m.pool.Exec(ctx, policyQuery)
	if err != nil {
		return fmt.Errorf("failed to create tenant policy on %s: %w", tableName, err)
	}
	return nil
}

// CreateSuperAdminBypassPolicy creates a policy that allows super admins to bypass RLS
func (m *RLSManager) CreateSuperAdminBypassPolicy(ctx context.Context, tableName string) error {
	// Drop existing policy if it exists
	dropQuery := fmt.Sprintf("DROP POLICY IF EXISTS super_admin_bypass ON %s", tableName)
	_, _ = m.pool.Exec(ctx, dropQuery)

	// Create bypass policy for super admins
	policyQuery := fmt.Sprintf(`
		CREATE POLICY super_admin_bypass ON %s
		USING (current_setting('app.is_super_admin', true)::boolean = true)
	`, tableName)

	_, err := m.pool.Exec(ctx, policyQuery)
	if err != nil {
		return fmt.Errorf("failed to create super admin bypass policy on %s: %w", tableName, err)
	}
	return nil
}

// SetupTableRLS enables RLS and creates policies for a table
func (m *RLSManager) SetupTableRLS(ctx context.Context, tableName string, includeSuperAdminBypass bool) error {
	// Enable RLS
	if err := m.EnableRLS(ctx, tableName); err != nil {
		return err
	}

	// Create tenant isolation policy
	if err := m.CreateTenantPolicy(ctx, tableName); err != nil {
		return err
	}

	// Optionally create super admin bypass policy
	if includeSuperAdminBypass {
		if err := m.CreateSuperAdminBypassPolicy(ctx, tableName); err != nil {
			return err
		}
	}

	return nil
}

// ExecuteWithTenantContext executes a function within a transaction with tenant context set
func (m *RLSManager) ExecuteWithTenantContext(ctx context.Context, tenantID uuid.UUID, fn func(pgx.Tx) error) error {
	return BeginFunc(ctx, func(tx pgx.Tx) error {
		// Set tenant context
		if err := m.SetTenantContext(ctx, tx, tenantID); err != nil {
			return err
		}

		// Execute function
		return fn(tx)
	})
}

// ExecuteWithSuperAdminContext executes a function within a transaction with super admin privileges
func (m *RLSManager) ExecuteWithSuperAdminContext(ctx context.Context, fn func(pgx.Tx) error) error {
	return BeginFunc(ctx, func(tx pgx.Tx) error {
		// Set super admin flag
		_, err := tx.Exec(ctx, "SET LOCAL app.is_super_admin = true")
		if err != nil {
			return fmt.Errorf("failed to set super admin context: %w", err)
		}

		// Execute function
		return fn(tx)
	})
}

// Global RLS manager instance
var RLS *RLSManager

// InitRLS initializes the global RLS manager
func InitRLS() {
	RLS = NewRLSManager(MainPool)
}
