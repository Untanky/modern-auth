package core

import (
	"context"

	"gorm.io/gorm"
)

// A repository persists entities of a given type.
// Entities must be unique by their id.
type Repository[IdType comparable, Type interface{}] interface {
	FindAll(ctx context.Context) ([]Type, error)
	FindById(ctx context.Context, id IdType) (Type, error)
	Save(ctx context.Context, entity Type) error
	Update(ctx context.Context, entity Type) error
	DeleteById(ctx context.Context, id IdType) error
}

type TenancyRepository[IdType comparable, Type interface{}] interface {
	Repository[IdType, Type]
	FindAllByTenantId(ctx context.Context, tenantId string) ([]Type, error)
}

type GormTenancyRepository[IdType comparable, Type interface{}] struct {
	DB *gorm.DB
}

func NewGormRepository[IdType comparable, Type interface{}](db *gorm.DB) *GormTenancyRepository[IdType, Type] {
	return &GormTenancyRepository[IdType, Type]{
		DB: db,
	}
}

func (repo *GormTenancyRepository[IdType, Type]) FindAll(ctx context.Context) ([]Type, error) {
	var entities []Type
	err := repo.DB.WithContext(ctx).Find(&entities).Error
	return entities, err
}

func (repo *GormTenancyRepository[IdType, Type]) FindById(ctx context.Context, id IdType) (Type, error) {
	var entity Type
	err := repo.DB.WithContext(ctx).First(&entity, id).Error
	return entity, err
}

func (repo *GormTenancyRepository[IdType, Type]) FindAllByTenantId(ctx context.Context, tenantId string) ([]Type, error) {
	var entities []Type
	err := repo.DB.WithContext(ctx).Where("tenant_id = ?", tenantId).Find(&entities).Error
	return entities, err
}

func (repo *GormTenancyRepository[IdType, Type]) Save(ctx context.Context, entity Type) error {
	err := repo.DB.WithContext(ctx).Create(&entity).Error
	return err
}

func (repo *GormTenancyRepository[IdType, Type]) Update(ctx context.Context, entity Type) error {
	err := repo.DB.WithContext(ctx).Save(&entity).Error
	return err
}

func (repo *GormTenancyRepository[IdType, Type]) DeleteById(ctx context.Context, id IdType) error {
	db := repo.DB.WithContext(ctx)
	entity, err := repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	err = db.Delete(&entity).Error
	return err
}
