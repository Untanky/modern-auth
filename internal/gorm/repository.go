package gorm

import (
	"context"

	"gorm.io/gorm"
)

type GormRepository[IdType comparable, GormType interface{}, ModelType interface{}] struct {
	db          *gorm.DB
	toGormModel func(entity ModelType) GormType
	toModel     func(entity GormType) ModelType
}

func NewGormRepository[IdType comparable, GormType interface{}, ModelType interface{}](db *gorm.DB, toGormModel func(entity ModelType) GormType, toModel func(entity GormType) ModelType) *GormRepository[IdType, GormType, ModelType] {
	return &GormRepository[IdType, GormType, ModelType]{
		db:          db,
		toGormModel: toGormModel,
		toModel:     toModel,
	}
}

func (repo *GormRepository[IdType, GormType, ModelType]) FindAll(ctx context.Context) ([]ModelType, error) {
	var entities []GormType
	err := repo.db.WithContext(ctx).Find(&entities).Error

	var models []ModelType
	for _, entity := range entities {
		models = append(models, repo.toModel(entity))
	}

	return models, err
}

func (repo *GormRepository[IdType, GormType, ModelType]) FindById(ctx context.Context, id IdType) (ModelType, error) {
	var entity GormType
	err := repo.db.WithContext(ctx).First(&entity, "id = ?", id).Error
	return repo.toModel(entity), err
}

func (repo *GormRepository[IdType, GormType, ModelType]) Save(ctx context.Context, entity ModelType) error {
	err := repo.db.WithContext(ctx).Create(repo.toGormModel(entity)).Error
	return err
}

func (repo *GormRepository[IdType, GormType, ModelType]) Update(ctx context.Context, entity ModelType) error {
	err := repo.db.WithContext(ctx).Save(repo.toGormModel(entity)).Error
	return err
}

func (repo *GormRepository[IdType, GormType, ModelType]) DeleteById(ctx context.Context, id IdType) error {
	db := repo.db.WithContext(ctx)
	entity, err := repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	err = db.Delete(&entity).Error
	return err
}
