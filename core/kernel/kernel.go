package kernel

import (
	"fmt"
	"sync"

	"github.com/aceextension/core/extension"
	"github.com/labstack/echo/v4"
)

// Module defines the interface for all AceExtension modules
type Module interface {
	Name() string
	Init() error
	RegisterRoutes(e *echo.Echo, g *echo.Group) error
	RegisterEvents() error
	RegisterPlugins(registry *extension.PluginRegistry) error
}

// Kernel is the core orchestrator of the framework
type Kernel struct {
	mu             sync.RWMutex
	modules        map[string]Module
	pluginManager  *extension.PluginManager
	pluginRegistry *extension.PluginRegistry
}

func NewKernel() *Kernel {
	registry := extension.NewPluginRegistry()
	return &Kernel{
		modules:        make(map[string]Module),
		pluginRegistry: registry,
		pluginManager:  extension.NewPluginManager(registry),
	}
}

func (k *Kernel) RegisterModule(m Module) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	name := m.Name()
	if _, ok := k.modules[name]; ok {
		return fmt.Errorf("module %s already registered", name)
	}

	k.modules[name] = m
	return nil
}

func (k *Kernel) Boot(e *echo.Echo) error {
	k.mu.RLock()
	defer k.mu.RUnlock()

	// 1. Initialize all modules
	for name, m := range k.modules {
		fmt.Printf("Initializing module: %s\n", name)
		if err := m.Init(); err != nil {
			return fmt.Errorf("failed to init module %s: %w", name, err)
		}
	}

	// 2. Register all plugins
	for name, m := range k.modules {
		fmt.Printf("Registering plugins for: %s\n", name)
		if err := m.RegisterPlugins(k.pluginRegistry); err != nil {
			return fmt.Errorf("failed to register plugins for %s: %w", name, err)
		}
	}

	// 3. Register all routes
	apiGroup := e.Group("/api/v1")
	for name, m := range k.modules {
		fmt.Printf("Registering routes for: %s\n", name)
		if err := m.RegisterRoutes(e, apiGroup); err != nil {
			return fmt.Errorf("failed to register routes for %s: %w", name, err)
		}
	}

	return nil
}

func (k *Kernel) GetPluginManager() *extension.PluginManager {
	return k.pluginManager
}
