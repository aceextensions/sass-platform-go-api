package catalog

import (
	"github.com/aceextension/catalog/handler"
	"github.com/aceextension/core/extension"
	"github.com/labstack/echo/v4"
)

type CatalogModule struct {
}

func NewCatalogModule() *CatalogModule {
	return &CatalogModule{}
}

func (m *CatalogModule) Name() string {
	return "catalog"
}

func (m *CatalogModule) Init() error {
	Init()
	return nil
}

func (m *CatalogModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	handler.RegisterRoutes(g, CategoryService, ProductService)
	return nil
}

func (m *CatalogModule) RegisterEvents() error {
	return nil
}

func (m *CatalogModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
