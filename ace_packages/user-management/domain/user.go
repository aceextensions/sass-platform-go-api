package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

type User struct {
	ID                uuid.UUID              `json:"id"`
	TenantID          uuid.UUID              `json:"tenantId"`
	Email             string                 `json:"email"`
	Name              string                 `json:"name"`
	Status            UserStatus             `json:"status"`
	Roles             []Role                 `json:"roles"`
	Permissions       []Permission           `json:"permissions"`
	APITokens         []APIToken             `json:"apiTokens"`
	Sessions          []Session              `json:"sessions"`
	TOTPSecret        *string                `json:"totpSecret,omitempty"`
	BackupCodes       []string               `json:"backupCodes,omitempty"`
	IPWhitelist       []string               `json:"ipWhitelist,omitempty"`
	LastLoginAt       *time.Time             `json:"lastLoginAt,omitempty"`
	LastLoginIP       string                 `json:"lastLoginIp"`
	LoginFailureCount int                    `json:"loginFailureCount"`
	Metadata          map[string]interface{} `json:"metadata"`
	Attributes        map[string]interface{} `json:"attributes"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
	DeletedAt         *time.Time             `json:"deletedAt,omitempty"`
}

type Role struct {
	ID          uuid.UUID    `json:"id"`
	TenantID    uuid.UUID    `json:"tenantId"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

type Permission struct {
	ID       uuid.UUID `json:"id"`
	Resource string    `json:"resource"` // "users", "tenants", "billing"
	Action   string    `json:"action"`   // "create", "read", "update", "delete"
}

type APIToken struct {
	ID          uuid.UUID    `json:"id"`
	UserID      uuid.UUID    `json:"userId"`
	Token       string       `json:"token"` // Hashed
	Name        string       `json:"name"`
	Permissions []Permission `json:"permissions"`
	LastUsedAt  *time.Time   `json:"lastUsedAt,omitempty"`
	ExpiresAt   *time.Time   `json:"expiresAt,omitempty"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

type Session struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	TokenHash  string    `json:"tokenHash"`
	IPAddress  string    `json:"ipAddress"`
	UserAgent  string    `json:"userAgent"`
	CreatedAt  time.Time `json:"createdAt"`
	ExpiresAt  time.Time `json:"expiresAt"`
}
