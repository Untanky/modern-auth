package gorm

import (
	"context"

	domain "github.com/Untanky/modern-auth/internal/user"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID []byte `gorm:"unique"`
	Status string `gorm:"not null"`
}

type GormUserRepo struct {
	GormRepository[string, *User, *domain.User]
}

func NewGormUserRepo(db *gorm.DB) *GormUserRepo {
	return &GormUserRepo{
		GormRepository: GormRepository[string, *User, *domain.User]{
			db: db,
			toGormModel: func(user *domain.User) *User {
				return &User{
					UserID: user.UserID,
					Status: user.Status,
				}
			},
			toModel: func(gormUser *User) *domain.User {
				return &domain.User{
					UserID: gormUser.UserID,
					Status: gormUser.Status,
				}
			},
		},
	}
}

func (r *GormUserRepo) FindByUserId(ctx context.Context, userId string) (*domain.User, error) {
	var gormUser User
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&gormUser).Error
	if err != nil {
		return nil, err
	}

	return r.toModel(&gormUser), nil
}
