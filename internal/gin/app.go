package gin

import (
	"context"

	"log/slog"

	. "github.com/Untanky/modern-auth/internal/app"

	"github.com/gin-gonic/gin"
)

type WriteFunc func([]byte) (int, error)

func (fn WriteFunc) Write(data []byte) (int, error) {
	return fn(data)
}

type ginApplication struct {
	modules []Module
}

func NewGinApplication() Module {
	return &ginApplication{}
}

type ginKey string

const ginEngine = ginKey("ginEngine")

func (app *ginApplication) Start(ctx context.Context) error {
	slog.Info("Application initialization starting")
	slog.Debug("Registered modules starting")
	for _, module := range app.modules {
		module.Start(ctx)
	}
	slog.Info("Registered modules started successful")

	gin.DefaultWriter = WriteFunc(func(data []byte) (int, error) {
		slog.Debug(string(data))
		return 0, nil
	})
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	ctx = context.WithValue(ctx, ginEngine, r)

	slog.Debug("Registration of controller routes starting")
	for _, controller := range app.GetControllers() {
		controller.RegisterRoutes(ctx)
	}
	slog.Debug("Registration of controller routes successful")
	slog.Info("Application initialization successful")

	slog.Info("Application starting to listen")
	return r.Run()
}

func (app *ginApplication) Stop(ctx context.Context) error {
	slog.Debug("Registered modules stopping")
	for _, module := range app.modules {
		module.Start(ctx)
	}
	slog.Info("Registered modules stopped successful")
	return nil
}

func (app *ginApplication) RegisterModule(module ...Module) {
	app.modules = append(app.modules, module...)
}

func (app *ginApplication) GetControllers() []Controller {
	return nil
}
