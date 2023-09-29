package main

import (
	"github.com/Untanky/modern-auth/internal/app"
	"github.com/Untanky/modern-auth/internal/core"
	gormLocal "github.com/Untanky/modern-auth/internal/gorm"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

const (
	ContextPath        = "/api/v1/oauth2"
	CacheControlHeader = "cache-control"
)

var (
	route gin.IRouter
	db    *gorm.DB
)

func main() {
	err := app.Sequence(
		"Application initialization",
		app.Step("Database initialization", initializeDatabase),
		app.Step("Service initialization", initializeServices),
		app.Step("Gin configuration", app.ConfigureGin),
		app.Step("Telemetry configuration", app.ConfigureTelemetry),
		func() error {
			route = app.GetRouter(ContextPath)
			return nil
		},
		app.Step("Routing configuration", configureRoutes),
	)
	if err != nil {
		panic(err)
	}

	err = app.AnnounceRun("Application", app.Start)
	if err != nil {
		panic(err)
	}
}

func initializeDatabase() error {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Berlin"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		return err
	}

	return nil
}

func initializeServices() error {
	clientRepo := gormLocal.NewGormRepository[string, *oauth2.ClientModel, *oauth2.ClientModel](
		db,
		func(a *oauth2.ClientModel) *oauth2.ClientModel {
			return a
		},
		func(a *oauth2.ClientModel) *oauth2.ClientModel {
			return a
		},
	)
	clientService = oauth2.NewClientService(clientRepo)

	authenticationVerifierStore := core.NewInMemoryKeyValueStore[[]byte]()

	meter := otel.GetMeterProvider().Meter("github.com/Untanky/modern-auth/oauth2")

	authorizationStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationRequest]()
	codeStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationRequest]()
	authorizationCodeInit, err := meter.Int64Counter("authorization_code_init")
	if err != nil {
		return err
	}
	authorizationCodeSuccess, err := meter.Int64Counter("authorization_code_success")
	if err != nil {
		return err
	}
	authorizationService = oauth2.NewAuthorizationService(authorizationStore, codeStore, authenticationVerifierStore, clientService, authorizationCodeInit, authorizationCodeSuccess)

	accessTokenStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationGrant]()
	accessTokensGenerated, err := meter.Int64Counter("access_tokens_generated")
	if err != nil {
		return err
	}
	accessTokenHandler := oauth2.NewRandomTokenHandler("access-token", 48, accessTokenStore, accessTokensGenerated)
	refreshTokenStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationGrant]()
	refreshTokensGenerated, err := meter.Int64Counter("refresh_tokens_generated")
	if err != nil {
		return err
	}
	refreshTokenHandler := oauth2.NewRandomTokenHandler("refresh-token", 64, refreshTokenStore, refreshTokensGenerated)
	tokenRequest, err := meter.Int64Counter("token_request")
	if err != nil {
		return err
	}
	tokenService = oauth2.NewOAuthTokenService(codeStore, accessTokenHandler, refreshTokenHandler, tokenRequest)

	return nil
}

func configureRoutes() error {
	route.Use(disableCaching)
	route.GET("/authorization", startAuthorization)
	route.POST("/authorization/succeed", succeedAuthorization)
	route.POST("/token", issueToken)
	route.POST("/token/validate", validateToken)
	route.GET("/client")
	route.GET("/client/:id")
	route.POST("/client")
	route.DELETE("/client/:id")

	return nil
}

func disableCaching(c *gin.Context) {
	c.Header(CacheControlHeader, "no-store")
	c.Next()
}
