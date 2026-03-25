package repository

import (
	"context"

	"github.com/aceextension/core/db"
	"github.com/aceextension/identity/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *models.Tenant) error
	GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error)
	GetTenantByName(ctx context.Context, name string) (*models.Tenant, error)
	UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) error

	// Membership support
	CreateMembership(ctx context.Context, membership *models.Membership) error
	GetMembershipsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Membership, error)

	// Transaction support
	WithTransaction(ctx context.Context, fn func(repo TenantRepository) error) error
	GetTx() pgx.Tx
}

type pgTenantRepository struct {
	tx pgx.Tx
}

func NewTenantRepository() TenantRepository {
	return &pgTenantRepository{}
}

func NewTenantRepositoryWithTx(tx pgx.Tx) TenantRepository {
	return &pgTenantRepository{tx: tx}
}

func (r *pgTenantRepository) GetTx() pgx.Tx {
	return r.tx
}

func (r *pgTenantRepository) getExecutor() db.QueryExecutor {
	if r.tx != nil {
		return r.tx
	}
	return db.MainPool
}

func (r *pgTenantRepository) CreateTenant(ctx context.Context, tenant *models.Tenant) error {
	query := `
		INSERT INTO tenants (name, business_name, status, category, fiscal_year_start, fiscal_year_end)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	if tenant.Category == "" {
		tenant.Category = models.TenantCategoryBusiness
	}

	return r.getExecutor().QueryRow(ctx, query,
		tenant.Name, tenant.BusinessName, tenant.Status, tenant.Category, tenant.FiscalYearStart, tenant.FiscalYearEnd,
	).Scan(&tenant.ID, &tenant.CreatedAt, &tenant.UpdatedAt)
}

func (r *pgTenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*models.Tenant, error) {
	query := `SELECT id, name, business_name, status, category, max_users, fiscal_year_start, fiscal_year_end, kyb_status, kyb_document_url, verified_at, is_active, created_at, updated_at FROM tenants WHERE id = $1`
	var tenant models.Tenant
	err := r.getExecutor().QueryRow(ctx, query, id).Scan(
		&tenant.ID, &tenant.Name, &tenant.BusinessName, &tenant.Status, &tenant.Category, &tenant.MaxUsers,
		&tenant.FiscalYearStart, &tenant.FiscalYearEnd, &tenant.KybStatus, &tenant.KybDocumentURL,
		&tenant.VerifiedAt, &tenant.IsActive, &tenant.CreatedAt, &tenant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *pgTenantRepository) GetTenantByName(ctx context.Context, name string) (*models.Tenant, error) {
	query := `SELECT id, name, business_name, status, category, max_users, fiscal_year_start, fiscal_year_end, kyb_status, kyb_document_url, verified_at, is_active, created_at, updated_at FROM tenants WHERE name = $1`
	var tenant models.Tenant
	err := r.getExecutor().QueryRow(ctx, query, name).Scan(
		&tenant.ID, &tenant.Name, &tenant.BusinessName, &tenant.Status, &tenant.Category, &tenant.MaxUsers,
		&tenant.FiscalYearStart, &tenant.FiscalYearEnd, &tenant.KybStatus, &tenant.KybDocumentURL,
		&tenant.VerifiedAt, &tenant.IsActive, &tenant.CreatedAt, &tenant.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *pgTenantRepository) UpdateTenantStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE tenants SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.getExecutor().Exec(ctx, query, status, id)
	return err
}

func (r *pgTenantRepository) CreateMembership(ctx context.Context, m *models.Membership) error {
	query := `
		INSERT INTO memberships (user_id, tenant_id, role, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	return r.getExecutor().QueryRow(ctx, query,
		m.UserID, m.TenantID, m.Role, m.Status,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *pgTenantRepository) GetMembershipsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Membership, error) {
	query := `SELECT id, user_id, tenant_id, role, status, created_at, updated_at FROM memberships WHERE user_id = $1`
	rows, err := r.getExecutor().Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []models.Membership
	for rows.Next() {
		var m models.Membership
		if err := rows.Scan(&m.ID, &m.UserID, &m.TenantID, &m.Role, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		memberships = append(memberships, m)
	}
	return memberships, nil
}

func (r *pgTenantRepository) WithTransaction(ctx context.Context, fn func(repo TenantRepository) error) error {
	if r.tx != nil {
		return fn(r)
	}

	return db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return fn(&pgTenantRepository{tx: tx})
	})
}
