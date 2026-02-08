package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aceextension/audit/domain"
	"github.com/aceextension/core/db"
	"github.com/google/uuid"
)

// PostgresAuditRepository implements AuditRepository using PostgreSQL
type PostgresAuditRepository struct{}

// NewPostgresAuditRepository creates a new PostgreSQL audit repository
func NewPostgresAuditRepository() *PostgresAuditRepository {
	return &PostgresAuditRepository{}
}

// Create inserts a new audit log entry
func (r *PostgresAuditRepository) Create(ctx context.Context, log *domain.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			id, tenant_id, user_id, action, entity, entity_id,
			ip_address, user_agent, details, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Convert details to JSON
	var detailsJSON []byte
	var err error
	if log.Details != nil {
		detailsJSON, err = json.Marshal(log.Details)
		if err != nil {
			return fmt.Errorf("failed to marshal details: %w", err)
		}
	}

	_, err = db.AuditPool.Exec(ctx, query,
		log.ID,
		log.TenantID,
		log.UserID,
		log.Action,
		log.Entity,
		log.EntityID,
		log.IPAddress,
		log.UserAgent,
		detailsJSON,
		log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetByID retrieves an audit log by ID
func (r *PostgresAuditRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity, entity_id,
		       ip_address, user_agent, details, created_at
		FROM audit_logs
		WHERE id = $1
	`

	var log domain.AuditLog
	var detailsJSON []byte

	err := db.AuditPool.QueryRow(ctx, query, id).Scan(
		&log.ID,
		&log.TenantID,
		&log.UserID,
		&log.Action,
		&log.Entity,
		&log.EntityID,
		&log.IPAddress,
		&log.UserAgent,
		&detailsJSON,
		&log.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}

	// Unmarshal details
	if len(detailsJSON) > 0 {
		if err := json.Unmarshal(detailsJSON, &log.Details); err != nil {
			return nil, fmt.Errorf("failed to unmarshal details: %w", err)
		}
	}

	return &log, nil
}

// GetByTenantID retrieves audit logs for a specific tenant
func (r *PostgresAuditRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity, entity_id,
		       ip_address, user_agent, details, created_at
		FROM audit_logs
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.AuditPool.Query(ctx, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// GetByEntity retrieves audit logs for a specific entity
func (r *PostgresAuditRepository) GetByEntity(ctx context.Context, entity string, entityID string, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity, entity_id,
		       ip_address, user_agent, details, created_at
		FROM audit_logs
		WHERE entity = $1 AND entity_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := db.AuditPool.Query(ctx, query, entity, entityID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// GetByUserID retrieves audit logs for a specific user
func (r *PostgresAuditRepository) GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity, entity_id,
		       ip_address, user_agent, details, created_at
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.AuditPool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// Search retrieves audit logs with filters
func (r *PostgresAuditRepository) Search(ctx context.Context, filters *AuditSearchFilters) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, tenant_id, user_id, action, entity, entity_id,
		       ip_address, user_agent, details, created_at
		FROM audit_logs
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 1

	if filters.TenantID != nil {
		query += fmt.Sprintf(" AND tenant_id = $%d", argCount)
		args = append(args, *filters.TenantID)
		argCount++
	}

	if filters.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argCount)
		args = append(args, *filters.UserID)
		argCount++
	}

	if filters.Action != nil {
		query += fmt.Sprintf(" AND action = $%d", argCount)
		args = append(args, *filters.Action)
		argCount++
	}

	if filters.Entity != nil {
		query += fmt.Sprintf(" AND entity = $%d", argCount)
		args = append(args, *filters.Entity)
		argCount++
	}

	if filters.EntityID != nil {
		query += fmt.Sprintf(" AND entity_id = $%d", argCount)
		args = append(args, *filters.EntityID)
		argCount++
	}

	if filters.StartDate != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argCount)
		args = append(args, *filters.StartDate)
		argCount++
	}

	if filters.EndDate != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argCount)
		args = append(args, *filters.EndDate)
		argCount++
	}

	query += " ORDER BY created_at DESC"

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filters.Limit)
		argCount++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filters.Offset)
	}

	rows, err := db.AuditPool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search audit logs: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

// scanRows is a helper function to scan multiple rows
func (r *PostgresAuditRepository) scanRows(rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}) ([]*domain.AuditLog, error) {
	logs := []*domain.AuditLog{}

	for rows.Next() {
		var log domain.AuditLog
		var detailsJSON []byte

		err := rows.Scan(
			&log.ID,
			&log.TenantID,
			&log.UserID,
			&log.Action,
			&log.Entity,
			&log.EntityID,
			&log.IPAddress,
			&log.UserAgent,
			&detailsJSON,
			&log.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		// Unmarshal details
		if len(detailsJSON) > 0 {
			if err := json.Unmarshal(detailsJSON, &log.Details); err != nil {
				return nil, fmt.Errorf("failed to unmarshal details: %w", err)
			}
		}

		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return logs, nil
}
