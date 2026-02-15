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

// PostgresCategoryRepository implements CategoryRepository using PostgreSQL
type PostgresCategoryRepository struct{}

// NewPostgresCategoryRepository creates a new PostgreSQL category repository
func NewPostgresCategoryRepository() *PostgresCategoryRepository {
	return &PostgresCategoryRepository{}
}

// Create creates a new category
func (r *PostgresCategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	attrsJSON, err := json.Marshal(category.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		INSERT INTO categories (
			id, tenant_id, category_code, name, description, parent_id,
			level, path, sort_order, is_active, custom_attributes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = db.MainPool.Exec(ctx, query,
		category.ID, category.TenantID, category.CategoryCode,
		category.Name, category.Description, category.ParentID,
		category.Level, category.Path, category.SortOrder, category.IsActive,
		attrsJSON, category.CreatedAt, category.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	return nil
}

// GetByID retrieves a category by ID
func (r *PostgresCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	query := `
		SELECT id, tenant_id, category_code, name, description, parent_id,
		       level, path, sort_order, is_active, custom_attributes, created_at, updated_at
		FROM categories
		WHERE id = $1
	`

	return r.scanCategory(db.MainPool.QueryRow(ctx, query, id))
}

// GetByCode retrieves a category by code
func (r *PostgresCategoryRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Category, error) {
	query := `
		SELECT id, tenant_id, category_code, name, description, parent_id,
		       level, path, sort_order, is_active, custom_attributes, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1 AND category_code = $2
	`

	return r.scanCategory(db.MainPool.QueryRow(ctx, query, tenantID, code))
}

// GetByTenantID retrieves all categories for a tenant
func (r *PostgresCategoryRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Category, error) {
	query := `
		SELECT id, tenant_id, category_code, name, description, parent_id,
		       level, path, sort_order, is_active, custom_attributes, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1
		ORDER BY sort_order, name
		LIMIT $2 OFFSET $3
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetRootCategories retrieves root categories (no parent)
func (r *PostgresCategoryRepository) GetRootCategories(ctx context.Context, tenantID uuid.UUID) ([]*domain.Category, error) {
	query := `
		SELECT id, tenant_id, category_code, name, description, parent_id,
		       level, path, sort_order, is_active, custom_attributes, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1 AND parent_id IS NULL
		ORDER BY sort_order, name
	`

	rows, err := db.MainPool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to query root categories: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// GetChildren retrieves child categories
func (r *PostgresCategoryRepository) GetChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.Category, error) {
	query := `
		SELECT id, tenant_id, category_code, name, description, parent_id,
		       level, path, sort_order, is_active, custom_attributes, created_at, updated_at
		FROM categories
		WHERE parent_id = $1
		ORDER BY sort_order, name
	`

	rows, err := db.MainPool.Query(ctx, query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query child categories: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// Update updates a category
func (r *PostgresCategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	attrsJSON, err := json.Marshal(category.CustomAttributes)
	if err != nil {
		return fmt.Errorf("failed to marshal custom attributes: %w", err)
	}

	query := `
		UPDATE categories
		SET name = $1, description = $2, parent_id = $3, level = $4,
		    path = $5, sort_order = $6, is_active = $7, custom_attributes = $8, updated_at = $9
		WHERE id = $10
	`

	_, err = db.MainPool.Exec(ctx, query,
		category.Name, category.Description, category.ParentID, category.Level,
		category.Path, category.SortOrder, category.IsActive, attrsJSON,
		category.UpdatedAt, category.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	return nil
}

// Delete deletes a category
func (r *PostgresCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM categories WHERE id = $1`

	_, err := db.MainPool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// Search searches categories by name
func (r *PostgresCategoryRepository) Search(ctx context.Context, tenantID uuid.UUID, query string, limit, offset int) ([]*domain.Category, error) {
	searchQuery := `
		SELECT id, tenant_id, category_code, name, description, parent_id,
		       level, path, sort_order, is_active, custom_attributes, created_at, updated_at
		FROM categories
		WHERE tenant_id = $1
		AND (name ILIKE $2 OR category_code ILIKE $2 OR description ILIKE $2)
		ORDER BY sort_order, name
		LIMIT $3 OFFSET $4
	`

	searchPattern := "%" + query + "%"
	rows, err := db.MainPool.Query(ctx, searchQuery, tenantID, searchPattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search categories: %w", err)
	}
	defer rows.Close()

	return r.scanCategories(rows)
}

// Count returns total number of categories for a tenant
func (r *PostgresCategoryRepository) Count(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM categories WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count categories: %w", err)
	}

	return count, nil
}

// GetNextCategoryNumber gets the next category number for code generation
func (r *PostgresCategoryRepository) GetNextCategoryNumber(ctx context.Context, tenantID uuid.UUID) (int64, error) {
	var count int64

	query := `SELECT COUNT(*) FROM categories WHERE tenant_id = $1`

	err := db.MainPool.QueryRow(ctx, query, tenantID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get next category number: %w", err)
	}

	return count + 1, nil
}

// scanCategory scans a single category row
func (r *PostgresCategoryRepository) scanCategory(row pgx.Row) (*domain.Category, error) {
	var category domain.Category
	var attrsJSON []byte

	err := row.Scan(
		&category.ID, &category.TenantID, &category.CategoryCode,
		&category.Name, &category.Description, &category.ParentID,
		&category.Level, &category.Path, &category.SortOrder, &category.IsActive,
		&attrsJSON, &category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to scan category: %w", err)
	}

	// Unmarshal custom attributes
	if len(attrsJSON) > 0 {
		if err := json.Unmarshal(attrsJSON, &category.CustomAttributes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
		}
	}

	if category.CustomAttributes == nil {
		category.CustomAttributes = make(map[string]interface{})
	}

	return &category, nil
}

// scanCategories scans multiple category rows
func (r *PostgresCategoryRepository) scanCategories(rows pgx.Rows) ([]*domain.Category, error) {
	categories := []*domain.Category{}

	for rows.Next() {
		var category domain.Category
		var attrsJSON []byte

		err := rows.Scan(
			&category.ID, &category.TenantID, &category.CategoryCode,
			&category.Name, &category.Description, &category.ParentID,
			&category.Level, &category.Path, &category.SortOrder, &category.IsActive,
			&attrsJSON, &category.CreatedAt, &category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		// Unmarshal custom attributes
		if len(attrsJSON) > 0 {
			if err := json.Unmarshal(attrsJSON, &category.CustomAttributes); err != nil {
				return nil, fmt.Errorf("failed to unmarshal custom attributes: %w", err)
			}
		}

		if category.CustomAttributes == nil {
			category.CustomAttributes = make(map[string]interface{})
		}

		categories = append(categories, &category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return categories, nil
}
