package main

import (
	"log"
	"log/slog"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/domain"
	gormLocal "github.com/Untanky/modern-auth/internal/gorm"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/Untanky/modern-auth/internal/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type App struct {
	db     *gorm.DB
	engine *gin.Engine
}

type entitiesKey string

var EntitiesKey entitiesKey = "entities"

type WriteFunc func([]byte) (int, error)

func (fn WriteFunc) Write(data []byte) (int, error) {
	return fn(data)
}

func (a *App) Start() {
	slog.Info("Application initialization starting")

	gin.DefaultWriter = WriteFunc(func(data []byte) (int, error) {
		slog.Debug(string(data))
		return 0, nil
	})

	a.db = a.connect()
	a.migrateEntities([]interface{}{
		oauth2.ClientModel{},
		gormLocal.User{},
		gormLocal.Credential{},
	})

	requestMetrics, err := newRequestTelemetry(meter)
	if err != nil {
		log.Fatal(err)
	}

	slog.Debug("Initialize services starting")
	clientRepo := gormLocal.NewGormRepository[string, *oauth2.ClientModel, *oauth2.ClientModel](
		a.db,
		func(a *oauth2.ClientModel) *oauth2.ClientModel {
			return a
		},
		func(a *oauth2.ClientModel) *oauth2.ClientModel {
			return a
		},
	)
	clientService := oauth2.NewClientService(clientRepo)
	clientController := oauth2.NewClientController(clientService)

	authorizationStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationRequest]()
	codeStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationRequest]()
	authorizationCodeInit, err := meter.Int64Counter("authorization_code_init")
	if err != nil {
		log.Fatal(err)
	}
	authorizationCodeSuccess, err := meter.Int64Counter("authorization_code_success")
	if err != nil {
		log.Fatal(err)
	}
	authorizationService := oauth2.NewAuthorizationService(authorizationStore, codeStore, clientService, authorizationCodeInit, authorizationCodeSuccess)
	authorizationController := oauth2.NewAuthorizationController(authorizationService)

	accessTokenStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationGrant]()
	accessTokensGenerated, err := meter.Int64Counter("access_tokens_generated")
	if err != nil {
		log.Fatal(err)
	}
	accessTokenHandler := oauth2.NewRandomTokenHandler("access-token", 48, accessTokenStore, accessTokensGenerated)
	refreshTokenStore := core.NewInMemoryKeyValueStore[*oauth2.AuthorizationGrant]()
	refreshTokensGenerated, err := meter.Int64Counter("refresh_tokens_generated")
	if err != nil {
		log.Fatal(err)
	}
	refreshTokenHandler := oauth2.NewRandomTokenHandler("refresh-token", 64, refreshTokenStore, refreshTokensGenerated)
	tokenRequest, err := meter.Int64Counter("token_request")
	if err != nil {
		log.Fatal(err)
	}
	oauthTokenService := oauth2.NewOAuthTokenService(codeStore, accessTokenHandler, refreshTokenHandler, tokenRequest)
	tokenController := oauth2.NewTokenController(oauthTokenService)

	initAuthnStore := core.NewInMemoryKeyValueStore[webauthn.CredentialOptions]()
	userRepo := gormLocal.NewGormUserRepo(a.db)
	userService := domain.NewUserService(userRepo)
	credentialRepo := gormLocal.NewGormCredentialRepo(a.db)
	credentialService := domain.NewCredentialService(credentialRepo)
	authenticationService := webauthn.NewAuthenticationService(initAuthnStore, userService, credentialService)
	authenticationController := webauthn.NewAuthenticationController(authenticationService)
	slog.Info("Initialize services successful")

	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	a.engine = r

	r.Use(gin.Recovery())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(a.handleRequestId)
	r.Use(otelgin.Middleware("modern-auth"))
	r.Use(requestMetrics.handleTelemetry())

	api := r.Group("/v1")
	slog.Debug("Router setup starting")
	clientController.RegisterRoutes(api.Group("/client"))
	oauth2Router := api.Group("/oauth2")
	authorizationController.RegisterRoutes(oauth2Router)
	tokenController.RegisterRoutes(oauth2Router)
	authenticationController.RegisterRoutes(api.Group("/webauthn"))
	slog.Info("Router setup successful")

	slog.Info("Application initialization successful")
	slog.Info("Application starting to listen")
	r.Run(":8080")
}

func (a *App) handleRequestId(c *gin.Context) {
	requestId := c.Request.Header.Get("Request-Id")
	if requestId == "" {
		requestId = uuid.New().String()
	}

	c.Set("requestId", requestId)
	c.Writer.Header().Set("Request-Id", c.GetString("requestId"))

	c.Next()
}

func (a *App) connect() *gorm.DB {
	slog.Debug("Database connection starting")
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	slog.Info("Database connection successful")
	return db
}

func (a *App) migrateEntities(entities []interface{}) {
	slog.Debug("Entity migration starting")
	for _, entity := range entities {
		err := a.db.AutoMigrate(entity)
		if err != nil {
			panic("failed to migrate entity")
		}
		slog.Info("Entity migration successful")
	}
}

func (a *App) Stop() {
	slog.Info("Application stopping")
	db, _ := a.db.DB()
	if db != nil {
		db.Close()
	}
}
