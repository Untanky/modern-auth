package core

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Module interface {
	Init(context.Context) context.Context
	GetEntities() []interface{}
	Shutdown()
}

type ParentModule struct {
	modules []Module
}

func NewParentModule(modules ...Module) *ParentModule {
	return &ParentModule{modules: modules}
}

func (pm *ParentModule) Init(ctx context.Context) context.Context {
	for _, module := range pm.modules {
		ctx = module.Init(ctx)
	}
	return ctx
}

func (pm *ParentModule) GetEntities() []interface{} {
	var entities []interface{}
	for _, module := range pm.modules {
		entities = append(entities, module.GetEntities()...)
	}
	return entities
}

func (pm *ParentModule) Shutdown() {
	for _, module := range pm.modules {
		module.Shutdown()
	}
}

type DatabaseModule struct {
	db     *gorm.DB
	config any
}

func NewDatabaseModule(config any) *DatabaseModule {
	return &DatabaseModule{config: config}
}

type databaseKey string

var DatabaseKey = databaseKey("database")

func (dm *DatabaseModule) Init(ctx context.Context) context.Context {
	dm.db = dm.connect()

	entities := ctx.Value(EntitiesKey).([]interface{})
	dm.migrateEntities(entities)

	return context.WithValue(ctx, DatabaseKey, dm.db)
}

func (dm *DatabaseModule) connect() *gorm.DB {
	logger.Debug("Database connection starting")
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	logger.Info("Database connection successful")
	return db
}

func (dm *DatabaseModule) migrateEntities(entities []interface{}) {
	logger.Debug("Entity migration starting")
	for _, entity := range entities {
		err := dm.db.AutoMigrate(entity)
		if err != nil {
			panic("failed to migrate entity")
		}
	}
	logger.Info("Entity migration successful")
}

func (dm *DatabaseModule) GetEntities() []interface{} {
	return []interface{}{}
}

func (dm *DatabaseModule) Shutdown() {
	db, _ := dm.db.DB()
	if db != nil {
		db.Close()
	}
}
