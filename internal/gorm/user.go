package gorm

import (
	"context"

	domain "github.com/Untanky/modern-auth/internal/user"
	"github.com/Untanky/modern-auth/internal/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID     string `gorm:"primaryKey"`
	UserID string `gorm:"type:varchar;unique;index;not null"`
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
					ID:     user.ID,
					UserID: string(utils.EncodeBase64(user.UserID)),
					Status: user.Status,
				}
			},
			toModel: func(gormUser *User) *domain.User {
				userID, _ := utils.DecodeBase64([]byte(gormUser.UserID))
				return &domain.User{
					ID:     gormUser.ID,
					UserID: userID,
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

func (r *GormUserRepo) ExistsUserId(ctx context.Context, userId string) (bool, error) {
	var gormUser User
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).First(&gormUser).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
