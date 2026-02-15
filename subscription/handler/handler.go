package handler

import (
	"net/http"

	"github.com/aceextension/core/db"
	authService "github.com/aceextension/identity/service"
	"github.com/aceextension/subscription/domain"
	"github.com/aceextension/subscription/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// PlanHandler handles plan administration
type PlanHandler struct {
	service service.SubscriptionService
}

func NewPlanHandler(service service.SubscriptionService) *PlanHandler {
	return &PlanHandler{service: service}
}

// CreatePlanRequest matches the request body
type CreatePlanRequest struct {
	Name        string          `json:"name" validate:"required"`
	Code        string          `json:"code" validate:"required"`
	Description string          `json:"description"`
	Price       float64         `json:"price" validate:"gte=0"`
	Interval    string          `json:"interval" validate:"oneof=MONTHLY YEARLY"`
	Features    map[string]bool `json:"features"`
	Limits      map[string]int  `json:"limits"`
}

// Create creates a new plan
// @Summary Create a plan
// @Description Create a new subscription plan (Admin)
// @Tags plans
// @Accept json
// @Produce json
// @Param request body CreatePlanRequest true "Plan Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/plans [post]
// @Security BearerAuth
func (h *PlanHandler) Create(c echo.Context) error {
	var req CreatePlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	plan := domain.NewPlan(req.Name, req.Code, req.Description, req.Price, req.Interval)
	if req.Features != nil {
		plan.Features = req.Features
	}
	if req.Limits != nil {
		plan.Limits = req.Limits
	}

	if err := h.service.CreatePlan(c.Request().Context(), plan); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": plan.ID.String()})
}

// List lists all plans
// @Summary List plans
// @Description List all available subscription plans
// @Tags plans
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /api/v1/plans [get]
func (h *PlanHandler) List(c echo.Context) error {
	plans, err := h.service.ListPlans(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	// Swagger local struct fix if needed, but domain.Plan is clean
	return c.JSON(http.StatusOK, plans)
}

// SubscriptionHandler handles subscriptions
type SubscriptionHandler struct {
	service     service.SubscriptionService
	authService authService.AuthService
}

func NewSubscriptionHandler(service service.SubscriptionService, authService authService.AuthService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, authService: authService}
}

// GetCurrentSubscription returns the current tenant's subscription
// @Summary Get current subscription
// @Description Get usage and subscription details for the tenant
// @Tags subscriptions
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]string
// @Router /api/v1/subscriptions/current [get]
// @Security BearerAuth
func (h *SubscriptionHandler) GetCurrentSubscription(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	sub, err := h.service.GetSubscription(c.Request().Context(), tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	if sub == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "No active subscription"})
	}

	return c.JSON(http.StatusOK, sub)
}

// SubscribeRequest for changing plans
type SubscribeRequest struct {
	PlanID   string `json:"planId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Subscribe subscribes tenant to a plan
// @Summary Subscribe to plan
// @Description Upgrade or change subscription plan
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body SubscribeRequest true "Subscribe Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/subscriptions/subscribe [post]
// @Security BearerAuth
func (h *SubscriptionHandler) Subscribe(c echo.Context) error {
	var req SubscribeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	planID, err := uuid.Parse(req.PlanID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Plan ID"})
	}

	sub, err := h.service.Subscribe(c.Request().Context(), tenantID, planID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, sub)
}
