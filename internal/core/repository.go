package core

import (
	"context"
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
