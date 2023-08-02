package user

import (
	"context"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
)

type UserRepository interface {
	core.Repository[string, *User]
	FindByUserId(ctx context.Context, userId []byte) (*User, error)
	ExistsUserId(ctx context.Context, userId []byte) (bool, error)
}

type User struct {
	ID     uuid.UUID
	UserID []byte
	Status string
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
	return s.repo.FindByUserId(ctx, s.hashUserId(userId))
}

func (s *UserService) ExistsUserId(ctx context.Context, userId []byte) (bool, error) {
	return s.repo.ExistsUserId(ctx, s.hashUserId(userId))
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	user.ID = uuid.New()
	user.UserID = s.hashUserId(user.UserID)
	user.Status = "active"

	return s.repo.Save(ctx, user)
}

func (s *UserService) DeleteById(ctx context.Context, userId string) error {
	return s.repo.DeleteById(ctx, userId)
}
