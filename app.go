package main

import (
	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
	engine *gin.Engine
}

var perfLogger *zap.Logger
var logger *zap.SugaredLogger

func init() {
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
	})

	logger.Debug("Initialize services starting")
	clientRepo := core.NewGormRepository[string, *oauth2.ClientModel](a.db)
	clientService := oauth2.NewClientService(clientRepo, logger.Named("ClientService"))
	clientController := oauth2.NewClientController(clientService)

	authorizationStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationRequest]()
	codeStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationRequest]()
	authorizationService := oauth2.NewAuthorizationService(authorizationStore, codeStore, clientService, logger.Named("AuthorizationService"))
	authorizationController := oauth2.NewAuthorizationController(authorizationService)

	accessTokenStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationGrant]()
	accessTokenHandler := oauth2.NewRandomTokenHandler(48, accessTokenStore, logger.Named("AccessTokenHandler"))
	refreshTokenStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationGrant]()
	refreshTokenHandler := oauth2.NewRandomTokenHandler(64, refreshTokenStore, logger.Named("RefreshTokenHandler"))
	oauthTokenService := oauth2.NewOAuthTokenService(codeStore, accessTokenHandler, refreshTokenHandler, logger.Named("TokenService"))
	tokenController := oauth2.NewTokenController(oauthTokenService)
	logger.Info("Initialize services successful")

	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(a.handleRequestId)
	api := r.Group("/v1")
	logger.Debug("Router setup starting")
	clientController.RegisterRoutes(api.Group("/client"))
	oauth2Router := api.Group("/oauth2")
	authorizationController.RegisterRoutes(oauth2Router)
	tokenController.RegisterRoutes(oauth2Router)
	logger.Info("Router setup successful")

	logger.Info("Application initialization successful")
	logger.Info("Application starting to listen")
	a.engine = r
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
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
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
