package core_test

import (
	"testing"

	"github.com/Untanky/modern-auth/internal/core"
)

type Person struct {
	Name string
}

func (p Person) String() string {
	return p.Name
}

func KeyValueStoreFullFlow(t *testing.T) {
	tests := []struct {
		name  string
		store core.KeyValueStore[string, Person]
	}{
		{
			name:  "In Memory Key Value Store",
			store: core.NewInMemoryKeyValueStore[Person](),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.store.Set("1", &Person{Name: "John"})
			if err != nil {
				t.Errorf("Set() error = %v", err)
			}

			value, err := tt.store.Get("1")
			if err != nil {
				t.Errorf("Get() error = %v", err)
			}
			if value == nil || value.Name != "John" {
				t.Errorf("Get() value = %v, want %v", value, "John")
			}

			err = tt.store.Set("1", &Person{Name: "Peter"})
			if err != nil {
				t.Errorf("Set() error = %v", err)
			}

			value, err = tt.store.Get("1")
			if err != nil {
				t.Errorf("Get() error = %v", err)
			}
			if value == nil || value.Name != "Peter" {
				t.Errorf("Get() value = %v, want %v", value, "Peter")
			}

			err = tt.store.Delete("1")
			if err != nil {
				t.Errorf("Delete() error = %v", err)
			}
		})
	}
}
