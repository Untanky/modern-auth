package user

import (
	"context"

	"github.com/Untanky/modern-auth/internal/core"
)

type UserRepository interface {
	core.Repository[string, *User]
	FindByUserId(ctx context.Context, userId string) (*User, error)
}

type User struct {
	ID     string `gorm:"primary_key"`
	UserID []byte `gorm:"unique"`
	Status string `gorm:"not null"`
}

type UserService struct {
	repo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		repo: userRepo,
	}
}

func (s *UserService) GetUserById(ctx context.Context, id string) (*User, error) {
	return s.repo.FindById(ctx, id)
}

func (s *UserService) GetUserByUserId(ctx context.Context, userId string) (*User, error) {
	return s.repo.FindByUserId(ctx, userId)
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	return s.repo.Save(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteById(ctx context.Context, userId string) error {
	return s.repo.DeleteById(ctx, userId)
}
