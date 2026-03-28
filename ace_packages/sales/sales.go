package sales

import (
	"github.com/aceextension/sales/service"
)

var (
	Service service.SalesService
)

// Init initializes the sales module
// For now, dependencies are injected in main.go, so this might be empty
// or used for other setup.
func Init(s service.SalesService) {
	Service = s
}
