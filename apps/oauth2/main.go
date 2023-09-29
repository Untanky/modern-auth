package main

import (
	"github.com/Untanky/modern-auth/internal/app"
	"github.com/gin-gonic/gin"
	"log/slog"
)

const (
	ContextPath        = "/api/v1/oauth2"
	CacheControlHeader = "cache-control"
)

func main() {
	slog.Info("Application initialization starting")

	slog.Debug("Gin configuration starting")
	engine := app.ConfigureGin()
	slog.Info("Gin configuration successful")
	route := engine.Group(ContextPath)

	slog.Debug("Telemetry configuration starting")
	err := app.ConfigureTelemetry(engine)
	if err != nil {
		slog.Error("Telemetry configuration failed", "error", err)
		panic(err)
	}
	slog.Info("Telemetry configuration successful")

	slog.Debug("Routing configuration starting")
	configureRoutes(route)
	slog.Info("Routing configuration successful")

	slog.Info("Application initialization successful")
	slog.Info("Application now starting to listen")

	err = engine.Run(":8080")
	if err != nil {
		slog.Error("Application did not run", "err", err)
	}
}

func configureRoutes(route gin.IRouter) {
	route.Use(disableCaching)
	route.GET("/authorization", startAuthorization)
	route.POST("/authorization/succeed", succeedAuthorization)
	route.POST("/token", issueToken)
	route.POST("/token/validate", validateToken)
	route.GET("/client")
	route.GET("/client/:id")
	route.POST("/client")
	route.DELETE("/client/:id")
}

func disableCaching(c *gin.Context) {
	c.Header(CacheControlHeader, "no-store")
	c.Next()
}
