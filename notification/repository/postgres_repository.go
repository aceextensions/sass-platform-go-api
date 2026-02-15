package repository

import (
	"context"
	"fmt"

	"github.com/aceextension/core/db"
	"github.com/aceextension/notification/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// PostgresTemplateRepository implements TemplateRepository
type PostgresTemplateRepository struct{}

// NewPostgresTemplateRepository creates a new PostgreSQL template repository
func NewPostgresTemplateRepository() *PostgresTemplateRepository {
	return &PostgresTemplateRepository{}
}

// Create creating new template
func (r *PostgresTemplateRepository) Create(ctx context.Context, template *domain.Template) error {
	query := `
		INSERT INTO templates (
			id, tenant_id, code, channel, subject, body, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := db.MainPool.Exec(ctx, query,
		template.ID, template.TenantID, template.Code, template.Channel,
		template.Subject, template.Body, template.IsActive, template.CreatedAt, template.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}
	return nil
}

// GetByID retrieving template by ID
func (r *PostgresTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Template, error) {
	query := `
		SELECT id, tenant_id, code, channel, subject, body, is_active, created_at, updated_at
		FROM templates WHERE id = $1
	`
	return r.scanTemplate(db.MainPool.QueryRow(ctx, query, id))
}

// GetByCode retrieving template by code and channel for a tenant
func (r *PostgresTemplateRepository) GetByCode(ctx context.Context, tenantID uuid.UUID, code string, channel domain.ChannelType) (*domain.Template, error) {
	query := `
		SELECT id, tenant_id, code, channel, subject, body, is_active, created_at, updated_at
		FROM templates WHERE tenant_id = $1 AND code = $2 AND channel = $3
	`
	return r.scanTemplate(db.MainPool.QueryRow(ctx, query, tenantID, code, channel))
}

// GetByTenantID retrieving all templates for a tenant
func (r *PostgresTemplateRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) ([]*domain.Template, error) {
	query := `
		SELECT id, tenant_id, code, channel, subject, body, is_active, created_at, updated_at
		FROM templates WHERE tenant_id = $1 ORDER BY code
	`
	rows, err := db.MainPool.Query(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to list templates: %w", err)
	}
	defer rows.Close()

	var templates []*domain.Template
	for rows.Next() {
		template, err := r.scanTemplateRow(rows)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
	}
	return templates, nil
}

// Update updating template
func (r *PostgresTemplateRepository) Update(ctx context.Context, template *domain.Template) error {
	query := `
		UPDATE templates SET
			subject = $1, body = $2, is_active = $3, updated_at = $4
		WHERE id = $5
	`
	_, err := db.MainPool.Exec(ctx, query,
		template.Subject, template.Body, template.IsActive, template.UpdatedAt, template.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}
	return nil
}

// Delete deleting template
func (r *PostgresTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := db.MainPool.Exec(ctx, "DELETE FROM templates WHERE id = $1", id)
	return err
}

func (r *PostgresTemplateRepository) scanTemplate(row pgx.Row) (*domain.Template, error) {
	var t domain.Template
	err := row.Scan(
		&t.ID, &t.TenantID, &t.Code, &t.Channel, &t.Subject, &t.Body,
		&t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *PostgresTemplateRepository) scanTemplateRow(rows pgx.Rows) (*domain.Template, error) {
	var t domain.Template
	err := rows.Scan(
		&t.ID, &t.TenantID, &t.Code, &t.Channel, &t.Subject, &t.Body,
		&t.IsActive, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// PostgresNotificationRepository implements NotificationRepository
type PostgresNotificationRepository struct{}

// NewPostgresNotificationRepository creates a new PostgreSQL notification repository
func NewPostgresNotificationRepository() *PostgresNotificationRepository {
	return &PostgresNotificationRepository{}
}

// Create creating new notification
func (r *PostgresNotificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (
			id, tenant_id, user_id, channel, recipient, subject, content,
			priority, status, retry_count, error_message, sent_at, template_id, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := db.MainPool.Exec(ctx, query,
		n.ID, n.TenantID, n.UserID, n.Channel, n.Recipient, n.Subject, n.Content,
		n.Priority, n.Status, n.RetryCount, n.ErrorMessage, n.SentAt, n.TemplateID, n.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

// GetByID retrieving notification by ID
func (r *PostgresNotificationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error) {
	query := `
		SELECT id, tenant_id, user_id, channel, recipient, subject, content,
		       priority, status, retry_count, error_message, sent_at, template_id, created_at
		FROM notifications WHERE id = $1
	`
	return r.scanNotification(db.MainPool.QueryRow(ctx, query, id))
}

// GetByTenantID retrieving notifications for a tenant
func (r *PostgresNotificationRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.Notification, error) {
	query := `
		SELECT id, tenant_id, user_id, channel, recipient, subject, content,
		       priority, status, retry_count, error_message, sent_at, template_id, created_at
		FROM notifications WHERE tenant_id = $1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3
	`
	rows, err := db.MainPool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n, err := r.scanNotificationRow(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

// Update updating notification status
func (r *PostgresNotificationRepository) Update(ctx context.Context, n *domain.Notification) error {
	query := `
		UPDATE notifications SET
			status = $1, retry_count = $2, error_message = $3, sent_at = $4
		WHERE id = $5
	`
	_, err := db.MainPool.Exec(ctx, query,
		n.Status, n.RetryCount, n.ErrorMessage, n.SentAt, n.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}
	return nil
}

// GetPending retrieving pending notifications for worker
func (r *PostgresNotificationRepository) GetPending(ctx context.Context, limit int) ([]*domain.Notification, error) {
	query := `
		SELECT id, tenant_id, user_id, channel, recipient, subject, content,
		       priority, status, retry_count, error_message, sent_at, template_id, created_at
		FROM notifications
		WHERE status IN ('PENDING', 'FAILED') AND retry_count < 3
		ORDER BY priority DESC, created_at ASC
		LIMIT $1
	`
	rows, err := db.MainPool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*domain.Notification
	for rows.Next() {
		n, err := r.scanNotificationRow(rows)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}
	return notifications, nil
}

func (r *PostgresNotificationRepository) scanNotification(row pgx.Row) (*domain.Notification, error) {
	var n domain.Notification
	err := row.Scan(
		&n.ID, &n.TenantID, &n.UserID, &n.Channel, &n.Recipient, &n.Subject, &n.Content,
		&n.Priority, &n.Status, &n.RetryCount, &n.ErrorMessage, &n.SentAt, &n.TemplateID, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *PostgresNotificationRepository) scanNotificationRow(rows pgx.Rows) (*domain.Notification, error) {
	var n domain.Notification
	err := rows.Scan(
		&n.ID, &n.TenantID, &n.UserID, &n.Channel, &n.Recipient, &n.Subject, &n.Content,
		&n.Priority, &n.Status, &n.RetryCount, &n.ErrorMessage, &n.SentAt, &n.TemplateID, &n.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &n, nil
}
