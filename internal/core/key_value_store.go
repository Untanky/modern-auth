package core

import "fmt"

// A Cache stores values of a given type under a given key.
type KeyValueStore[Key comparable, Type interface{}] interface {
	// Get returns the value associated with the given key.
	Get(key Key) (*Type, error)
	// Set associates the given value with the given value.
	Set(key Key, value *Type) error
	// Delete removes the value associated with the given key.
	Delete(key Key) error
}

type InMemoryKeyValueStore[Type fmt.Stringer] struct {
	storage map[string]Type
}

func NewInMemoryKeyValueStore[Type fmt.Stringer]() *InMemoryKeyValueStore[Type] {
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
