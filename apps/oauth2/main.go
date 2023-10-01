package main

import (
	"fmt"
	"github.com/Untanky/modern-auth/internal/app"
	"github.com/Untanky/modern-auth/internal/core"
	ginApp "github.com/Untanky/modern-auth/internal/gin"
	gormLocal "github.com/Untanky/modern-auth/internal/gorm"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"net/http"
	"strings"
)

const (
	ContextPath        = "/api/v1/oauth2"
	CacheControlHeader = "cache-control"
)

var (
	db *gorm.DB
)

func main() {
	err := app.Sequence(
		"Application initialization",
		app.Step("Database initialization", initializeDatabase),
		app.Step("Service initialization", initializeServices),
		app.Step("Gin configuration", ginApp.ConfigureGin),
		app.Step("Telemetry configuration", ginApp.ConfigureTelemetry),
		app.Step("Routing configuration", configureRoutes),
	)
	if err != nil {
		panic(err)
	}

	err = app.AnnounceRun("Application", ginApp.Start)
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
	clientService := oauth2.NewClientService(clientRepo)

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
	authorizationService := oauth2.NewAuthorizationService(authorizationStore, codeStore, authenticationVerifierStore, clientService, authorizationCodeInit, authorizationCodeSuccess)

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
	tokenService := oauth2.NewOAuthTokenService(codeStore, accessTokenHandler, refreshTokenHandler, tokenRequest)

	controllerInstance = newController(authorizationService, clientService, tokenService)

	return nil
}

func configureRoutes() error {
	route := ginApp.GetRouter(ContextPath)

	route.Use(disableCaching)
	route.GET("/authorization", controllerInstance.startAuthorization)
	route.POST("/authorization/succeed", controllerInstance.succeedAuthorization)
	route.POST("/token", controllerInstance.issueToken)
	route.POST("/token/validate", controllerInstance.handleAuthorization, controllerInstance.returnGrant)
	route.GET("/client", controllerInstance.handleAuthorization, controllerInstance.listClients)
	route.GET("/client/:id", controllerInstance.handleAuthorization, controllerInstance.getClient)
	route.POST("/client", controllerInstance.handleAuthorization, controllerInstance.createClient)
	route.DELETE("/client/:id", controllerInstance.handleAuthorization, controllerInstance.deleteClient)

	return nil
}

func disableCaching(c *gin.Context) {
	c.Header(CacheControlHeader, "no-store")
	c.Next()
}

var controllerInstance *controller

type controller struct {
	authorizationService *oauth2.AuthorizationService
	clientService        *oauth2.ClientService
	tokenService         *oauth2.OAuthTokenService
}

func newController(authorizationService *oauth2.AuthorizationService, clientService *oauth2.ClientService, tokenService *oauth2.OAuthTokenService) *controller {
	return &controller{
		authorizationService: authorizationService,
		clientService:        clientService,
		tokenService:         tokenService,
	}
}

func (controller *controller) handleAuthorization(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderParts) != 2 || authorizationHeaderParts[0] != "Bearer" {
		err := fmt.Errorf("invalid authorization header")
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	token := authorizationHeaderParts[1]

	grant, err := controller.tokenService.Validate(ctx.Request.Context(), token)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.Set("grant", grant)
	ctx.Next()
}
