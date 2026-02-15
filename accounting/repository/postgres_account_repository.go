package repository

import (
	"context"
	"time"

	"github.com/aceextension/accounting/domain"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type postgresAccountRepository struct {
	pool db.QueryExecutor
}

func NewPostgresAccountRepository(pool db.QueryExecutor) AccountRepository {
	return &postgresAccountRepository{pool: pool}
}

func (r *postgresAccountRepository) Create(ctx context.Context, account *domain.Account) error {
	query := `
		INSERT INTO accounts (id, tenant_id, code, name, type, parent_id, is_active, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.pool.Exec(ctx, query,
		account.ID, account.TenantID, account.Code, account.Name, account.Type,
		account.ParentID, account.IsActive, account.Description, account.CreatedAt, account.UpdatedAt,
	)
	return err
}

func (r *postgresAccountRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	query := `
		SELECT id, tenant_id, code, name, type, parent_id, is_active, description, created_at, updated_at
		FROM accounts
		WHERE id = $1
	`
	var acc domain.Account
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&acc.ID, &acc.TenantID, &acc.Code, &acc.Name, &acc.Type,
		&acc.ParentID, &acc.IsActive, &acc.Description, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

func (r *postgresAccountRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Account, error) {
	query := `
		SELECT id, tenant_id, code, name, type, parent_id, is_active, description, created_at, updated_at
		FROM accounts
		WHERE tenant_id = $1 AND code = $2
	`
	var acc domain.Account
	err := r.pool.QueryRow(ctx, query, tenantID, code).Scan(
		&acc.ID, &acc.TenantID, &acc.Code, &acc.Name, &acc.Type,
		&acc.ParentID, &acc.IsActive, &acc.Description, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

func (r *postgresAccountRepository) List(ctx context.Context, tenantID uuid.UUID) ([]*domain.Account, error) {
	query := `
		SELECT id, tenant_id, code, name, type, parent_id, is_active, description, created_at, updated_at
		FROM accounts
		WHERE tenant_id = $1
		ORDER BY code ASC
	`
	rows, err := r.pool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*domain.Account
	for rows.Next() {
		var acc domain.Account
		if err := rows.Scan(
			&acc.ID, &acc.TenantID, &acc.Code, &acc.Name, &acc.Type,
			&acc.ParentID, &acc.IsActive, &acc.Description, &acc.CreatedAt, &acc.UpdatedAt,
		); err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}
	return accounts, nil
}

func (r *postgresAccountRepository) Update(ctx context.Context, account *domain.Account) error {
	query := `
		UPDATE accounts
		SET code=$2, name=$3, type=$4, parent_id=$5, is_active=$6, description=$7, updated_at=$8
		WHERE id=$1
	`
	_, err := r.pool.Exec(ctx, query,
		account.ID, account.Code, account.Name, account.Type,
		account.ParentID, account.IsActive, account.Description, time.Now(),
	)
	return err
}

func (r *postgresAccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Soft delete by setting is_active = false
	// Alternatively, implement strict deletion logic if no transactions exist
	query := `UPDATE accounts SET is_active=false, updated_at=$2 WHERE id=$1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}
