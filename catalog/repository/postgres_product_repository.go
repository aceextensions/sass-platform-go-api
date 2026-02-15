package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aceextension/catalog/domain"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgresProductRepository implements ProductRepository using PostgreSQL
type PostgresProductRepository struct{}

// NewPostgresProductRepository creates a new PostgreSQL product repository
func NewPostgresProductRepository() *PostgresProductRepository {
	return &PostgresProductRepository{}
}

// Create creates a new product
func (r *PostgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	attrsJSON, err := json.Marshal(product.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		INSERT INTO products (
			id, tenant_id, product_code, name, description, category_id,
			cost_price, selling_price, mrp, tax_rate,
			sku, barcode, unit, status, is_active,
			custom_attributes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err = db.MainPool.Exec(ctx, query,
		product.ID, product.TenantID, product.ProductCode,
		product.Name, product.Description, product.CategoryID,
		product.CostPrice, product.SellingPrice, product.MRP, product.TaxRate,
		product.SKU, product.Barcode, product.Unit, product.Status, product.IsActive,
		attrsJSON, product.CreatedAt, product.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

// GetByID retrieves a product by ID
func (r *PostgresProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, product_code, name, description, category_id,
		       cost_price, selling_price, mrp, tax_rate,
		       sku, barcode, unit, status, is_active,
		       custom_attributes, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	return r.scanProduct(db.MainPool.QueryRow(ctx, query, id))
}

// GetByCode retrieves a product by code
func (r *PostgresProductRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, product_code, name, description, category_id,
		       cost_price, selling_price, mrp, tax_rate,
		       sku, barcode, unit, status, is_active,
		       custom_attributes, created_at, updated_at
		FROM products
		WHERE tenant_id = $1 AND product_code = $2
	`

	return r.scanProduct(db.MainPool.QueryRow(ctx, query, tenantID, code))
}

// GetBySKU retrieves a product by SKU
func (r *PostgresProductRepository) GetBySKU(ctx context.Context, tenantID uuid.UUID, sku string) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, product_code, name, description, category_id,
		       cost_price, selling_price, mrp, tax_rate,
		       sku, barcode, unit, status, is_active,
		       custom_attributes, created_at, updated_at
		FROM products
		WHERE tenant_id = $1 AND sku = $2
	`

	return r.scanProduct(db.MainPool.QueryRow(ctx, query, tenantID, sku))
}

// GetByBarcode retrieves a product by barcode
func (r *PostgresProductRepository) GetByBarcode(ctx context.Context, tenantID uuid.UUID, barcode string) (*domain.Product, error) {
	query := `
		SELECT id, tenant_id, product_code, name, description, category_id,
		       cost_price, selling_price, mrp, tax_rate,
		       sku, barcode, unit, status, is_active,
		       custom_attributes, created_at, updated_at
		FROM products
		WHERE tenant_id = $1 AND barcode = $2
	`

	return r.scanProduct(db.MainPool.QueryRow(ctx, query, tenantID, barcode))
}

// GetByCategory retrieves products by category
func (r *PostgresProductRepository) GetByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, tenant_id, product_code, name, description, category_id,
		       cost_price, selling_price, mrp, tax_rate,
		       sku, barcode, unit, status, is_active,
		       custom_attributes, created_at, updated_at
		FROM products
		WHERE category_id = $1
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := db.MainPool.Query(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query products by category: %w", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

// GetByTenantID retrieves all products for a tenant
func (r *PostgresProductRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, tenant_id, product_code, name, description, category_id,
		       cost_price, selling_price, mrp, tax_rate,
		       sku, barcode, unit, status, is_active,
		       custom_attributes, created_at, updated_at
		FROM products
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

// Update updates a product
func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	attrsJSON, err := json.Marshal(product.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		UPDATE products
		SET name = $1, description = $2, category_id = $3,
		    cost_price = $4, selling_price = $5, mrp = $6, tax_rate = $7,
		    sku = $8, barcode = $9, unit = $10, status = $11, is_active = $12,
		    custom_attributes = $13, updated_at = $14
		WHERE id = $15
	`

	_, err = db.MainPool.Exec(ctx, query,
		product.Name, product.Description, product.CategoryID,
		product.CostPrice, product.SellingPrice, product.MRP, product.TaxRate,
		product.SKU, product.Barcode, product.Unit, product.Status, product.IsActive,
		attrsJSON, product.UpdatedAt, product.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

// Delete deletes a product
func (r *PostgresProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	_, err := db.MainPool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

// Search searches products by name, SKU, or barcode
func (r *PostgresProductRepository) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Product, error) {
	searchQuery := `
		SELECT id, tenant_id, product_code, name, description, category_id,
		       cost_price, selling_price, mrp, tax_rate,
		       sku, barcode, unit, status, is_active,
		       custom_attributes, created_at, updated_at
		FROM products
		WHERE tenant_id = $1
		AND (
			name ILIKE $2
			OR product_code ILIKE $2
			OR sku ILIKE $2
			OR barcode ILIKE $2
			OR description ILIKE $2
		)
		ORDER BY name
		LIMIT $3 OFFSET $4
	`

	searchPattern := "%" + query + "%"
	rows, err := db.MainPool.Query(ctx, searchQuery, tenantID, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	return r.scanProducts(rows)
}

// Count returns total number of products for a tenant
func (r *PostgresProductRepository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM products WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return count, nil
}

// GetNextProductNumber gets the next product number for code generation
func (r *PostgresProductRepository) GetNextProductNumber(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM products WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get next product number: %w", err)
	}

	return count + 1, nil
}

// scanProduct scans a single product row
func (r *PostgresProductRepository) scanProduct(row pgx.Row) (*domain.Product, error) {
	var product domain.Product
	var attrsJSON []byte

	err := row.Scan(
		&product.ID, &product.TenantID, &product.ProductCode,
		&product.Name, &product.Description, &product.CategoryID,
		&product.CostPrice, &product.SellingPrice, &product.MRP, &product.TaxRate,
		&product.SKU, &product.Barcode, &product.Unit, &product.Status, &product.IsActive,
		&attrsJSON, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan product: %w", err)
	}

	// Unmarshal custom attributes
	if len(attrsJSON) > 0 {
		if err := json.Unmarshal(attrsJSON, &product.CustomAttributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
		}
	}

	if product.CustomAttributes == nil {
		product.CustomAttributes = make(map[string]interface{})
	}

	return &product, nil
}

// scanProducts scans multiple product rows
func (r *PostgresProductRepository) scanProducts(rows pgx.Rows) ([]*domain.Product, error) {
	products := []*domain.Product{}

	for rows.Next() {
		var product domain.Product
		var attrsJSON []byte

		err := rows.Scan(
			&product.ID, &product.TenantID, &product.ProductCode,
			&product.Name, &product.Description, &product.CategoryID,
			&product.CostPrice, &product.SellingPrice, &product.MRP, &product.TaxRate,
			&product.SKU, &product.Barcode, &product.Unit, &product.Status, &product.IsActive,
			&attrsJSON, &product.CreatedAt, &product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		// Unmarshal custom attributes
		if len(attrsJSON) > 0 {
			if err := json.Unmarshal(attrsJSON, &product.CustomAttributes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
			}
		}

		if product.CustomAttributes == nil {
			product.CustomAttributes = make(map[string]interface{})
		}

		products = append(products, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return products, nil
}
