package handler

import (
	"net/http"

	"github.com/aceextension/core/db"
	"github.com/aceextension/notification/domain"
	"github.com/aceextension/notification/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NotificationHandler handles notification requests
type NotificationHandler struct {
	service service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(service service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

// SendNotificationRequest represents the request body
type SendNotificationRequest struct {
	UserID     *string                `json:"userId"`
	Channel    string                 `json:"channel" validate:"required"`
	Recipient  string                 `json:"recipient" validate:"required"`
	TemplateID *string                `json:"templateId"`
	Content    string                 `json:"content"`
	Priority   string                 `json:"priority"`
	Variables  map[string]interface{} `json:"variables"`
}

// Send sends a notification
// @Summary Send a notification
// @Description Send a notification (instant or queued)
// @Tags notifications
// @Accept json
// @Produce json
// @Param request body SendNotificationRequest true "Notification Request"
// @Success 202 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/notifications/send [post]
// @Security BearerAuth
func (h *NotificationHandler) Send(c echo.Context) error {
	var req SendNotificationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tenantID, ok := db.GetTenantID(c.Request().Context())
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Tenant not found"})
	}

	channel := domain.ChannelType(req.Channel)
	priority := domain.Priority(req.Priority)
	if priority == "" {
		priority = domain.PriorityLow
	}

	serviceReq := service.SendRequest{
		TenantID:  tenantID,
		Channel:   channel,
		Recipient: req.Recipient,
		Content:   req.Content,
		Variables: req.Variables,
		Priority:  priority,
	}

	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err == nil {
			serviceReq.UserID = &uid
		}
	}

	if req.TemplateID != nil {
		tid, err := uuid.Parse(*req.TemplateID)
		if err == nil {
			serviceReq.TemplateID = &tid
		}
	}

	notification, err := h.service.Send(c.Request().Context(), serviceReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"id":     notification.ID.String(),
		"status": string(notification.Status),
	})
}

// GetQueue retrieves pending notifications
// @Summary Get notification queue
// @Description Get pending notifications in the queue
// @Tags notifications
// @Produce json
// @Success 200 {array} domain.Notification
// @Failure 401 {object} map[string]string
// @Router /api/v1/notifications/queue [get]
// @Security BearerAuth
func (h *NotificationHandler) GetQueue(c echo.Context) error {
	notifications, err := h.service.GetPendingNotifications(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, notifications)
}
