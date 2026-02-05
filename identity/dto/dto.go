package dto

import "github.com/google/uuid"

type RegisterTenantDTO struct {
	TenantName string  `json:"tenantName" validate:"required"`
	OwnerName  string  `json:"ownerName" validate:"required"`
	Phone      string  `json:"phone" validate:"required"`
	Email      *string `json:"email"`
	Password   string  `json:"password" validate:"required,min=6"`
}

type LoginDTO struct {
	Phone    string `json:"phone" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type VerifyOTPDTO struct {
	Phone string `json:"phone" validate:"required"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type UserResponse struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Email    *string    `json:"email"`
	Phone    string     `json:"phone"`
	Role     string     `json:"role"`
	TenantID *uuid.UUID `json:"tenantId"`
}

type AuthResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

type TokenPayload struct {
	UserID   uuid.UUID  `json:"userId"`
	TenantID *uuid.UUID `json:"tenantId"`
	Role     string     `json:"role"`
}

type ForgotPasswordDTO struct {
	Identifier string `json:"identifier" validate:"required"`
}

type ResetPasswordDTO struct {
	Identifier  string `json:"identifier" validate:"required"`
	OTP         string `json:"otp" validate:"required,len=6"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}
