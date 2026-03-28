package accounting

import (
	"github.com/aceextension/accounting/handler"
	"github.com/aceextension/core/extension"
	"github.com/labstack/echo/v4"
)

type AccountingModule struct {
}

func NewAccountingModule() *AccountingModule {
	return &AccountingModule{}
}

func (m *AccountingModule) Name() string {
	return "accounting"
}

func (m *AccountingModule) Init() error {
	Init()
	return nil
}

func (m *AccountingModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	accAccountHandler := handler.NewAccountHandler(Service)
	accJournalHandler := handler.NewJournalHandler(Service)
	accReportHandler := handler.NewReportHandler(Service)

	handler.RegisterRoutes(g, accAccountHandler, accJournalHandler, accReportHandler)
	return nil
}

func (m *AccountingModule) RegisterEvents() error {
	return nil
}

func (m *AccountingModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
