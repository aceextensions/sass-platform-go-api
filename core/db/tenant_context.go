package db

import (
	"context"

	"github.com/google/uuid"
)

// TenantContextKey is the key for storing tenant ID in context
type contextKey string

const (
	TenantIDKey     contextKey = "tenant_id"
	UserIDKey       contextKey = "user_id"
	IsSuperAdminKey contextKey = "is_super_admin"
)

// WithTenantID adds tenant ID to context
func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, TenantIDKey, tenantID)
}

// GetTenantID retrieves tenant ID from context
func GetTenantID(ctx context.Context) (uuid.UUID, bool) {
	tenantID, ok := ctx.Value(TenantIDKey).(uuid.UUID)
	return tenantID, ok
}

// WithUserID adds user ID to context
func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID retrieves user ID from context
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// WithSuperAdmin marks context as super admin
func WithSuperAdmin(ctx context.Context, isSuperAdmin bool) context.Context {
	return context.WithValue(ctx, IsSuperAdminKey, isSuperAdmin)
}

// IsSuperAdmin checks if context has super admin privileges
func IsSuperAdmin(ctx context.Context) bool {
	isSuperAdmin, ok := ctx.Value(IsSuperAdminKey).(bool)
	return ok && isSuperAdmin
}
