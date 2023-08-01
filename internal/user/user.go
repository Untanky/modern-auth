package user

import (
	"context"
	"log"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/google/uuid"
)

type UserRepository interface {
	core.Repository[string, *User]
	FindByUserId(ctx context.Context, userId string) (*User, error)
	ExistsUserId(ctx context.Context, userId string) (bool, error)
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
	hashedUserID := utils.HashShake256([]byte(userId))
	return s.repo.FindByUserId(ctx, string(hashedUserID))
}

func (s *UserService) ExistsUserId(ctx context.Context, userId string) (bool, error) {
	hashedUserID := utils.HashShake256([]byte(userId))
	return s.repo.ExistsUserId(ctx, string(hashedUserID))
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	log.Println("Creating user")
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	user.ID = uuid.String()
	user.UserID = utils.HashShake256(user.UserID)
	user.Status = "active"

	return s.repo.Save(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user *User) error {
	return s.repo.Update(ctx, user)
}

func (s *UserService) DeleteById(ctx context.Context, userId string) error {
	return s.repo.DeleteById(ctx, userId)
}
