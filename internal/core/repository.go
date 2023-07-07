package core

// A repository persists entities of a given type.
// Entities must be unique by their id.
type Repository[IdType comparable, Type interface{}] interface {
	FindAll() ([]Type, error)
	FindById(id IdType) (Type, error)
	Save(entity Type) error
	Update(entity Type) error
	DeleteById(id IdType) error
}

type TenancyRepository[IdType comparable, Type interface{}] interface {
	Repository[IdType, Type]
	FindAllByTenantId(tenantId string) ([]Type, error)
}
