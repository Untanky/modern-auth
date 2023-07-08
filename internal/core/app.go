package core

import (
	"context"

	"go.uber.org/zap"
)

type App struct {
	module Module
}

func NewApp(module Module) *App {
	return &App{module: module}
}

var logger *zap.SugaredLogger

func init() {
	log, _ := zap.NewDevelopment()
	logger = log.Sugar()
}

type entitiesKey string

var EntitiesKey entitiesKey = "entities"

func (a *App) Start() {
	logger.Info("Application starting")

	moduleCtx := context.Background()

	entities := a.module.GetEntities()
	moduleCtx = context.WithValue(moduleCtx, EntitiesKey, entities)

	logger.Debug("Module initialization starting")
	a.module.Init(moduleCtx)
	logger.Info("Module initialization successful")
}

func (a *App) Stop() {
	a.module.Shutdown()
}
