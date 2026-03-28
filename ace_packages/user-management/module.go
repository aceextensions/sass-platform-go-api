package usermanagement

import (
	"github.com/aceextension/core/extension"
	"github.com/labstack/echo/v4"
)

type UserManagementModule struct {
}

func NewUserManagementModule() *UserManagementModule {
	return &UserManagementModule{}
}

func (m *UserManagementModule) Name() string {
	return "user-management"
}

func (m *UserManagementModule) Init() error {
	return nil
}

func (m *UserManagementModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	return nil
}

func (m *UserManagementModule) RegisterEvents() error {
	return nil
}

func (m *UserManagementModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
