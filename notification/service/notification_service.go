package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aceextension/core/config"
	"github.com/aceextension/notification/domain"
	"github.com/aceextension/notification/repository"
	"github.com/aceextension/notifier"
	"github.com/aceextension/notifier/email"
	"github.com/aceextension/notifier/slack"
	"github.com/google/uuid"
)

type notificationService struct {
	repo         repository.NotificationRepository
	templateRepo repository.TemplateRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(repo repository.NotificationRepository, templateRepo repository.TemplateRepository) NotificationService {
	return &notificationService{
		repo:         repo,
		templateRepo: templateRepo,
	}
}

// Send sends a notification
func (s *notificationService) Send(ctx context.Context, req SendRequest) (*domain.Notification, error) {
	content := req.Content

	// Render template if ID is provided
	if req.TemplateID != nil {
		template, err := s.templateRepo.GetByID(ctx, *req.TemplateID)
		if err != nil {
			return nil, fmt.Errorf("failed to get template: %w", err)
		}
		content = renderTemplate(template.Body, req.Variables)
	}

	// Create notification record
	notification := domain.NewNotification(req.TenantID, req.Channel, req.Recipient, content)
	notification.UserID = req.UserID
	notification.Priority = req.Priority
	notification.TemplateID = req.TemplateID
	notification.Subject = getSubject(req, content)

	if err := s.repo.Create(ctx, notification); err != nil {
		return nil, fmt.Errorf("failed to create notification record: %w", err)
	}

	// Instant trigger for High Priority
	if req.Priority == domain.PriorityHigh {
		// Launch in goroutine to not block, but it's "instant" from user perspective
		// In a real system, we might want to wait for result or use a fast queue
		go func() {
			// Create a background context for async execution
			bgCtx := context.Background()
			if err := s.sendInstant(bgCtx, notification); err != nil {
				log.Printf("Failed to send instant notification %s: %v", notification.ID, err)
			}
		}()
	}

	// Low priority notifications are picked up by the background worker

	return notification, nil
}

// ProcessPending processes pending notifications
func (s *notificationService) ProcessPending(ctx context.Context) error {
	// Fetch pending notifications
	notifications, err := s.repo.GetPending(ctx, 10) // Process batch of 10
	if err != nil {
		return fmt.Errorf("failed to get pending notifications: %w", err)
	}

	for _, n := range notifications {
		if err := s.sendInstant(ctx, n); err != nil {
			log.Printf("Failed to process notification %s: %v", n.ID, err)
		}
	}

	return nil
}

// sendInstant actually sends the notification via provider
func (s *notificationService) sendInstant(ctx context.Context, n *domain.Notification) error {
	// Update status to PROCESSING
	n.Status = domain.StatusProcessing
	if err := s.repo.Update(ctx, n); err != nil {
		return fmt.Errorf("failed to update status to processing: %w", err)
	}

	// 4. Send using real providers
	var ntr notifier.Notifier
	if n.Channel == domain.ChannelEmail {
		ntr = email.New(email.Config{
			SMTPHost: config.GlobalConfig.SMTPHost,
			SMTPPort: config.GlobalConfig.SMTPPort,
			Username: config.GlobalConfig.SMTPUser,
			Password: config.GlobalConfig.SMTPPass,
			From:     config.GlobalConfig.SMTPFrom,
			To:       []string{n.Recipient},
		})
	} else if n.Channel == domain.ChannelSMS {
		// Mock SMS OR fallback to Slack if webhook config exists
		if config.GlobalConfig.SlackWebhookURL != "" {
			ntr = slack.New(config.GlobalConfig.SlackWebhookURL)
		} else if config.GlobalConfig.Env == "development" {
			log.Printf("MOCK SMS to %s: %s", n.Recipient, n.Content)
		}
	}

	if ntr != nil {
		subject := "Update from Practixa"
		if n.Subject != nil {
			subject = *n.Subject
		}
		if err := ntr.Send(subject, n.Content); err != nil {
			n.Status = domain.StatusFailed
			errMsg := err.Error()
			n.ErrorMessage = &errMsg
			s.repo.Update(ctx, n)
			return err
		}
	}

	// Update status to SENT
	n.Status = domain.StatusSent
	now := time.Now()
	n.SentAt = &now

	if err := s.repo.Update(ctx, n); err != nil {
		return fmt.Errorf("failed to update status to sent: %w", err)
	}

	return nil
}

func (s *notificationService) GetTemplates(ctx context.Context, tenantID uuid.UUID) ([]*domain.Template, error) {
	return s.templateRepo.GetByTenantID(ctx, tenantID)
}

func (s *notificationService) CreateTemplate(ctx context.Context, template *domain.Template) error {
	return s.templateRepo.Create(ctx, template)
}

// Simple template renderer {{key}} -> value
func renderTemplate(body string, variables map[string]interface{}) string {
	for k, v := range variables {
		placeholder := fmt.Sprintf("{{%s}}", k)
		body = strings.ReplaceAll(body, placeholder, fmt.Sprintf("%v", v))
	}
	return body
}

func (s *notificationService) GetPendingNotifications(ctx context.Context) ([]*domain.Notification, error) {
	// For inspection, just get first 50 pending items
	return s.repo.GetPending(ctx, 50)
}

func getSubject(req SendRequest, content string) *string {
	// Logic to extract subject, e.g., from template or request
	// For now, return nil or a pointer to a string if provided
	return nil
}
