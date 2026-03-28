package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/tenancy/domain"
	"github.com/google/uuid"
)

type PostgresTenantRepository struct{}

func NewPostgresTenantRepository() *PostgresTenantRepository {
	return &PostgresTenantRepository{}
}

// GetByID retrieves a tenant by ID
func (r *PostgresTenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	query := `
		SELECT id, name, email, phone, status, database_url, schema_name, metadata, settings, created_at, updated_at
		FROM tenants
		WHERE id = $1
	`
	
	var tenant domain.Tenant
	err := db.GetExecutor(ctx).QueryRow(ctx, query, id).Scan(
		&tenant.ID, &tenant.Name, &tenant.Email, &tenant.Phone, 
		&tenant.Status, &tenant.DatabaseURL, &tenant.SchemaName, 
		&tenant.Metadata, &tenant.Settings, &tenant.CreatedAt, &tenant.UpdatedAt,
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	
	return &tenant, nil
}

// GetAll retrieves all active tenants
func (r *PostgresTenantRepository) GetAll(ctx context.Context) ([]*domain.Tenant, error) {
	query := `
		SELECT id, name, email, phone, status, database_url, schema_name, metadata, settings, created_at, updated_at
		FROM tenants
		WHERE status = 'active'
	`
	
	rows, err := db.GetExecutor(ctx).Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}
	defer rows.Close()
	
	var tenants []*domain.Tenant
	for rows.Next() {
		var tenant domain.Tenant
		err := rows.Scan(
			&tenant.ID, &tenant.Name, &tenant.Email, &tenant.Phone, 
			&tenant.Status, &tenant.DatabaseURL, &tenant.SchemaName, 
			&tenant.Metadata, &tenant.Settings, &tenant.CreatedAt, &tenant.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tenant: %w", err)
		}
		tenants = append(tenants, &tenant)
	}
	
	return tenants, nil
}

// Update updates a tenant profile
func (r *PostgresTenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	query := `
		UPDATE tenants 
		SET name = $2, email = $3, phone = $4, database_url = $5, schema_name = $6, metadata = $7, settings = $8, updated_at = NOW()
		WHERE id = $1
	`
	_, err := db.GetExecutor(ctx).Exec(ctx, query,
		tenant.ID, tenant.Name, tenant.Email, tenant.Phone,
		tenant.DatabaseURL, tenant.SchemaName, tenant.Metadata, tenant.Settings,
	)
	if err != nil {
		return fmt.Errorf("failed to update tenant: %w", err)
	}
	return nil
}

// UpdateSettings updates only the settings JSONB field
func (r *PostgresTenantRepository) UpdateSettings(ctx context.Context, id uuid.UUID, settings interface{}) error {
	query := `UPDATE tenants SET settings = $2, updated_at = NOW() WHERE id = $1`
	_, err := db.GetExecutor(ctx).Exec(ctx, query, id, settings)
	if err != nil {
		return fmt.Errorf("failed to update tenant settings: %w", err)
	}
	return nil
}
