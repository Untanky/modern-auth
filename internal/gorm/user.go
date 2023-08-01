package gorm

import (
	"context"

	domain "github.com/Untanky/modern-auth/internal/user"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID     uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID []byte    `gorm:"type:bytea;unique;index;not null"`
	Status string    `gorm:"not null"`
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
					ID:     user.ID,
					UserID: user.UserID,
					Status: user.Status,
				}
			},
			toModel: func(gormUser *User) *domain.User {
				return &domain.User{
					ID:     gormUser.ID,
					UserID: gormUser.UserID,
					Status: gormUser.Status,
				}
			},
		},
	}
}

func (r *GormUserRepo) FindByUserId(ctx context.Context, userId []byte) (*domain.User, error) {
	var gormUser User
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&gormUser).Error
	if err != nil {
		return nil, err
	}

	return r.toModel(&gormUser), nil
}

func (r *GormUserRepo) ExistsUserId(ctx context.Context, userId []byte) (bool, error) {
	var gormUser User
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&gormUser).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
