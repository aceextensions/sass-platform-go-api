package handler

import (
	"net/http"
	"time"

	"github.com/aceextension/core/db"
	"github.com/aceextension/core/middleware"
	"github.com/aceextension/notification"
	"github.com/aceextension/notification/domain"
	"github.com/aceextension/notification/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// TemplateHandler handles template requests
type TemplateHandler struct {
	service service.NotificationService
}

// NewTemplateHandler creates a new template handler
func NewTemplateHandler(service service.NotificationService) *TemplateHandler {
	return &TemplateHandler{service: service}
}

// CreateTemplateRequest request body
type CreateTemplateRequest struct {
	Code    string `json:"code" validate:"required"`
	Channel string `json:"channel" validate:"required"`
	Subject string `json:"subject"`
	Body    string `json:"body" validate:"required"`
}

// TemplateResponse represents the template response structure for Swagger
type TemplateResponse struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenantId"`
	Code      string    `json:"code"`
	Channel   string    `json:"channel"`
	Subject   *string   `json:"subject,omitempty"`
	Body      string    `json:"body"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Create creates a new template
// @Summary Create a template
// @Description Create a new notification template
// @Tags templates
// @Accept json
// @Produce json
// @Param request body CreateTemplateRequest true "Template Request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/notifications/templates [post]
// @Security BearerAuth
func (h *TemplateHandler) Create(c echo.Context) error {
	var req CreateTemplateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	template := domain.NewTemplate(tenantID, req.Code, domain.ChannelType(req.Channel), req.Body)
	if req.Subject != "" {
		template.Subject = &req.Subject
	}

	if err := h.service.CreateTemplate(c.Request().Context(), template); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"id": template.ID.String()})
}

// List lists templates
// @Summary List templates
// @Description List all templates for tenant
// @Tags templates
// @Produce json
// @Success 200 {array} TemplateResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/notifications/templates [get]
// @Security BearerAuth
func (h *TemplateHandler) List(c echo.Context) error {
	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	templates, err := h.service.GetTemplates(c.Request().Context(), tenantID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var response []TemplateResponse
	for _, t := range templates {
		response = append(response, TemplateResponse{
			ID:        t.ID,
			TenantID:  t.TenantID,
			Code:      t.Code,
			Channel:   string(t.Channel),
			Subject:   t.Subject,
			Body:      t.Body,
			IsActive:  t.IsActive,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterRoutes registers routes
func RegisterRoutes(e *echo.Echo) {
	// Ensure service is initialized if not already
	if notification.Service == nil {
		notification.Init()
	}

	svc := notification.Service

	nHandler := NewNotificationHandler(svc)
	tHandler := NewTemplateHandler(svc)

	v1 := e.Group("/api/v1/notifications")
	// Add TenantMiddleware to ensure tenant context is present
	v1.Use(middleware.TenantMiddleware)

	v1.POST("/send", nHandler.Send)
	v1.GET("/queue", nHandler.GetQueue)
	v1.POST("/templates", tHandler.Create)
	v1.GET("/templates", tHandler.List)
}
