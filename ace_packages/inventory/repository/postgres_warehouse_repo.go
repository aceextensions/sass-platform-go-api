package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/inventory/domain"
	"github.com/google/uuid"
)

type postgresWarehouseRepository struct{}

// NewPostgresWarehouseRepository creates a new postgres repository
func NewPostgresWarehouseRepository() WarehouseRepository {
	return &postgresWarehouseRepository{}
}

func (r *postgresWarehouseRepository) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `
		INSERT INTO inventory_warehouses (
			id, tenant_id, name, location, description, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := db.GetExecutor(ctx).Exec(ctx, query,
		warehouse.ID, warehouse.TenantID, warehouse.Name, warehouse.Location, warehouse.Description,
		warehouse.IsActive, warehouse.CreatedAt, warehouse.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create warehouse: %w", err)
	}
	return nil
}

func (r *postgresWarehouseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Warehouse, error) {
	query := `SELECT * FROM inventory_warehouses WHERE id = $1`
	var w domain.Warehouse
	err := db.GetExecutor(ctx).QueryRow(ctx, query, id).Scan(
		&w.ID, &w.TenantID, &w.Name, &w.Location, &w.Description,
		&w.IsActive, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse: %w", err)
	}
	return &w, nil
}

func (r *postgresWarehouseRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Warehouse, error) {
	query := `SELECT * FROM inventory_warehouses WHERE tenant_id = $1 ORDER BY created_at DESC`
	rows, err := db.GetExecutor(ctx).Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list warehouses: %w", err)
	}
	defer rows.Close()

	var warehouses []*domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		if err := rows.Scan(
			&w.ID, &w.TenantID, &w.Name, &w.Location, &w.Description,
			&w.IsActive, &w.CreatedAt, &w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &w)
	}
	return warehouses, nil
}

func (r *postgresWarehouseRepository) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `
		UPDATE inventory_warehouses 
		SET name = $1, location = $2, description = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`
	_, err := db.GetExecutor(ctx).Exec(ctx, query,
		warehouse.Name, warehouse.Location, warehouse.Description, warehouse.IsActive,
		warehouse.UpdatedAt, warehouse.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update warehouse: %w", err)
	}
	return nil
}

func (r *postgresWarehouseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM inventory_warehouses WHERE id = $1`
	_, err := db.GetExecutor(ctx).Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}
	return nil
}
