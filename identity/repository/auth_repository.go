package repository

import (
	"context"
	"time"

	"github.com/aceextension/core/db"
	"github.com/aceextension/identity/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByPhone(ctx context.Context, phone string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUserVerification(ctx context.Context, userID uuid.UUID, isVerified bool) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	UpdateOTP(ctx context.Context, userID uuid.UUID, otp *string, expiresAt *time.Time) error

	// Session management
	CreateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, userID uuid.UUID, refreshToken string) error
	GetSessionByToken(ctx context.Context, refreshToken string) (*models.Session, error)
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error

	// Transaction support for registration
	WithTransaction(ctx context.Context, fn func(repo AuthRepository) error) error
}

type pgAuthRepository struct {
	tx pgx.Tx // For transactions
}

func NewAuthRepository() AuthRepository {
	return &pgAuthRepository{}
}

func (r *pgAuthRepository) getExecutor() db.QueryExecutor {
	if r.tx != nil {
		return r.tx
	}
	return db.MainPool
}

func (r *pgAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (tenant_id, name, email, phone, password_hash, role, is_verified, otp, otp_expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.getExecutor().QueryRow(ctx, query,
		user.TenantID, user.Name, user.Email, user.Phone, user.PasswordHash,
		user.Role, user.IsVerified, user.OTP, user.OTPExpiresAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *pgAuthRepository) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `SELECT id, tenant_id, name, email, phone, password_hash, role, is_verified, otp, otp_expires_at, is_active, last_login, created_at, updated_at FROM users WHERE phone = $1`
	var user models.User
	err := r.getExecutor().QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.TenantID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash,
		&user.Role, &user.IsVerified, &user.OTP, &user.OTPExpiresAt, &user.IsActive,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *pgAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, tenant_id, name, email, phone, password_hash, role, is_verified, otp, otp_expires_at, is_active, last_login, created_at, updated_at FROM users WHERE email = $1`
	var user models.User
	err := r.getExecutor().QueryRow(ctx, query, email).Scan(
		&user.ID, &user.TenantID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash,
		&user.Role, &user.IsVerified, &user.OTP, &user.OTPExpiresAt, &user.IsActive,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *pgAuthRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, tenant_id, name, email, phone, password_hash, role, is_verified, otp, otp_expires_at, is_active, last_login, created_at, updatedAt FROM users WHERE id = $1`
	var user models.User
	err := r.getExecutor().QueryRow(ctx, query, id).Scan(
		&user.ID, &user.TenantID, &user.Name, &user.Email, &user.Phone, &user.PasswordHash,
		&user.Role, &user.IsVerified, &user.OTP, &user.OTPExpiresAt, &user.IsActive,
		&user.LastLogin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *pgAuthRepository) UpdateUserVerification(ctx context.Context, userID uuid.UUID, isVerified bool) error {
	query := `UPDATE users SET is_verified = $1, otp = NULL, otp_expires_at = NULL, updated_at = NOW() WHERE id = $2`
	_, err := r.getExecutor().Exec(ctx, query, isVerified, userID)
	return err
}

func (r *pgAuthRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.getExecutor().Exec(ctx, query, passwordHash, userID)
	return err
}

func (r *pgAuthRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = $1`
	_, err := r.getExecutor().Exec(ctx, query, userID)
	return err
}

func (r *pgAuthRepository) UpdateOTP(ctx context.Context, userID uuid.UUID, otp *string, expiresAt *time.Time) error {
	query := `UPDATE users SET otp = $1, otp_expires_at = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.getExecutor().Exec(ctx, query, otp, expiresAt, userID)
	return err
}

func (r *pgAuthRepository) CreateSession(ctx context.Context, session *models.Session) error {
	query := `
		INSERT INTO sessions (user_id, refresh_token, device_fingerprint, ip_address, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	return r.getExecutor().QueryRow(ctx, query,
		session.UserID, session.RefreshToken, session.DeviceFingerprint,
		session.IPAddress, session.ExpiresAt,
	).Scan(&session.ID, &session.CreatedAt)
}

func (r *pgAuthRepository) DeleteSession(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	query := `DELETE FROM sessions WHERE user_id = $1 AND refresh_token = $2`
	_, err := r.getExecutor().Exec(ctx, query, userID, refreshToken)
	return err
}

func (r *pgAuthRepository) GetSessionByToken(ctx context.Context, refreshToken string) (*models.Session, error) {
	query := `SELECT id, user_id, refresh_token, device_fingerprint, ip_address, expires_at, created_at FROM sessions WHERE refresh_token = $1`
	var session models.Session
	err := r.getExecutor().QueryRow(ctx, query, refreshToken).Scan(
		&session.ID, &session.UserID, &session.RefreshToken, &session.DeviceFingerprint,
		&session.IPAddress, &session.ExpiresAt, &session.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *pgAuthRepository) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM sessions WHERE user_id = $1`
	_, err := r.getExecutor().Exec(ctx, query, userID)
	return err
}

func (r *pgAuthRepository) WithTransaction(ctx context.Context, fn func(repo AuthRepository) error) error {
	if r.tx != nil {
		return fn(r) // Already in a transaction
	}

	return db.BeginFunc(ctx, func(tx pgx.Tx) error {
		return fn(&pgAuthRepository{tx: tx})
	})
}
