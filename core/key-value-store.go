package core

// A Cache stores values of a given type under a given key.
type Cache[Key comparable, Type interface{}] interface {
	// Get returns the value associated with the given key.
	Get(key Key) (Type, error)
	// Set associates the given value with the given value.
	Set(key Key, value Type) error
	// Delete removes the value associated with the given key.
	Delete(key Key) error
}
