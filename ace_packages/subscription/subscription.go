package subscription

import (
	"github.com/aceextension/subscription/repository"
	"github.com/aceextension/subscription/service"
)

var (
	Service service.SubscriptionService
)

// Init initializes the subscription module
func Init() {
	if Service != nil {
		return
	}

	planRepo := repository.NewPostgresPlanRepository()
	subRepo := repository.NewPostgresSubscriptionRepository()
	Service = service.NewSubscriptionService(planRepo, subRepo)
}
