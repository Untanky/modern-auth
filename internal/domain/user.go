package domain

import (
	"context"
	"log/slog"

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
	ID     uuid.UUID
	UserID []byte
	Status string
}

type UserService struct {
	repo   UserRepository
	logger *slog.Logger
}

func NewUserService(userRepo UserRepository) *UserService {
	logger := slog.Default().With(slog.String("service", "user"))

	return &UserService{
		repo:   userRepo,
		logger: logger,
	}
}

func (s *UserService) GetUserById(ctx context.Context, id string) (*User, error) {
	s.logger.InfoContext(ctx, "Finding user", slog.String("id", id))
	return s.repo.FindById(ctx, id)
}

func (s *UserService) hashUserId(userId []byte) []byte {
	hashedUserID := make([]byte, 32)
	sha3.ShakeSum256(hashedUserID, userId)
	return hashedUserID
}

func (s *UserService) GetUserByUserID(ctx context.Context, userId []byte) (*User, error) {
	hashedUserId := s.hashUserId(userId)
	s.logger.InfoContext(ctx, "Finding user", "userId", utils.EncodeBase64(hashedUserId))
	return s.repo.FindByUserId(ctx, hashedUserId)
}

func (s *UserService) CreateUser(ctx context.Context, user *User) error {
	user.ID = uuid.New()
	user.UserID = s.hashUserId(user.UserID)
	user.Status = "active"

	s.logger.InfoContext(ctx, "Creating user", "id", user.ID, "userId", utils.EncodeBase64(user.UserID))

	return s.repo.Save(ctx, user)
}

func (s *UserService) DeleteById(ctx context.Context, id string) error {
	s.logger.InfoContext(ctx, "Deleting user", "id", id)

	return s.repo.DeleteById(ctx, id)
}
