package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/aceextension/identity/dto"
	"github.com/aceextension/identity/models"
	"github.com/aceextension/identity/repository"
	"github.com/google/uuid"
)

type AuthService interface {
	RegisterTenant(ctx context.Context, data dto.RegisterTenantDTO) (*dto.UserResponse, error)
	VerifyOTP(ctx context.Context, data dto.VerifyOTPDTO) (*dto.AuthResponse, error)
	Login(ctx context.Context, data dto.LoginDTO) (*dto.AuthResponse, error)
	Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	ForgotPassword(ctx context.Context, data dto.ForgotPasswordDTO) error
	ResetPassword(ctx context.Context, data dto.ResetPasswordDTO) error
	Impersonate(ctx context.Context, tenantID uuid.UUID, adminUserID uuid.UUID) (*dto.AuthResponse, error)
	GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error)
}

type authService struct {
	authRepo   repository.AuthRepository
	tenantRepo repository.TenantRepository
}

func NewAuthService(authRepo repository.AuthRepository, tenantRepo repository.TenantRepository) AuthService {
	return &authService{
		authRepo:   authRepo,
		tenantRepo: tenantRepo,
	}
}

func (s *authService) RegisterTenant(ctx context.Context, data dto.RegisterTenantDTO) (*dto.UserResponse, error) {
	// 1. Hash Password
	passwordHash, err := HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	// 2. Generate OTP
	otp := "123456" // Default for dev as per Bun implementation
	otpExpiresAt := time.Now().Add(10 * time.Minute)

	var user models.User

	// 3. Execution in Transaction
	err = s.tenantRepo.WithTransaction(ctx, func(tr repository.TenantRepository) error {
		// Create Tenant
		tenant := models.Tenant{
			Name:         data.TenantName,
			BusinessName: &data.TenantName,
			Status:       "trial",
		}

		// Set fiscal year defaults
		now := time.Now()
		fs := time.Date(now.Year(), 4, 1, 0, 0, 0, 0, time.Local)
		fe := time.Date(now.Year()+1, 3, 31, 23, 59, 59, 0, time.Local)
		tenant.FiscalYearStart = &fs
		tenant.FiscalYearEnd = &fe

		if err := tr.CreateTenant(ctx, &tenant); err != nil {
			return err
		}

		// Create Owner User using the shared transaction if possible
		// Since AuthRepo and TenantRepo use the same pool, we can share the Tx
		// For now, I'll pass the tenant ID to the authRepo
		user = models.User{
			TenantID:     &tenant.ID,
			Name:         data.OwnerName,
			Phone:        data.Phone,
			Email:        data.Email,
			PasswordHash: &passwordHash,
			Role:         "owner",
			IsVerified:   false,
			OTP:          &otp,
			OTPExpiresAt: &otpExpiresAt,
		}

		// We need to ensure AuthRepo uses the SAME transaction
		authRepoTx := repository.NewAuthRepositoryWithTx(tr.GetTx())
		return authRepoTx.CreateUser(ctx, &user)
	})

	if err != nil {
		// Handle Postgres unique violations (simplified for now)
		return nil, err
	}

	fmt.Printf("📱 OTP for %s: %s (expires in 10 minutes)\n", data.Phone, otp)

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		TenantID: user.TenantID,
	}, nil
}

func (s *authService) VerifyOTP(ctx context.Context, data dto.VerifyOTPDTO) (*dto.AuthResponse, error) {
	user, err := s.authRepo.GetUserByPhone(ctx, data.Phone)
	if err != nil {
		return nil, errors.New("invalid OTP or user already verified")
	}

	if user.IsVerified {
		return nil, errors.New("user already verified")
	}

	if user.OTP == nil || *user.OTP != data.OTP {
		return nil, errors.New("invalid OTP")
	}

	if user.OTPExpiresAt != nil && time.Now().After(*user.OTPExpiresAt) {
		return nil, errors.New("OTP expired")
	}

	// Update verification status
	if err := s.authRepo.UpdateUserVerification(ctx, user.ID, true); err != nil {
		return nil, err
	}

	// Generate Tokens
	payload := dto.TokenPayload{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
	}

	accessToken, err := GenerateAccessToken(payload)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateRefreshToken(payload)
	if err != nil {
		return nil, err
	}

	// Create Session
	session := models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.authRepo.CreateSession(ctx, &session); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
	}, nil
}

func (s *authService) Login(ctx context.Context, data dto.LoginDTO) (*dto.AuthResponse, error) {
	// Support both phone and email login
	var user *models.User
	var err error

	user, err = s.authRepo.GetUserByPhone(ctx, data.Phone)
	if err != nil {
		user, err = s.authRepo.GetUserByEmail(ctx, data.Phone)
	}

	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsVerified {
		return nil, errors.New("account not verified")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify Password
	if user.PasswordHash == nil || !ComparePassword(data.Password, *user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Update last login
	_ = s.authRepo.UpdateLastLogin(ctx, user.ID)

	// Generate Tokens
	payload := dto.TokenPayload{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
	}

	accessToken, _ := GenerateAccessToken(payload)
	refreshToken, _ := GenerateRefreshToken(payload)

	// Create Session
	session := models.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	_ = s.authRepo.CreateSession(ctx, &session)

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID, refreshToken string) error {
	return s.authRepo.DeleteSession(ctx, userID, refreshToken)
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	session, err := s.authRepo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.authRepo.DeleteSession(ctx, session.UserID, refreshToken)
		return nil, errors.New("session expired")
	}

	user, err := s.authRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	payload := dto.TokenPayload{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Role:     user.Role,
	}

	newAccessToken, _ := GenerateAccessToken(payload)
	newRefreshToken, _ := GenerateRefreshToken(payload)

	// Replace old session with new one
	_ = s.authRepo.DeleteSession(ctx, session.UserID, refreshToken)
	newSession := models.Session{
		UserID:       user.ID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(7 * 24 * time.Hour),
	}
	_ = s.authRepo.CreateSession(ctx, &newSession)

	return &dto.AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
	}, nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.PasswordHash == nil || !ComparePassword(oldPassword, *user.PasswordHash) {
		return errors.New("invalid old password")
	}

	newHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.authRepo.UpdateUserPassword(ctx, userID, newHash)
}

func (s *authService) ForgotPassword(ctx context.Context, data dto.ForgotPasswordDTO) error {
	user, err := s.authRepo.GetUserByIdentifier(ctx, data.Identifier)
	if err != nil {
		return errors.New("user not found")
	}

	otp := "123456" // Default for dev
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := s.authRepo.UpdateOTP(ctx, user.ID, &otp, &expiresAt); err != nil {
		return err
	}

	fmt.Printf("📧 Password Reset OTP for %s: %s\n", data.Identifier, otp)
	return nil
}

func (s *authService) ResetPassword(ctx context.Context, data dto.ResetPasswordDTO) error {
	user, err := s.authRepo.GetUserByIdentifier(ctx, data.Identifier)
	if err != nil {
		return errors.New("user not found")
	}

	if user.OTP == nil || *user.OTP != data.OTP {
		return errors.New("invalid code")
	}

	if user.OTPExpiresAt != nil && time.Now().After(*user.OTPExpiresAt) {
		return errors.New("code expired")
	}

	newHash, err := HashPassword(data.NewPassword)
	if err != nil {
		return err
	}

	return s.authRepo.WithTransaction(ctx, func(repo repository.AuthRepository) error {
		if err := repo.UpdateUserPassword(ctx, user.ID, newHash); err != nil {
			return err
		}
		// Clear OTP
		return repo.UpdateOTP(ctx, user.ID, nil, nil)
	})
}

func (s *authService) Impersonate(ctx context.Context, tenantID uuid.UUID, adminUserID uuid.UUID) (*dto.AuthResponse, error) {
	// 1. Get owner user of the tenant
	// Note: In a real app we'd find the primary owner. For now, get any owner.
	users, err := s.authRepo.GetUsersByTenantID(ctx, tenantID)
	if err != nil || len(users) == 0 {
		return nil, errors.New("tenant user not found")
	}

	var targetUser *models.User
	for _, u := range users {
		if u.Role == "owner" {
			targetUser = u
			break
		}
	}
	if targetUser == nil {
		targetUser = users[0]
	}

	payload := dto.TokenPayload{
		UserID:   targetUser.ID,
		TenantID: targetUser.TenantID,
		Role:     targetUser.Role,
	}

	accessToken, _ := GenerateAccessToken(payload)
	refreshToken, _ := GenerateRefreshToken(payload)

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       targetUser.ID,
			Name:     targetUser.Name,
			Email:    targetUser.Email,
			Phone:    targetUser.Phone,
			Role:     targetUser.Role,
			TenantID: targetUser.TenantID,
		},
	}, nil
}

func (s *authService) GetMe(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.authRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		TenantID: user.TenantID,
	}, nil
}

func generateRandomOTP() string {
	// ... existing implementation ...
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
