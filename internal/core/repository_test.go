package core_test

import (
	"context"
	"os"
	"testing"

	"github.com/Untanky/modern-auth/internal/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormPerson struct {
	gorm.Model
	ID   string
	Name string
}

func (p GormPerson) String() string {
	return p.Name
}

func TestRepositoryFullFlow(t *testing.T) {
	type args struct {
		dbFilename string
	}
	tests := []struct {
		name    string
		args    args
		init    func(args) core.Repository[string, GormPerson]
		cleanUp func(args)
	}{
		{
			name: "In Memory Key Value Store",
			args: args{
				dbFilename: "gorm.db",
			},
			init: func(args args) core.Repository[string, GormPerson] {
				db, err := gorm.Open(sqlite.Open(args.dbFilename), &gorm.Config{})
				if err != nil {
					panic("failed to connect database")
				}

				// Migrate the schema
				db.AutoMigrate(&GormPerson{})
				return core.NewGormRepository[string, GormPerson](db)
			},
			cleanUp: func(args args) {
				os.Remove(args.dbFilename)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.init(tt.args)

			err := repository.Save(context.TODO(), GormPerson{ID: "1", Name: "John"})
			if err != nil {
				t.Errorf("Save() error = %v", err)
				return
			}

			list, err := repository.FindAll(context.TODO())
			if err != nil {
				t.Errorf("FindAll() error = %v", err)
				return
			}
			if len(list) != 1 {
				t.Errorf("FindAll() got = %v, want %v", len(list), 1)
				return
			}

			value, err := repository.FindById(context.TODO(), "1")
			if err != nil {
				t.Errorf("FindById() error = %v", err)
				return
			}
			if value.Name != "John" {
				t.Errorf("FindById() got = %v, want %v", value.Name, "John")
				return
			}

			value.Name = "Peter"
			err = repository.Update(context.TODO(), value)
			if err != nil {
				t.Errorf("Update() error = %v", err)
				return
			}

			value, err = repository.FindById(context.TODO(), "1")
			if err != nil {
				t.Errorf("FindById() error = %v", err)
				return
			}
			if value.Name != "Peter" {
				t.Errorf("FindById() got = %v, want %v", value.Name, "Peter")
				return
			}

			err = repository.DeleteById(context.TODO(), "1")
			if err != nil {
				t.Errorf("DeleteById() error = %v", err)
				return
			}

			list, err = repository.FindAll(context.TODO())
			if err != nil {
				t.Errorf("FindAll() error = %v", err)
				return
			}
			if len(list) != 0 {
				t.Errorf("FindAll() got = %v, want %v", len(list), 0)
				return
			}

			tt.cleanUp(tt.args)
		})
	}
}
