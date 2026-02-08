package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/crm/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgresSupplierRepository implements SupplierRepository using PostgreSQL
type PostgresSupplierRepository struct{}

// NewPostgresSupplierRepository creates a new PostgreSQL supplier repository
func NewPostgresSupplierRepository() *PostgresSupplierRepository {
	return &PostgresSupplierRepository{}
}

// Create creates a new supplier
func (r *PostgresSupplierRepository) Create(ctx context.Context, supplier *domain.Supplier) error {
	attrsJSON, err := json.Marshal(supplier.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		INSERT INTO suppliers (
			id, tenant_id, supplier_code, name, email, phone,
			supplier_type, status, custom_attributes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err = db.MainPool.Exec(ctx, query,
		supplier.ID, supplier.TenantID, supplier.SupplierCode,
		supplier.Name, supplier.Email, supplier.Phone,
		supplier.SupplierType, supplier.Status, attrsJSON,
		supplier.CreatedAt, supplier.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create supplier: %w", err)
	}

	return nil
}

// GetByID retrieves a supplier by ID
func (r *PostgresSupplierRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Supplier, error) {
	query := `
		SELECT id, tenant_id, supplier_code, name, email, phone,
		       supplier_type, status, custom_attributes, created_at, updated_at
		FROM suppliers
		WHERE id = $1
	`

	return r.scanSupplier(db.MainPool.QueryRow(ctx, query, id))
}

// GetByCode retrieves a supplier by supplier code
func (r *PostgresSupplierRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Supplier, error) {
	query := `
		SELECT id, tenant_id, supplier_code, name, email, phone,
		       supplier_type, status, custom_attributes, created_at, updated_at
		FROM suppliers
		WHERE tenant_id = $1 AND supplier_code = $2
	`

	return r.scanSupplier(db.MainPool.QueryRow(ctx, query, tenantID, code))
}

// GetByTenantID retrieves all suppliers for a tenant
func (r *PostgresSupplierRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Supplier, error) {
	query := `
		SELECT id, tenant_id, supplier_code, name, email, phone,
		       supplier_type, status, custom_attributes, created_at, updated_at
		FROM suppliers
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query suppliers: %w", err)
	}
	defer rows.Close()

	return r.scanSuppliers(rows)
}

// Update updates a supplier
func (r *PostgresSupplierRepository) Update(ctx context.Context, supplier *domain.Supplier) error {
	attrsJSON, err := json.Marshal(supplier.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		UPDATE suppliers
		SET name = $1, email = $2, phone = $3, supplier_type = $4,
		    status = $5, custom_attributes = $6, updated_at = $7
		WHERE id = $8
	`

	_, err = db.MainPool.Exec(ctx, query,
		supplier.Name, supplier.Email, supplier.Phone, supplier.SupplierType,
		supplier.Status, attrsJSON, supplier.UpdatedAt, supplier.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update supplier: %w", err)
	}

	return nil
}

// Delete deletes a supplier
func (r *PostgresSupplierRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM suppliers WHERE id = $1`

	_, err := db.MainPool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete supplier: %w", err)
	}

	return nil
}

// Search searches suppliers by name, email, or phone
func (r *PostgresSupplierRepository) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Supplier, error) {
	searchQuery := `
		SELECT id, tenant_id, supplier_code, name, email, phone,
		       supplier_type, status, custom_attributes, created_at, updated_at
		FROM suppliers
		WHERE tenant_id = $1
		AND (
			name ILIKE $2
			OR email ILIKE $2
			OR phone ILIKE $2
			OR supplier_code ILIKE $2
		)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	searchPattern := "%" + query + "%"
	rows, err := db.MainPool.Query(ctx, searchQuery, tenantID, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search suppliers: %w", err)
	}
	defer rows.Close()

	return r.scanSuppliers(rows)
}

// SearchByCustomAttribute searches suppliers by custom attribute
func (r *PostgresSupplierRepository) SearchByCustomAttribute(ctx context.Context, tenantID uuid.UUID, key, value string) ([]*domain.Supplier, error) {
	query := `
		SELECT id, tenant_id, supplier_code, name, email, phone,
		       supplier_type, status, custom_attributes, created_at, updated_at
		FROM suppliers
		WHERE tenant_id = $1
		AND custom_attributes->>$2 = $3
		ORDER BY created_at DESC
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID, key, value)
	if err != nil {
		return nil, fmt.Errorf("failed to search suppliers by custom attribute: %w", err)
	}
	defer rows.Close()

	return r.scanSuppliers(rows)
}

// Count returns total number of suppliers for a tenant
func (r *PostgresSupplierRepository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM suppliers WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count suppliers: %w", err)
	}

	return count, nil
}

// GetNextSupplierNumber gets the next supplier number for code generation
func (r *PostgresSupplierRepository) GetNextSupplierNumber(ctx context.Context, tenantID uuid.UUID) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM suppliers WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get next supplier number: %w", err)
	}

	return count + 1, nil
}

// scanSupplier scans a single supplier row
func (r *PostgresSupplierRepository) scanSupplier(row pgx.Row) (*domain.Supplier, error) {
	var supplier domain.Supplier
	var attrsJSON []byte

	err := row.Scan(
		&supplier.ID, &supplier.TenantID, &supplier.SupplierCode,
		&supplier.Name, &supplier.Email, &supplier.Phone,
		&supplier.SupplierType, &supplier.Status, &attrsJSON,
		&supplier.CreatedAt, &supplier.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan supplier: %w", err)
	}

	// Unmarshal custom attributes
	if len(attrsJSON) > 0 {
		if err := json.Unmarshal(attrsJSON, &supplier.CustomAttributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
		}
	}

	if supplier.CustomAttributes == nil {
		supplier.CustomAttributes = make(map[string]interface{})
	}

	return &supplier, nil
}

// scanSuppliers scans multiple supplier rows
func (r *PostgresSupplierRepository) scanSuppliers(rows pgx.Rows) ([]*domain.Supplier, error) {
	suppliers := []*domain.Supplier{}

	for rows.Next() {
		var supplier domain.Supplier
		var attrsJSON []byte

		err := rows.Scan(
			&supplier.ID, &supplier.TenantID, &supplier.SupplierCode,
			&supplier.Name, &supplier.Email, &supplier.Phone,
			&supplier.SupplierType, &supplier.Status, &attrsJSON,
			&supplier.CreatedAt, &supplier.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan supplier: %w", err)
		}

		// Unmarshal custom attributes
		if len(attrsJSON) > 0 {
			if err := json.Unmarshal(attrsJSON, &supplier.CustomAttributes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
			}
		}

		if supplier.CustomAttributes == nil {
			supplier.CustomAttributes = make(map[string]interface{})
		}

		suppliers = append(suppliers, &supplier)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return suppliers, nil
}
