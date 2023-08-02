package main

import (
	"log"

	"github.com/Untanky/modern-auth/internal/core"
	gormLocal "github.com/Untanky/modern-auth/internal/gorm"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/Untanky/modern-auth/internal/user"
	"github.com/Untanky/modern-auth/internal/webauthn"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type App struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
	engine *gin.Engine
}

var perfLogger *zap.Logger
var logger *zap.SugaredLogger

func init() {
	traceExporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		log.Fatal(err)
	}
	mergedResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ModernAuth"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(mergedResource),
	)

	otel.SetTracerProvider(tracerProvider)

	meterExporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	meterProvider := metric.NewMeterProvider(metric.WithReader(meterExporter))

	otel.SetMeterProvider(meterProvider)

	perfLogger, _ = zap.NewDevelopment()
	logger = perfLogger.Sugar()
}

type entitiesKey string

var EntitiesKey entitiesKey = "entities"

func (a *App) Start() {
	a.logger = logger
	logger.Info("Application initialization starting")

	a.db = a.connect()
	a.migrateEntities([]interface{}{
		oauth2.ClientModel{},
		gormLocal.User{},
		gormLocal.Credential{},
	})

	meter := otel.GetMeterProvider().Meter("github.com/Untanky/modern-auth")

	requestMetrics, err := newRequestTelemetry(meter, perfLogger.Named("RequestTelemetry"))
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("Initialize services starting")
	clientRepo := gormLocal.NewGormRepository[string, *oauth2.ClientModel, *oauth2.ClientModel](
		a.db,
		func(a *oauth2.ClientModel) *oauth2.ClientModel {
			return a
		},
		func(a *oauth2.ClientModel) *oauth2.ClientModel {
			return a
		},
	)
	clientService := oauth2.NewClientService(clientRepo, logger.Named("ClientService"))
	clientController := oauth2.NewClientController(clientService)

	authorizationStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationRequest]()
	codeStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationRequest]()
	authorizationCodeInit, err := meter.Int64Counter("authorization_code_init")
	if err != nil {
		log.Fatal(err)
	}
	authorizationCodeSuccess, err := meter.Int64Counter("authorization_code_success")
	if err != nil {
		log.Fatal(err)
	}
	authorizationService := oauth2.NewAuthorizationService(authorizationStore, codeStore, clientService, logger.Named("AuthorizationService"), authorizationCodeInit, authorizationCodeSuccess)
	authorizationController := oauth2.NewAuthorizationController(authorizationService)

	accessTokenStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationGrant]()
	accessTokensGenerated, err := meter.Int64Counter("access_tokens_generated")
	if err != nil {
		log.Fatal(err)
	}
	accessTokenHandler := oauth2.NewRandomTokenHandler(48, accessTokenStore, logger.Named("AccessTokenHandler"), accessTokensGenerated)
	refreshTokenStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationGrant]()
	refreshTokensGenerated, err := meter.Int64Counter("refresh_tokens_generated")
	if err != nil {
		log.Fatal(err)
	}
	refreshTokenHandler := oauth2.NewRandomTokenHandler(64, refreshTokenStore, logger.Named("RefreshTokenHandler"), refreshTokensGenerated)
	tokenRequest, err := meter.Int64Counter("token_request")
	if err != nil {
		log.Fatal(err)
	}
	oauthTokenService := oauth2.NewOAuthTokenService(codeStore, accessTokenHandler, refreshTokenHandler, logger.Named("TokenService"), tokenRequest)
	tokenController := oauth2.NewTokenController(oauthTokenService)

	initAuthnStore := core.NewInMemoryKeyValueStore[webauthn.InitiateAuthenticationResponse]()
	userRepo := gormLocal.NewGormUserRepo(a.db)
	userService := user.NewUserService(userRepo)
	credentialRepo := gormLocal.NewGormCredentialRepo(a.db)
	credentialService := user.NewCredentialService(credentialRepo)
	authenticationService := webauthn.NewAuthenticationService(initAuthnStore, userService, credentialService)
	authenticationController := webauthn.NewAuthenticationController(authenticationService)
	logger.Info("Initialize services successful")

	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	a.engine = r

	r.Use(gin.Recovery())
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Use(a.handleRequestId)
	r.Use(otelgin.Middleware("modern-auth"))
	r.Use(requestMetrics.handleTelemetry())

	api := r.Group("/v1")
	logger.Debug("Router setup starting")
	clientController.RegisterRoutes(api.Group("/client"))
	oauth2Router := api.Group("/oauth2")
	authorizationController.RegisterRoutes(oauth2Router)
	tokenController.RegisterRoutes(oauth2Router)
	authenticationController.RegisterRoutes(api.Group("/webauthn"))
	logger.Info("Router setup successful")

	logger.Info("Application initialization successful")
	logger.Info("Application starting to listen")
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
	logger.Debug("Database connection starting")
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Berlin"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	logger.Info("Database connection successful")
	return db
}

func (a *App) migrateEntities(entities []interface{}) {
	logger.Debug("Entity migration starting")
	for _, entity := range entities {
		err := a.db.AutoMigrate(entity)
		if err != nil {
			panic("failed to migrate entity")
		}
		logger.Info("Entity migration successful")
	}
}

func (a *App) Stop() {
	logger.Info("Application stopping")
	db, _ := a.db.DB()
	if db != nil {
		db.Close()
	}
}
