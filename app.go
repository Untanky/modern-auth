package main

import (
	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
	engine *gin.Engine
}

var logger *zap.SugaredLogger

func init() {
	log, _ := zap.NewDevelopment()
	logger = log.Sugar()
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
	clientService := oauth2.NewClientService(clientRepo)
	clientController := oauth2.NewClientController(clientService)

	authorizationStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationRequest]()
	codeStore := core.NewInMemoryKeyValueStore[oauth2.AuthorizationRequest]()
	authorizationService := oauth2.NewAuthorizationService(authorizationStore, codeStore, clientService)
	authorizationController := oauth2.NewAuthorizationController(authorizationService)
	logger.Info("Initialize services successful")

	r := gin.Default()
	api := r.Group("/v1")
	logger.Debug("Router setup starting")
	clientController.RegisterRoutes(api.Group("/client"))
	authorizationController.RegisterRoutes(api.Group("/oauth2"))
	logger.Info("Router setup successful")

	logger.Info("Application initialization successful")
	logger.Info("Application starting to listen")
	r.Run()
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
