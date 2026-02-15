package notification

import (
	"github.com/aceextension/notification/repository"
	"github.com/aceextension/notification/service"
)

var (
	// TemplateRepo instance
	TemplateRepo repository.TemplateRepository
	// NotificationRepo instance
	NotificationRepo repository.NotificationRepository
	// Service instance
	Service service.NotificationService
)

// Init initializes the notification module
func Init() {
	TemplateRepo = repository.NewPostgresTemplateRepository()
	NotificationRepo = repository.NewPostgresNotificationRepository()
	Service = service.NewNotificationService(NotificationRepo, TemplateRepo)
}
