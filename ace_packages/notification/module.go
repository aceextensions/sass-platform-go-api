package notification

import (
	"github.com/aceextension/core/extension"
	"github.com/aceextension/notification/handler"
	"github.com/labstack/echo/v4"
)

type NotificationModule struct {
}

func NewNotificationModule() *NotificationModule {
	return &NotificationModule{}
}

func (m *NotificationModule) Name() string {
	return "notification"
}

func (m *NotificationModule) Init() error {
	Init()
	return nil
}

func (m *NotificationModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	handler.RegisterRoutes(e, Service)
	return nil
}

func (m *NotificationModule) RegisterEvents() error {
	return nil
}

func (m *NotificationModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
