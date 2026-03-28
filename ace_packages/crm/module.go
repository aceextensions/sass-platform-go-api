package crm

import (
	"github.com/aceextension/core/extension"
	"github.com/aceextension/crm/handler"
	"github.com/aceextension/notification"
	"github.com/labstack/echo/v4"
)

type CRMModule struct {
}

func NewCRMModule() *CRMModule {
	return &CRMModule{}
}

func (m *CRMModule) Name() string {
	return "crm"
}

func (m *CRMModule) Init() error {
	Init(notification.Service)
	return nil
}

func (m *CRMModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	handler.RegisterRoutes(e, CustomerService, SupplierService)
	return nil
}

func (m *CRMModule) RegisterEvents() error {
	return nil
}

func (m *CRMModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
