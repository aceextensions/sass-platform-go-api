package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/identity/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserRepository interface {
	// User Management
	ListUsers(ctx context.Context, tenantID uuid.UUID, query db.BuiltQuery) ([]models.User, int, error)
	GetUserCountByTenant(ctx context.Context, tenantID uuid.UUID) (int, error)

	// Invitation Management
	CreateInvitation(ctx context.Context, invite *models.Invitation) error
	GetInvitationByToken(ctx context.Context, token string) (*models.Invitation, error)
	UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status string) error

	// Transaction support
	WithTransaction(ctx context.Context, fn func(repo UserRepository) error) error
	GetTx() pgx.Tx
}

type pgUserRepository struct {
	tx pgx.Tx
}

func NewUserRepository() UserRepository {
	return &pgUserRepository{}
}

func NewUserRepositoryWithTx(tx pgx.Tx) UserRepository {
	return &pgUserRepository{tx: tx}
}

func (r *pgUserRepository) GetTx() pgx.Tx {
	return r.tx
}

func (r *pgUserRepository) getExecutor() db.QueryExecutor {
	if r.tx != nil {
		return r.tx
	}
	return db.MainPool
}

func (r *pgUserRepository) ListUsers(ctx context.Context, tenantID uuid.UUID, bq db.BuiltQuery) ([]models.User, int, error) {
	// 1. Get total count with filters (but without limit/offset)

	// Inject tenantID into args for security
	// Our BuildQuery generated $1, $2 etc starting from an index.
	// We need to ensure tenant_id is always present.

	// Wait, BuildQuery might not have tenant_id in WhereClause.
	// Let's refine the approach: we always append tenant_id.

	finalWhere := "WHERE tenant_id = $1"
	if bq.WhereClause != "" {
		finalWhere += " AND " + bq.WhereClause[6:] // Remove "WHERE "
	}

	finalArgs := append([]interface{}{tenantID}, bq.Args...)

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM users %s", finalWhere)
	var total int
	err := r.getExecutor().QueryRow(ctx, countSQL, finalArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// 2. Get data
	dataSQL := fmt.Sprintf("SELECT id, tenant_id, name, email, phone, role, is_verified, is_active, last_login, created_at, updated_at FROM users %s %s LIMIT %d OFFSET %d",
		finalWhere, bq.OrderBy, bq.Limit, bq.Offset)

	rows, err := r.getExecutor().Query(ctx, dataSQL, finalArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.TenantID, &u.Name, &u.Email, &u.Phone, &u.Role,
			&u.IsVerified, &u.IsActive, &u.LastLogin, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *pgUserRepository) GetUserCountByTenant(ctx context.Context, tenantID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM users WHERE tenant_id = $1`
	var count int
	err := r.getExecutor().QueryRow(ctx, query, tenantID).Scan(&count)
	return count, err
}

func (r *pgUserRepository) CreateInvitation(ctx context.Context, invite *models.Invitation) error {
	query := `
		INSERT INTO invitations (tenant_id, email, phone, role, token, expires_at, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	return r.getExecutor().QueryRow(ctx, query,
		invite.TenantID, invite.Email, invite.Phone, invite.Role, invite.Token, invite.ExpiresAt, invite.Status,
	).Scan(&invite.ID, &invite.CreatedAt)
}

func (r *pgUserRepository) GetInvitationByToken(ctx context.Context, token string) (*models.Invitation, error) {
	query := `SELECT id, tenant_id, email, phone, role, token, expires_at, status, created_at FROM invitations WHERE token = $1`
	var i models.Invitation
	err := r.getExecutor().QueryRow(ctx, query, token).Scan(
		&i.ID, &i.TenantID, &i.Email, &i.Phone, &i.Role, &i.Token, &i.ExpiresAt, &i.Status, &i.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *pgUserRepository) UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE invitations SET status = $1 WHERE id = $2`
	_, err := r.getExecutor().Exec(ctx, query, status, id)
	return err
}

func (r *pgUserRepository) WithTransaction(ctx context.Context, fn func(repo UserRepository) error) error {
	if r.tx != nil {
		return fn(r)
	}

	return db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return fn(&pgUserRepository{tx: tx})
	})
}
