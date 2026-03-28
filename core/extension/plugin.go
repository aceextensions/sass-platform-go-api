package extension

import (
	"context"
	"fmt"
	"sync"
)

// PluginHook defines all extension points in the system
type PluginHook string

const (
	// User Hooks
	BeforeUserCreate PluginHook = "user:before_create"
	AfterUserCreate  PluginHook = "user:after_create"
	BeforeUserUpdate PluginHook = "user:before_update"
	AfterUserUpdate  PluginHook = "user:after_update"

	// Tenant Hooks
	BeforeTenantCreate PluginHook = "tenant:before_create"
	AfterTenantCreate  PluginHook = "tenant:after_create"

	// Auth Hooks
	BeforeLogin  PluginHook = "auth:before_login"
	AfterLogin   PluginHook = "auth:after_login"
	BeforeLogout PluginHook = "auth:before_logout"

	// Data Hooks (generic)
	BeforeCreate PluginHook = "data:before_create"
	AfterCreate  PluginHook = "data:after_create"
)

// Plugin defines a pluggable function
type Plugin interface {
	Name() string
	Version() string
	Execute(ctx context.Context, payload map[string]interface{}) (map[string]interface{}, error)
}

// PluginRegistry manages all plugins
type PluginRegistry struct {
	mu      sync.RWMutex
	plugins map[PluginHook][]Plugin
}

func NewPluginRegistry() *PluginRegistry {
	return &PluginRegistry{
		plugins: make(map[PluginHook][]Plugin),
	}
}

func (r *PluginRegistry) Register(hook PluginHook, plugin Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[hook] = append(r.plugins[hook], plugin)
	return nil
}

func (r *PluginRegistry) Execute(ctx context.Context, hook PluginHook, payload map[string]interface{}) (map[string]interface{}, error) {
	r.mu.RLock()
	plugins, ok := r.plugins[hook]
	r.mu.RUnlock()

	if !ok {
		return payload, nil
	}

	for _, plugin := range plugins {
		result, err := plugin.Execute(ctx, payload)
		if err != nil {
			return nil, fmt.Errorf("plugin %s failed: %w", plugin.Name(), err)
		}
		payload = result
	}
	return payload, nil
}

// PluginManager is a high-level orchestrator
type PluginManager struct {
	Registry *PluginRegistry
}

func NewPluginManager(registry *PluginRegistry) *PluginManager {
	return &PluginManager{
		Registry: registry,
	}
}

func (m *PluginManager) Execute(ctx context.Context, hook PluginHook, payload map[string]interface{}) (map[string]interface{}, error) {
	return m.Registry.Execute(ctx, hook, payload)
}
