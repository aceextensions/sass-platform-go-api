package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBRouter manages multiple database connection pools for multi-tenancy
type DBRouter struct {
	mu          sync.RWMutex
	tenantPools map[uuid.UUID]*pgxpool.Pool
	defaultPool *pgxpool.Pool
}

// NewDBRouter creates a new database router
func NewDBRouter(defaultPool *pgxpool.Pool) *DBRouter {
	return &DBRouter{
		tenantPools: make(map[uuid.UUID]*pgxpool.Pool),
		defaultPool: defaultPool,
	}
}

// GetExecutor returns a QueryExecutor (Pool or Tx) for the current context
func (r *DBRouter) GetExecutor(ctx context.Context) QueryExecutor {
	// 1. Check if there's a transaction in the context
	// (Transaction management will be handled in a follow-up step)

	// 2. Get TenantID from context
	tenantID, ok := GetTenantID(ctx)
	if !ok {
		return r.defaultPool
	}

	// 3. Resolve pool for tenant
	r.mu.RLock()
	pool, exists := r.tenantPools[tenantID]
	r.mu.RUnlock()

	if exists {
		return pool
	}

	// 4. Fallback to default pool
	return r.defaultPool
}

// RegisterTenantPool registers a specific connection pool for a tenant
func (r *DBRouter) RegisterTenantPool(tenantID uuid.UUID, connStr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if already registered
	if _, exists := r.tenantPools[tenantID]; exists {
		return nil
	}

	// Create new pool
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create pool for tenant %s: %w", tenantID, err)
	}

	// Ping to verify
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping tenant %s pool: %w", tenantID, err)
	}

	r.tenantPools[tenantID] = pool
	return nil
}

// UnregisterTenantPool closes and removes a tenant pool
func (r *DBRouter) UnregisterTenantPool(tenantID uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if pool, exists := r.tenantPools[tenantID]; exists {
		pool.Close()
		delete(r.tenantPools, tenantID)
	}
}

// GetExecutor is a global helper to get the correct executor for the context
func GetExecutor(ctx context.Context) QueryExecutor {
	if Router == nil {
		return MainPool
	}
	return Router.GetExecutor(ctx)
}
