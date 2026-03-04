package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/inventory/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresInventoryRepository struct {
	db *pgxpool.Pool
}

// NewPostgresInventoryRepository creates a new postgres repository
func NewPostgresInventoryRepository(db *pgxpool.Pool) InventoryRepository {
	return &postgresInventoryRepository{db: db}
}

func (r *postgresInventoryRepository) GetStock(ctx context.Context, warehouseID, productID uuid.UUID) (*domain.InventoryItem, error) {
	query := `
		SELECT * FROM inventory_items 
		WHERE warehouse_id = $1 AND product_id = $2
	`
	var i domain.InventoryItem
	err := r.db.QueryRow(ctx, query, warehouseID, productID).Scan(
		&i.ID, &i.TenantID, &i.WarehouseID, &i.ProductID, &i.Quantity,
		&i.ReorderLevel, &i.LastStock, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Return nil if record doesn't exist yet
		}
		return nil, fmt.Errorf("failed to get stock: %w", err)
	}
	return &i, nil
}

func (r *postgresInventoryRepository) GetStockByWarehouse(ctx context.Context, warehouseID uuid.UUID) ([]*domain.InventoryItem, error) {
	query := `SELECT * FROM inventory_items WHERE warehouse_id = $1`
	rows, err := r.db.Query(ctx, query, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list stock: %w", err)
	}
	defer rows.Close()

	var items []*domain.InventoryItem
	for rows.Next() {
		var i domain.InventoryItem
		if err := rows.Scan(
			&i.ID, &i.TenantID, &i.WarehouseID, &i.ProductID, &i.Quantity,
			&i.ReorderLevel, &i.LastStock, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	return items, nil
}

func (r *postgresInventoryRepository) CreateStock(ctx context.Context, item *domain.InventoryItem) error {
	query := `
		INSERT INTO inventory_items (
			id, tenant_id, warehouse_id, product_id, quantity, reorder_level, last_restocked, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		item.ID, item.TenantID, item.WarehouseID, item.ProductID, item.Quantity,
		item.ReorderLevel, item.LastStock, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create stock record: %w", err)
	}
	return nil
}

func (r *postgresInventoryRepository) UpdateStock(ctx context.Context, item *domain.InventoryItem) error {
	query := `
		UPDATE inventory_items 
		SET quantity = $1, reorder_level = $2, last_restocked = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query,
		item.Quantity, item.ReorderLevel, item.LastStock, item.UpdatedAt, item.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}
	return nil
}

func (r *postgresInventoryRepository) LogTransaction(ctx context.Context, t *domain.StockTransaction) error {
	query := `
		INSERT INTO inventory_transactions (
			id, tenant_id, warehouse_id, product_id, type, quantity, direction,
			reference_id, reference_type, notes, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(ctx, query,
		t.ID, t.TenantID, t.WarehouseID, t.ProductID, t.Type, t.Quantity, t.Direction,
		t.ReferenceID, t.ReferenceType, t.Notes, t.CreatedBy, t.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to log stock transaction: %w", err)
	}
	return nil
}

func (r *postgresInventoryRepository) GetTransactions(ctx context.Context, warehouseID, productID uuid.UUID, limit, offset int) ([]*domain.StockTransaction, error) {
	query := `
		SELECT * FROM inventory_transactions 
		WHERE warehouse_id = $1 AND product_id = $2 
		ORDER BY created_at DESC 
		LIMIT $3 OFFSET $4
	`
	rows, err := r.db.Query(ctx, query, warehouseID, productID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.StockTransaction
	for rows.Next() {
		var t domain.StockTransaction
		if err := rows.Scan(
			&t.ID, &t.TenantID, &t.WarehouseID, &t.ProductID, &t.Type, &t.Quantity, &t.Direction,
			&t.ReferenceID, &t.ReferenceType, &t.Notes, &t.CreatedBy, &t.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}
