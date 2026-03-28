package purchase

import (
	"github.com/aceextension/purchase/service"
)

var (
	Service service.PurchaseService
)

func Init(s service.PurchaseService) {
	Service = s
}
