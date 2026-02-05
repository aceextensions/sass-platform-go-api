package models

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents the tenants table
type Tenant struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	Name               string     `json:"name" db:"name"`
	BusinessName       *string    `json:"businessName" db:"business_name"`
	TradeName          *string    `json:"tradeName" db:"trade_name"`
	PanNumber          *string    `json:"panNumber" db:"pan_number"`
	VatNumber          *string    `json:"vatNumber" db:"vat_number"`
	RegistrationNumber *string    `json:"registrationNumber" db:"registration_number"`
	Address            *string    `json:"address" db:"address"`
	Phone              *string    `json:"phone" db:"phone"`
	Email              *string    `json:"email" db:"email"`
	Status             string     `json:"status" db:"status"`
	MaxUsers           string     `json:"maxUsers" db:"max_users"`
	FiscalYearStart    *time.Time `json:"fiscalYearStart" db:"fiscal_year_start"`
	FiscalYearEnd      *time.Time `json:"fiscalYearEnd" db:"fiscal_year_end"`
	KybStatus          string     `json:"kybStatus" db:"kyb_status"`
	KybDocumentURL     *string    `json:"kybDocumentUrl" db:"kyb_document_url"`
	VerifiedAt         *time.Time `json:"verifiedAt" db:"verified_at"`
	IsActive           bool       `json:"isActive" db:"is_active"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time  `json:"updatedAt" db:"updated_at"`
}

// User represents the users table
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	TenantID     *uuid.UUID `json:"tenantId" db:"tenant_id"`
	Name         string     `json:"name" db:"name"`
	Email        *string    `json:"email" db:"email"`
	Phone        string     `json:"phone" db:"phone"`
	PasswordHash *string    `json:"-" db:"password_hash"`
	Role         string     `json:"role" db:"role"`
	IsVerified   bool       `json:"isVerified" db:"is_verified"`
	OTP          *string    `json:"-" db:"otp"`
	OTPExpiresAt *time.Time `json:"-" db:"otp_expires_at"`
	IsActive     bool       `json:"isActive" db:"is_active"`
	LastLogin    *time.Time `json:"lastLogin" db:"last_login"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

// Session represents the sessions table
type Session struct {
	ID                uuid.UUID `json:"id" db:"id"`
	UserID            uuid.UUID `json:"userId" db:"user_id"`
	RefreshToken      string    `json:"refreshToken" db:"refresh_token"`
	DeviceFingerprint *string   `json:"deviceFingerprint" db:"device_fingerprint"`
	IPAddress         *string   `json:"ipAddress" db:"ip_address"`
	ExpiresAt         time.Time `json:"expiresAt" db:"expires_at"`
	CreatedAt         time.Time `json:"createdAt" db:"created_at"`
}
