package purchase

import (
	"github.com/aceextension/core/extension"
	"github.com/aceextension/inventory"
	"github.com/aceextension/accounting"
	"github.com/aceextension/purchase/handler"
	"github.com/aceextension/purchase/repository"
	"github.com/aceextension/purchase/service"
	"github.com/labstack/echo/v4"
)

type PurchaseModule struct {
}

func NewPurchaseModule() *PurchaseModule {
	return &PurchaseModule{}
}

func (m *PurchaseModule) Name() string {
	return "purchase"
}

func (m *PurchaseModule) Init() error {
	pRepo := repository.NewPostgresPurchaseRepository()
	pService := service.NewPurchaseService(pRepo, inventory.Service, accounting.Service)
	Init(pService)
	return nil
}

func (m *PurchaseModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	handler.RegisterRoutes(g, handler.NewPurchaseHandler(Service))
	return nil
}

func (m *PurchaseModule) RegisterEvents() error {
	return nil
}

func (m *PurchaseModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
