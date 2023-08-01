package user

import (
	"context"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

type UserRepository interface {
	core.Repository[string, *User]
	FindByUserId(ctx context.Context, userId []byte) (*User, error)
	ExistsUserId(ctx context.Context, userId []byte) (bool, error)
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

func (s *UserService) hashUserId(userId []byte) []byte {
	hashedUserID := make([]byte, 32)
	sha3.ShakeSum256(hashedUserID, userId)
	return hashedUserID
}

func (s *UserService) GetUserByUserID(ctx context.Context, userId []byte) (*User, error) {
	return s.repo.FindByUserId(ctx, utils.EncodeBase64(s.hashUserId(userId)))
}

func (s *UserService) ExistsUserId(ctx context.Context, userId []byte) (bool, error) {
	return s.repo.ExistsUserId(ctx, utils.EncodeBase64(s.hashUserId(userId)))
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	user.ID = uuid.String()
	user.UserID = utils.EncodeBase64(s.hashUserId(user.UserID))
	user.Status = "active"

	return s.repo.Save(ctx, user)
}

func (s *UserService) DeleteById(ctx context.Context, userId string) error {
	return s.repo.DeleteById(ctx, userId)
}
