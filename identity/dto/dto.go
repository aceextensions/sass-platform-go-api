package dto

import (
	"time"

	"github.com/google/uuid"
)

type RegisterTenantDTO struct {
	TenantName string  `json:"tenantName" validate:"required,min=3,max=50"`
	OwnerName  string  `json:"ownerName" validate:"required,min=2,max=100"`
	Phone      string  `json:"phone" validate:"required,min=10,max=15"`
	Email      *string `json:"email" validate:"omitempty,email"`
	Password   string  `json:"password" validate:"required,min=6,max=50"`
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
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     *string    `json:"email"`
	Phone     string     `json:"phone"`
	Role      string     `json:"role"`
	TenantID  *uuid.UUID `json:"tenantId"`
	IsActive  bool       `json:"isActive"`
	CreatedAt time.Time  `json:"createdAt"`
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
	NewPassword string `json:"newPassword" validate:"required,min=6,max=50"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutDTO struct {
	RefreshToken *string `json:"refreshToken"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,min=6"`
}

type InviteUserDTO struct {
	Email string `json:"email" validate:"required_without=Phone,omitempty,email"`
	Phone string `json:"phone" validate:"required_without=Email,omitempty,min=10,max=15"`
	Role  string `json:"role" validate:"required,oneof=owner manager staff admin"`
}

type JoinTenantDTO struct {
	Token    string `json:"token" validate:"required"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserListResponse struct {
	Data       []UserResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
	TotalItems  int `json:"totalItems"`
	TotalPages  int `json:"totalPages"`
	CurrentPage int `json:"currentPage"`
	Limit       int `json:"limit"`
}
