package domain

import (
	"time"

	"github.com/google/uuid"
)

// Plan represents a subscription plan
type Plan struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Code        string          `json:"code"` // e.g., "silver-monthly"
	Description string          `json:"description"`
	Price       float64         `json:"price"`
	Currency    string          `json:"currency"` // "NPR", "USD"
	Interval    string          `json:"interval"` // "MONTHLY", "YEARLY"
	Features    map[string]bool `json:"features"`
	Limits      map[string]int  `json:"limits"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// SubscriptionStatus defines the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionStatusPastDue   SubscriptionStatus = "PAST_DUE"
)

// Subscription represents a tenant's subscription
type Subscription struct {
	ID        uuid.UUID          `json:"id"`
	TenantID  uuid.UUID          `json:"tenantId"`
	PlanID    uuid.UUID          `json:"planId"`
	Plan      *Plan              `json:"plan,omitempty" gorm:"-"` // Loaded relation
	Status    SubscriptionStatus `json:"status"`
	StartDate time.Time          `json:"startDate"`
	EndDate   time.Time          `json:"endDate"`
	AutoRenew bool               `json:"autoRenew"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

// NewPlan creates a new plan
func NewPlan(name, code, description string, price float64, interval string) *Plan {
	return &Plan{
		ID:          uuid.New(),
		Name:        name,
		Code:        code,
		Description: description,
		Price:       price,
		Currency:    "NPR",
		Interval:    interval,
		Features:    make(map[string]bool),
		Limits:      make(map[string]int),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewSubscription creates a new subscription
func NewSubscription(tenantID, planID uuid.UUID, startDate, endDate time.Time) *Subscription {
	return &Subscription{
		ID:        uuid.New(),
		TenantID:  tenantID,
		PlanID:    planID,
		Status:    SubscriptionStatusActive,
		StartDate: startDate,
		EndDate:   endDate,
		AutoRenew: true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
