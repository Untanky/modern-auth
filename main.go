package main

import (
	"context"
	"log"

	"github.com/Untanky/modern-auth/internal/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Person struct {
	Name string
}

func (p Person) String() string {
	return p.Name
}

func main() {
	app := core.NewApp(core.NewParentModule(core.NewDatabaseModule(nil)))
	app.Start()
	app.Stop()

	// runInMemoryKeyValueStore()
	// runGormRepository()
}

func runInMemoryKeyValueStore() {
	log.Println("In Memory Key Value Store Example")
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

type Product struct {
	gorm.Model
	ID    string
	Code  string
	Price uint
}

func runGormRepository() {
	log.Println("Gorm Repository Example")
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Product{})

	repository := core.NewGormRepository[string, Product](db)

	err = repository.Save(context.TODO(), Product{ID: "1", Code: "D42", Price: 100})
	log.Println("Save:", err)

	list, err := repository.FindAll(context.TODO())
	log.Println("Find:", list, err)

	value, err := repository.FindById(context.TODO(), "1")
	log.Println("FindById:", value, err)

	value.Price = 800
	err = repository.Update(context.TODO(), value)
	log.Println("Update:", err)

	err = repository.DeleteById(context.TODO(), "1")
	log.Println("Delete:", err)

	list, err = repository.FindAll(context.TODO())
	log.Println("Find:", list, err)
}
