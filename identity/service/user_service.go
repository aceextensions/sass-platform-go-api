package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/aceextension/core/db"
	"github.com/aceextension/identity/dto"
	"github.com/aceextension/identity/models"
	"github.com/aceextension/identity/repository"
	"github.com/google/uuid"
)

type UserService interface {
	ListUsers(ctx context.Context, tenantID uuid.UUID, options db.QueryOptions) (*dto.UserListResponse, error)
	InviteUser(ctx context.Context, actorID uuid.UUID, tenantID uuid.UUID, role string, data dto.InviteUserDTO) (*models.Invitation, error)
	JoinTenant(ctx context.Context, data dto.JoinTenantDTO) error
}

type userService struct {
	userRepo   repository.UserRepository
	tenantRepo repository.TenantRepository
	authRepo   repository.AuthRepository
}

func NewUserService(userRepo repository.UserRepository, tenantRepo repository.TenantRepository, authRepo repository.AuthRepository) UserService {
	return &userService{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		authRepo:   authRepo,
	}
}

func (s *userService) ListUsers(ctx context.Context, tenantID uuid.UUID, options db.QueryOptions) (*dto.UserListResponse, error) {
	// Start args from $2 because $1 is tenantID
	bq := db.BuildQuery(options, 2)

	users, total, err := s.userRepo.ListUsers(ctx, tenantID, bq)
	if err != nil {
		return nil, err
	}

	resData := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resData = append(resData, dto.UserResponse{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Phone:     u.Phone,
			Role:      u.Role,
			TenantID:  u.TenantID,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt,
		})
	}

	limit := bq.Limit
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.UserListResponse{
		Data: resData,
		Pagination: dto.Pagination{
			TotalItems:  total,
			TotalPages:  totalPages,
			CurrentPage: options.Page,
			Limit:       limit,
		},
	}, nil
}

func (s *userService) InviteUser(ctx context.Context, actorID uuid.UUID, tenantID uuid.UUID, actorRole string, data dto.InviteUserDTO) (*models.Invitation, error) {
	// 1. RBAC Check
	if actorRole != "owner" && actorRole != "manager" {
		return nil, errors.New("unauthorized: only owners and managers can invite users")
	}

	// 2. Seat Limit Check
	tenant, err := s.tenantRepo.GetTenantByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	maxUsers, _ := strconv.Atoi(tenant.MaxUsers)
	if maxUsers == 0 {
		maxUsers = 5 // Default
	}

	currentCount, err := s.userRepo.GetUserCountByTenant(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if currentCount >= maxUsers {
		return nil, errors.New("seat limit reached for this tenant")
	}

	// 3. Create Invitation
	token, _ := generateRandomToken(32)
	invite := &models.Invitation{
		TenantID:  tenantID,
		Role:      data.Role,
		Token:     token,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		Status:    "pending",
	}

	if data.Email != "" {
		invite.Email = &data.Email
	}
	if data.Phone != "" {
		invite.Phone = &data.Phone
	}

	if err := s.userRepo.CreateInvitation(ctx, invite); err != nil {
		return nil, err
	}

	// In real app: Send notification (email/sms)
	return invite, nil
}

func (s *userService) JoinTenant(ctx context.Context, data dto.JoinTenantDTO) error {
	// 1. Validate Token
	invite, err := s.userRepo.GetInvitationByToken(ctx, data.Token)
	if err != nil {
		return errors.New("invalid or expired invitation")
	}

	if invite.Status != "pending" || time.Now().After(invite.ExpiresAt) {
		return errors.New("invitation is no longer valid")
	}

	// 2. Prepare User
	hash, err := HashPassword(data.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		TenantID:     &invite.TenantID,
		Name:         data.Name,
		Email:        invite.Email,
		Phone:        "", // Optional or from invite
		PasswordHash: &hash,
		Role:         invite.Role,
		IsVerified:   true,
		IsActive:     true,
	}
	if invite.Phone != nil {
		user.Phone = *invite.Phone
	}

	// 3. Execute in Transaction
	return s.userRepo.WithTransaction(ctx, func(tr repository.UserRepository) error {
		// Create User (using auth repository logic but shared transaction)
		authRepoTx := repository.NewAuthRepositoryWithTx(tr.GetTx())
		if err := authRepoTx.CreateUser(ctx, user); err != nil {
			return err
		}

		// Update Invitation
		return tr.UpdateInvitationStatus(ctx, invite.ID, "accepted")
	})
}

func generateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
