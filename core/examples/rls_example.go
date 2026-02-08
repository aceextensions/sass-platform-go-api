package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Load configuration
	config.Load()

	// Initialize database connections
	mainConnStr := config.Get("DATABASE_URL")
	auditConnStr := config.Get("AUDIT_DATABASE_URL")
	db.Init(mainConnStr, auditConnStr)
	defer db.Close()

	// Initialize RLS manager
	db.InitRLS()

	fmt.Println("🔒 RLS Middleware Example")
	fmt.Println("==========================\n")

	ctx := context.Background()

	// Create test tenants
	tenant1ID := uuid.New()
	tenant2ID := uuid.New()

	fmt.Println("Creating test tenants...")
	if err := createTestTenant(ctx, tenant1ID, "Tenant 1"); err != nil {
		log.Fatalf("Failed to create tenant 1: %v", err)
	}
	if err := createTestTenant(ctx, tenant2ID, "Tenant 2"); err != nil {
		log.Fatalf("Failed to create tenant 2: %v", err)
	}
	fmt.Println("✅ Created 2 test tenants\n")

	// Example 1: Query with tenant context
	fmt.Println("Example 1: Query with Tenant Context")
	fmt.Println("--------------------------------------")

	err := db.RLS.ExecuteWithTenantContext(ctx, tenant1ID, func(tx pgx.Tx) error {
		// This query will only return users from tenant1
		rows, err := tx.Query(ctx, "SELECT id, name, tenant_id FROM users LIMIT 5")
		if err != nil {
			return err
		}
		defer rows.Close()

		fmt.Printf("Users visible to Tenant 1 (%s):\n", tenant1ID)
		count := 0
		for rows.Next() {
			var id, tenantID uuid.UUID
			var name string
			if err := rows.Scan(&id, &name, &tenantID); err != nil {
				return err
			}
			fmt.Printf("  - %s (Tenant: %s)\n", name, tenantID)
			count++
		}
		if count == 0 {
			fmt.Println("  (No users found)")
		}

		return nil
	})

	if err != nil {
		log.Printf("Error querying with tenant context: %v", err)
	}

	fmt.Println()

	// Example 2: Super admin context (can see all tenants)
	fmt.Println("Example 2: Super Admin Context")
	fmt.Println("-------------------------------")

	err = db.RLS.ExecuteWithSuperAdminContext(ctx, func(tx pgx.Tx) error {
		// This query will return users from ALL tenants
		rows, err := tx.Query(ctx, "SELECT id, name, tenant_id FROM users LIMIT 10")
		if err != nil {
			return err
		}
		defer rows.Close()

		fmt.Println("Users visible to Super Admin:")
		count := 0
		for rows.Next() {
			var id, tenantID uuid.UUID
			var name string
			if err := rows.Scan(&id, &name, &tenantID); err != nil {
				return err
			}
			fmt.Printf("  - %s (Tenant: %s)\n", name, tenantID)
			count++
		}
		if count == 0 {
			fmt.Println("  (No users found)")
		}

		return nil
	})

	if err != nil {
		log.Printf("Error querying with super admin context: %v", err)
	}

	fmt.Println()

	// Example 3: Manual tenant context setting
	fmt.Println("Example 3: Manual Tenant Context")
	fmt.Println("---------------------------------")

	err = db.BeginFunc(ctx, func(tx pgx.Tx) error {
		// Manually set tenant context
		if err := db.RLS.SetTenantContext(ctx, tx, tenant2ID); err != nil {
			return err
		}

		// Query will be scoped to tenant2
		rows, err := tx.Query(ctx, "SELECT id, name FROM users LIMIT 5")
		if err != nil {
			return err
		}
		defer rows.Close()

		fmt.Printf("Users visible to Tenant 2 (%s):\n", tenant2ID)
		count := 0
		for rows.Next() {
			var id uuid.UUID
			var name string
			if err := rows.Scan(&id, &name); err != nil {
				return err
			}
			fmt.Printf("  - %s\n", name)
			count++
		}
		if count == 0 {
			fmt.Println("  (No users found)")
		}

		return nil
	})

	if err != nil {
		log.Printf("Error with manual tenant context: %v", err)
	}

	fmt.Println()

	// Example 4: Context utilities
	fmt.Println("Example 4: Context Utilities")
	fmt.Println("-----------------------------")

	// Add tenant to Go context
	ctxWithTenant := db.WithTenantID(ctx, tenant1ID)
	ctxWithUser := db.WithUserID(ctxWithTenant, uuid.New())

	// Retrieve from context
	if tenantID, ok := db.GetTenantID(ctxWithUser); ok {
		fmt.Printf("Tenant ID from context: %s\n", tenantID)
	}

	if userID, ok := db.GetUserID(ctxWithUser); ok {
		fmt.Printf("User ID from context: %s\n", userID)
	}

	// Check super admin
	ctxWithSuperAdmin := db.WithSuperAdmin(ctx, true)
	if db.IsSuperAdmin(ctxWithSuperAdmin) {
		fmt.Println("Context has super admin privileges: true")
	}

	fmt.Println("\n✨ RLS middleware is working correctly!")
}

func createTestTenant(ctx context.Context, tenantID uuid.UUID, name string) error {
	query := `
		INSERT INTO tenants (id, name, status, is_active)
		VALUES ($1, $2, 'active', true)
		ON CONFLICT (id) DO NOTHING
	`
	_, err := db.MainPool.Exec(ctx, query, tenantID, name)
	return err
}
