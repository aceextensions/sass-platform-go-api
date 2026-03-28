package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/aceextension/core/config"
	"github.com/aceextension/core/db"
	"github.com/aceextension/identity/dto"
	"github.com/aceextension/identity/models"
	"github.com/aceextension/identity/repository"
	notifDomain "github.com/aceextension/notification/domain"
	notifService "github.com/aceextension/notification/service"
	"github.com/aceextension/common/templates"
	"github.com/google/uuid"
)

type UserService interface {
	ListUsers(ctx context.Context, tenantID uuid.UUID, options db.QueryOptions) (*dto.UserListResponse, error)
	InviteUser(ctx context.Context, actorID uuid.UUID, tenantID uuid.UUID, role string, data dto.InviteUserDTO) (*models.Invitation, error)
	JoinTenant(ctx context.Context, data dto.JoinTenantDTO) error
	ListInvitations(ctx context.Context, tenantID uuid.UUID) ([]dto.InvitationResponse, error)
	RevokeInvitation(ctx context.Context, actorRole string, tenantID uuid.UUID, id uuid.UUID) error
	ResendInvitation(ctx context.Context, actorRole string, tenantID uuid.UUID, id uuid.UUID) error
	RemoveMember(ctx context.Context, actorID uuid.UUID, actorRole string, tenantID uuid.UUID, userID uuid.UUID) error
}

type userService struct {
	userRepo   repository.UserRepository
	tenantRepo repository.TenantRepository
	authRepo   repository.AuthRepository
	notifServ  notifService.NotificationService
}

func NewUserService(userRepo repository.UserRepository, tenantRepo repository.TenantRepository, authRepo repository.AuthRepository, notifServ notifService.NotificationService) UserService {
	return &userService{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		authRepo:   authRepo,
		notifServ:  notifServ,
	}
}

func (s *userService) ListUsers(ctx context.Context, tenantID uuid.UUID, options db.QueryOptions) (*dto.UserListResponse, error) {
	// Start args from $2 because $1 is tenantID (if not Nil)
	startIndex := 1
	if tenantID != uuid.Nil {
		startIndex = 2
	}
	bq := db.BuildQuery(options, startIndex)

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

	// 4. Send Notification
	channel := notifDomain.ChannelEmail
	recipient := data.Email
	if recipient == "" && data.Phone != "" {
		channel = notifDomain.ChannelSMS
		recipient = data.Phone
	}

	inviteLink := config.GlobalConfig.AppBaseURL + "/join?token=" + invite.Token

	// Create Pongo2 context for the template
	templateVars := map[string]interface{}{
		"inviter_name": "Administrator",
		"tenant_name":  "Your Workspace",
		"role":         data.Role,
		"invite_link":  inviteLink,
		"subject":      "You've been invited to Practixa",
	}

	// Fetch tenant info for the beautiful template
	if tenant, err := s.tenantRepo.GetTenantByID(ctx, invite.TenantID); err == nil {
		templateVars["tenant_name"] = tenant.Name
	}

	// Render the content using the modular templates
	// We wrap the invitation template in the base layout
	fullTemplate := templates.Wrap(InvitationTemplate)

	s.notifServ.Send(ctx, notifService.SendRequest{
		TenantID:  tenantID,
		Channel:   channel,
		Recipient: recipient,
		Content:   fullTemplate,
		Variables: templateVars,
		Priority:  notifDomain.PriorityHigh,
	})

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

	// 2. Check if user already exists
	var existingUser *models.User
	if invite.Email != nil {
		existingUser, _ = s.authRepo.GetUserByEmail(ctx, *invite.Email)
	} else if invite.Phone != nil {
		existingUser, _ = s.authRepo.GetUserByPhone(ctx, *invite.Phone)
	}

	// 3. Process Join within Transaction
	return s.userRepo.WithTransaction(ctx, func(tr repository.UserRepository) error {
		var userID uuid.UUID

		if existingUser != nil {
			// User exists, just use their ID
			userID = existingUser.ID
			// Note: We ignore the provided password/name as they already have an account
		} else {
			// New user, create them
			hash, err := HashPassword(data.Password)
			if err != nil {
				return err
			}

			user := &models.User{
				TenantID:     &invite.TenantID,
				Name:         data.Name,
				PasswordHash: &hash,
				Role:         invite.Role,
				IsVerified:   true,
				IsActive:     true,
			}
			if invite.Email != nil {
				user.Email = invite.Email
			}
			if invite.Phone != nil {
				user.Phone = invite.Phone
			}

			authRepoTx := repository.NewAuthRepositoryWithTx(tr.GetTx())
			if err := authRepoTx.CreateUser(ctx, user); err != nil {
				return err
			}
			userID = user.ID
		}

		// 4. Create Membership
		tenantRepoTx := repository.NewTenantRepositoryWithTx(tr.GetTx())
		membership := &models.Membership{
			UserID:   userID,
			TenantID: invite.TenantID,
			Role:     invite.Role,
			Status:   "active",
		}

		if err := tenantRepoTx.CreateMembership(ctx, membership); err != nil {
			return err
		}

		// 5. Update Invitation
		return tr.UpdateInvitationStatus(ctx, invite.ID, "accepted")
	})
}

func (s *userService) ListInvitations(ctx context.Context, tenantID uuid.UUID) ([]dto.InvitationResponse, error) {
	invites, err := s.userRepo.ListInvitations(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.InvitationResponse, 0, len(invites))
	for _, i := range invites {
		res = append(res, dto.InvitationResponse{
			ID:        i.ID,
			Email:     i.Email,
			Phone:     i.Phone,
			Role:      i.Role,
			Status:    i.Status,
			ExpiresAt: i.ExpiresAt,
			CreatedAt: i.CreatedAt,
		})
	}
	return res, nil
}

func (s *userService) RevokeInvitation(ctx context.Context, actorRole string, tenantID uuid.UUID, id uuid.UUID) error {
	// 1. RBAC Check
	if actorRole != "owner" && actorRole != "manager" {
		return errors.New("unauthorized: only owners and managers can revoke invitations")
	}

	return s.userRepo.DeleteInvitation(ctx, id, tenantID)
}

func (s *userService) ResendInvitation(ctx context.Context, actorRole string, tenantID uuid.UUID, id uuid.UUID) error {
	// 1. RBAC Check
	if actorRole != "owner" && actorRole != "manager" {
		return errors.New("unauthorized: only owners and managers can resend invitations")
	}

	// 2. Fetch Translation
	invite, err := s.userRepo.GetInvitationByID(ctx, id, tenantID)
	if err != nil {
		return errors.New("invitation not found")
	}

	if invite.Status != "pending" {
		return errors.New("cannot resend: invitation is already " + invite.Status)
	}

	// 3. Re-send Notification
	channel := notifDomain.ChannelEmail
	recipient := ""
	if invite.Email != nil {
		recipient = *invite.Email
	} else if invite.Phone != nil {
		channel = notifDomain.ChannelSMS
		recipient = *invite.Phone
	}

	if recipient == "" {
		return errors.New("no recipient found for invitation")
	}

	inviteLink := config.GlobalConfig.AppBaseURL + "/join?token=" + invite.Token

	templateVars := map[string]interface{}{
		"inviter_name": "Administrator",
		"tenant_name":  "Your Workspace",
		"role":         invite.Role,
		"invite_link":  inviteLink,
		"subject":      "Reminder: You've been invited to Practixa",
	}

	if tenant, err := s.tenantRepo.GetTenantByID(ctx, invite.TenantID); err == nil {
		templateVars["tenant_name"] = tenant.Name
	}

	fullTemplate := templates.Wrap(InvitationTemplate)

	_, err = s.notifServ.Send(ctx, notifService.SendRequest{
		TenantID:  tenantID,
		Channel:   channel,
		Recipient: recipient,
		Content:   fullTemplate,
		Variables: templateVars,
		Priority:  notifDomain.PriorityHigh,
	})
	return err
}

func (s *userService) RemoveMember(ctx context.Context, actorID uuid.UUID, actorRole string, tenantID uuid.UUID, userID uuid.UUID) error {
	// 1. RBAC Check
	if actorRole != "owner" && actorRole != "manager" {
		return errors.New("unauthorized: only owners and managers can remove members")
	}

	// 2. Prevent self-removal (they should delete tenant or leave instead)
	if actorID == userID {
		return errors.New("cannot remove yourself from the workspace")
	}

	// 3. Prevent non-owners from removing owners
	targetUser, err := s.authRepo.GetUserByID(ctx, userID)
	if err == nil && targetUser.Role == "owner" && actorRole != "owner" {
		return errors.New("only owners can remove other owners")
	}

	// 4. Delete Membership
	return s.userRepo.DeleteMembership(ctx, userID, tenantID)
}

func generateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
