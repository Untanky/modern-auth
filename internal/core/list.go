package core

import "context"

type List[Type interface{}] interface {
	Index(Type) (int64, error)
	Len() int64
	Append(Type) (int64, error)
	Remove(int64) error
	WithContext(ctx context.Context) List[Type]
	Slice() []Type
}
