package sales

import (
	"github.com/aceextension/accounting"
	"github.com/aceextension/core/extension"
	"github.com/aceextension/inventory"
	"github.com/aceextension/sales/handler"
	"github.com/aceextension/sales/repository"
	"github.com/aceextension/sales/service"
	"github.com/labstack/echo/v4"
)

type SalesModule struct {
}

func NewSalesModule() *SalesModule {
	return &SalesModule{}
}

func (m *SalesModule) Name() string {
	return "sales"
}

func (m *SalesModule) Init() error {
	sRepo := repository.NewPostgresSalesRepository()
	sService := service.NewSalesService(sRepo, inventory.Service, accounting.Service)
	Init(sService)
	return nil
}

func (m *SalesModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	handler.RegisterRoutes(g, handler.NewSalesHandler(Service))
	return nil
}

func (m *SalesModule) RegisterEvents() error {
	return nil
}

func (m *SalesModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
