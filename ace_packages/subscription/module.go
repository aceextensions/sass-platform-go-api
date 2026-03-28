package subscription

import (
	"github.com/aceextension/core/extension"
	"github.com/aceextension/subscription/handler"
	"github.com/labstack/echo/v4"
)

type SubscriptionModule struct {
}

func NewSubscriptionModule() *SubscriptionModule {
	return &SubscriptionModule{}
}

func (m *SubscriptionModule) Name() string {
	return "subscription"
}

func (m *SubscriptionModule) Init() error {
	Init()
	return nil
}

func (m *SubscriptionModule) RegisterRoutes(e *echo.Echo, g *echo.Group) error {
	subPlanHandler := handler.NewPlanHandler(Service)
	subHandler := handler.NewSubscriptionHandler(Service, nil) // Authentication service will be injected via plugin or shared kernel state later

	plans := g.Group("/plans")
	plans.POST("", subPlanHandler.Create)
	plans.GET("", subPlanHandler.List)

	subs := g.Group("/subscriptions")
	subs.GET("/current", subHandler.GetCurrentSubscription)
	subs.POST("/subscribe", subHandler.Subscribe)

	return nil
}

func (m *SubscriptionModule) RegisterEvents() error {
	return nil
}

func (m *SubscriptionModule) RegisterPlugins(registry *extension.PluginRegistry) error {
	return nil
}
