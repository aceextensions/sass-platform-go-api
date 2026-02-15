package subscription

import (
	"github.com/aceextension/core/db"
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

	planRepo := repository.NewPostgresPlanRepository(db.MainPool)
	subRepo := repository.NewPostgresSubscriptionRepository(db.MainPool)
	Service = service.NewSubscriptionService(planRepo, subRepo)
}
