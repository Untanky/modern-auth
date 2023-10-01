package main

//import (
//	"log"
//	"log/slog"
//
//	"github.com/Untanky/modern-auth/apps/oauth2/internal/oauth2"
//	"github.com/Untanky/modern-auth/internal/core"
//	"github.com/Untanky/modern-auth/internal/domain"
//	gormLocal "github.com/Untanky/modern-auth/internal/gorm"
//	"github.com/gin-gonic/gin"
//	"github.com/google/uuid"
//	"github.com/prometheus/client_golang/prometheus/promhttp"
//	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	"gorm.io/plugin/opentelemetry/tracing"
//)
//
//type App struct {
//	db     *gorm.DB
//	engine *gin.Engine
//}
//
//type entitiesKey string
//
//var EntitiesKey entitiesKey = "entities"
//
//type WriteFunc func([]byte) (int, error)
//
//func (fn WriteFunc) Write(data []byte) (int, error) {
//	return fn(data)
//}
//
//func (a *App) Start() {
//	slog.Info("Application initialization starting")
//
//	gin.DefaultWriter = WriteFunc(func(data []byte) (int, error) {
//		slog.Debug(string(data))
//		return 0, nil
//	})
//
//	a.db = a.connect()
//	a.migrateEntities([]interface{}{
//		oauth2.ClientModel{},
//		gormLocal.User{},
//		gormLocal.Credential{},
//	})
//
//	requestMetrics, err := newRequestTelemetry(meter)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	slog.Debug("Initialize services starting")
//	authenticationVerifierStore := core.NewInMemoryKeyValueStore[[]byte]()
//
//	initAuthnStore := core.NewInMemoryKeyValueStore[webauthn.CredentialOptions]()
//	userRepo := gormLocal.NewGormUserRepo(a.db)
//	userService := domain.NewUserService(userRepo)
//	credentialRepo := gormLocal.NewGormCredentialRepo(a.db)
//	credentialService := domain.NewCredentialService(credentialRepo)
//	authenticationService := webauthn.NewAuthenticationService(initAuthnStore, authenticationVerifierStore, userService, credentialService)
//	authenticationController := webauthn.NewAuthenticationController(authenticationService)
//	slog.Info("Initialize services successful")
//
//	// gin.SetMode(gin.ReleaseMode)
//	r := gin.New()
//	a.engine = r
//
//	r.Use(gin.Recovery())
//	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
//
//	r.Use(a.handleRequestId)
//	r.Use(otelgin.Middleware("modern-auth"))
//	r.Use(requestMetrics.handleTelemetry())
//
//	api := r.Group("/v1")
//	slog.Debug("Router setup starting")
//	authenticationController.RegisterRoutes(api.Group("/webauthn"))
//	slog.Info("Router setup successful")
//
//	slog.Info("Application initialization successful")
//	slog.Info("Application starting to listen")
//	r.Run(":8080")
//}
//
//func (a *App) handleRequestId(c *gin.Context) {
//	requestId := c.Request.Header.Get("Request-Id")
//	if requestId == "" {
//		requestId = uuid.New().String()
//	}
//
//	c.Set("requestId", requestId)
//	c.Writer.Header().Set("Request-Id", c.GetString("requestId"))
//
//	c.Next()
//}
//
//func (a *App) connect() *gorm.DB {
//	slog.Debug("Database connection starting")
//	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Berlin"
//	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect database")
//	}
//	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
//		panic(err)
//	}
//	slog.Info("Database connection successful")
//	return db
//}
//
//func (a *App) migrateEntities(entities []interface{}) {
//	slog.Debug("Entity migration starting")
//	for _, entity := range entities {
//		err := a.db.AutoMigrate(entity)
//		if err != nil {
//			panic("failed to migrate entity")
//		}
//		slog.Info("Entity migration successful")
//	}
//}
//
//func (a *App) Stop() {
//	slog.Info("Application stopping")
//	db, _ := a.db.DB()
//	if db != nil {
//		db.Close()
//	}
//}
