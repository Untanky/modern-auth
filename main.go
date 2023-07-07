package main

import (
	"log"

	"github.com/Untanky/modern-auth/internal/core"
)

type Person struct {
	Name string
}

func (p Person) String() string {
	return p.Name
}

func main() {
	store := core.NewInMemoryKeyValueStore[Person]()

	err := store.Set("1", &Person{Name: "John"})
	log.Println("Set:", err)

	value, err := store.Get("1")
	log.Println("Get:", value, err)

	err = store.Set("1", &Person{Name: "Peter"})
	log.Println("Set:", err)

	value, err = store.Get("1")
	log.Println("Get:", value, err)

	err = store.Delete("1")
	log.Println("Delete:", err)
}
