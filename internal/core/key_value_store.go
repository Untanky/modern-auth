package core

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func init() {
	tp := otel.GetTracerProvider()
	tracer = tp.Tracer("KeyValueStore")
}

// A Cache stores values of a given type under a given key.
type KeyValueStore[Key comparable, Type interface{}] interface {
	// Get returns the value associated with the given key.
	Get(key Key) (*Type, error)
	// Set associates the given value with the given value.
	Set(key Key, value *Type) error
	// Delete removes the value associated with the given key.
	Delete(key Key) error

	WithContext(ctx context.Context) KeyValueStore[Key, Type]
}

type contextKeyValueStore[Key comparable, Type interface{}] struct {
	ctx   context.Context
	store KeyValueStore[Key, Type]
}

func (store *contextKeyValueStore[Key, Type]) Get(key Key) (*Type, error) {
	_, span := tracer.Start(store.ctx, "keyValueStore.Get")
	defer span.End()
	return store.store.Get(key)
}

func (store *contextKeyValueStore[Key, Type]) Set(key Key, value *Type) error {
	_, span := tracer.Start(store.ctx, "keyValueStore.Set")
	defer span.End()
	return store.store.Set(key, value)
}

func (store *contextKeyValueStore[Key, Type]) Delete(key Key) error {
	_, span := tracer.Start(store.ctx, "keyValueStore.Delete")
	defer span.End()
	return store.store.Delete(key)
}

func (store *contextKeyValueStore[Key, Type]) WithContext(ctx context.Context) KeyValueStore[Key, Type] {
	return &contextKeyValueStore[Key, Type]{
		ctx:   ctx,
		store: store,
	}
}

type InMemoryKeyValueStore[Type interface{}] struct {
	storage map[string]Type
}

func NewInMemoryKeyValueStore[Type interface{}]() *InMemoryKeyValueStore[Type] {
	return &InMemoryKeyValueStore[Type]{
		storage: make(map[string]Type),
	}
}

func (store *InMemoryKeyValueStore[Type]) Get(key string) (*Type, error) {
	value, ok := store.storage[key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}
	return &value, nil
}

func (store *InMemoryKeyValueStore[Type]) Set(key string, value *Type) error {
	store.storage[key] = *value
	return nil
}

func (store *InMemoryKeyValueStore[Type]) Delete(key string) error {
	delete(store.storage, key)
	return nil
}

func (store *InMemoryKeyValueStore[Type]) WithContext(ctx context.Context) KeyValueStore[string, Type] {
	return &contextKeyValueStore[string, Type]{
		ctx:   ctx,
		store: store,
	}
}
