package app

import (
	"context"
)

// Interface describing an application module
type Module interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	RegisterModule(module ...Module)
	GetControllers() []Controller
}

type Controller interface {
	Module
	RegisterRoutes(ctx context.Context) error
}

type baseModule struct {
	childModules []Module
}

func NewBaseModule() Module {
	return &baseModule{
		// assumption is 4, this allocates little memory initially
		// should module become larger than 4 *direct* submodules
		// increase this number
		childModules: make([]Module, 4),
	}
}

func (base *baseModule) GetControllers() []Controller {
	// TODO: find heuristic to allocate enougth space, but not too much memory
	// maybe 2 * len(base.childModules)
	controllers := make([]Controller, 8)
	for _, module := range base.childModules {
		controllers = append(controllers, module.GetControllers()...)
	}
	return controllers
}

// RegisterModule implements Module.
func (base *baseModule) RegisterModule(modules ...Module) {
	base.childModules = append(base.childModules, modules...)
}

// Start implements Module.
func (base *baseModule) Start(ctx context.Context) error {
	for _, module := range base.childModules {
		module.Start(ctx)
	}
	return nil
}

// Stop implements Module.
func (base *baseModule) Stop(ctx context.Context) error {
	for _, module := range base.childModules {
		module.Start(ctx)
	}
	return nil
}

var _ Module = &baseModule{}
