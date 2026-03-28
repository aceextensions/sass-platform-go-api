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
	GetInvitationByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*models.Invitation, error)
	UpdateInvitationStatus(ctx context.Context, id uuid.UUID, status string) error
	ListInvitations(ctx context.Context, tenantID uuid.UUID) ([]models.Invitation, error)
	DeleteInvitation(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error
	DeleteMembership(ctx context.Context, userID uuid.UUID, tenantID uuid.UUID) error

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
	var finalWhere string
	var finalArgs []interface{}

	if tenantID == uuid.Nil {
		if bq.WhereClause != "" {
			finalWhere = bq.WhereClause
		}
		finalArgs = bq.Args
	} else {
		finalWhere = "WHERE tenant_id = $1"
		if bq.WhereClause != "" {
			// bq.WhereClause starts with "WHERE ", so we replace it with " AND "
			finalWhere += " AND " + bq.WhereClause[6:]
		}
		finalArgs = append([]interface{}{tenantID}, bq.Args...)
	}

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

func (r *pgUserRepository) GetInvitationByID(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) (*models.Invitation, error) {
	query := `SELECT id, tenant_id, email, phone, role, token, expires_at, status, created_at FROM invitations WHERE id = $1 AND tenant_id = $2`
	var i models.Invitation
	err := r.getExecutor().QueryRow(ctx, query, id, tenantID).Scan(
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

func (r *pgUserRepository) ListInvitations(ctx context.Context, tenantID uuid.UUID) ([]models.Invitation, error) {
	query := `SELECT id, tenant_id, email, phone, role, token, expires_at, status, created_at FROM invitations WHERE tenant_id = $1 ORDER BY created_at DESC`
	rows, err := r.getExecutor().Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invites []models.Invitation
	for rows.Next() {
		var i models.Invitation
		err := rows.Scan(
			&i.ID, &i.TenantID, &i.Email, &i.Phone, &i.Role, &i.Token, &i.ExpiresAt, &i.Status, &i.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		invites = append(invites, i)
	}
	return invites, nil
}

func (r *pgUserRepository) DeleteInvitation(ctx context.Context, id uuid.UUID, tenantID uuid.UUID) error {
	query := `DELETE FROM invitations WHERE id = $1 AND tenant_id = $2`
	_, err := r.getExecutor().Exec(ctx, query, id, tenantID)
	return err
}

func (r *pgUserRepository) DeleteMembership(ctx context.Context, userID uuid.UUID, tenantID uuid.UUID) error {
	query := `DELETE FROM memberships WHERE user_id = $1 AND tenant_id = $2`
	_, err := r.getExecutor().Exec(ctx, query, userID, tenantID)
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
