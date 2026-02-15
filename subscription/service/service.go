package service

import (
	"context"
	"errors"
	"time"

	"github.com/aceextension/subscription/domain"
	"github.com/aceextension/subscription/repository"
	"github.com/google/uuid"
)

// SubscriptionService defines the interface
type SubscriptionService interface {
	// Plan Management
	CreatePlan(ctx context.Context, plan *domain.Plan) error
	GetPlan(ctx context.Context, id uuid.UUID) (*domain.Plan, error)
	ListPlans(ctx context.Context) ([]*domain.Plan, error)
	UpdatePlan(ctx context.Context, plan *domain.Plan) error

	// Subscription Management
	Subscribe(ctx context.Context, tenantID, planID uuid.UUID) (*domain.Subscription, error)
	GetSubscription(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error)

	// Feature Gating
	HasFeature(ctx context.Context, tenantID uuid.UUID, feature string) (bool, error)
	CheckLimit(ctx context.Context, tenantID uuid.UUID, limitKey string, currentValue int) (bool, error)
}

type subscriptionService struct {
	planRepo repository.PlanRepository
	subRepo  repository.SubscriptionRepository
}

func NewSubscriptionService(planRepo repository.PlanRepository, subRepo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		planRepo: planRepo,
		subRepo:  subRepo,
	}
}

// Plan Management Implementation

func (s *subscriptionService) CreatePlan(ctx context.Context, plan *domain.Plan) error {
	// Check if code exists
	existing, _ := s.planRepo.GetByCode(ctx, plan.Code)
	if existing != nil {
		return errors.New("plan with this code already exists")
	}
	return s.planRepo.Create(ctx, plan)
}

func (s *subscriptionService) GetPlan(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	return s.planRepo.GetByID(ctx, id)
}

func (s *subscriptionService) ListPlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.List(ctx)
}

func (s *subscriptionService) UpdatePlan(ctx context.Context, plan *domain.Plan) error {
	return s.planRepo.Update(ctx, plan)
}

// Subscription Implementation

func (s *subscriptionService) Subscribe(ctx context.Context, tenantID, planID uuid.UUID) (*domain.Subscription, error) {
	plan, err := s.planRepo.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if !plan.IsActive {
		return nil, errors.New("plan is not active")
	}

	// Calculate dates
	startDate := time.Now()
	var endDate time.Time
	if plan.Interval == "YEARLY" {
		endDate = startDate.AddDate(1, 0, 0)
	} else {
		endDate = startDate.AddDate(0, 1, 0)
	}

	// Create subscription
	sub := domain.NewSubscription(tenantID, planID, startDate, endDate)

	// In a real system, we'd cancel existing active subscriptions first or queue this one
	// For now, simpler: just create new one which becomes the "latest"
	if err := s.subRepo.Create(ctx, sub); err != nil {
		return nil, err
	}

	sub.Plan = plan // Attach plan for return
	return sub, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error) {
	return s.subRepo.GetActiveByTenantID(ctx, tenantID)
}

// Feature Gating Implementation

func (s *subscriptionService) HasFeature(ctx context.Context, tenantID uuid.UUID, feature string) (bool, error) {
	sub, err := s.subRepo.GetActiveByTenantID(ctx, tenantID)
	if err != nil {
		return false, err
	}
	if sub == nil {
		// No active subscription -> No features (or maybe a default fallback?)
		// For stricter SaaS, no sub = no access. For freemium, might have a default free plan.
		// Assuming every tenant gets a Free plan on registration.
		return false, nil
	}

	// Check plan features
	allowed, ok := sub.Plan.Features[feature]
	if !ok {
		return false, nil // Feature not listed = not allowed
	}
	return allowed, nil
}

func (s *subscriptionService) CheckLimit(ctx context.Context, tenantID uuid.UUID, limitKey string, currentValue int) (bool, error) {
	sub, err := s.subRepo.GetActiveByTenantID(ctx, tenantID)
	if err != nil {
		return false, err
	}
	if sub == nil {
		return false, nil
	}

	limit, ok := sub.Plan.Limits[limitKey]
	if !ok {
		// Limit not defined = Unlimited? Or 0?
		// Let's assume -1 or missing means unlimited for now, OR 0 means defined as 0.
		// Safer: If limit key exists, check it. If not, maybe it's unrestricted or restricted.
		// Convention: "max_users" present with value 5 means 5. If missing, assume default low limit?
		// Let's assume if missing, it's 0 (not allowed). unlimited should be represented by -1.
		return false, nil
	}

	if limit == -1 {
		return true, nil // Unlimited
	}

	return currentValue < limit, nil
}
