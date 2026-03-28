package fiscal

import (
	"github.com/aceextension/fiscal/repository"
	"github.com/aceextension/fiscal/service"
)

// Global fiscal year service instance
var Service service.FiscalYearService

// Init initializes the fiscal module
func Init() {
	repo := repository.NewPostgresFiscalYearRepository()
	Service = service.NewFiscalYearService(repo)
}
