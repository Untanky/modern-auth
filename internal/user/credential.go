package user

import (
	"context"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/google/uuid"
)

type CredentialRepository interface {
	core.Repository[string, *Credential]
	FindByCredentialId(ctx context.Context, credentialId []byte) (*Credential, error)
}

type Credential struct {
	ID           uuid.UUID
	CredentialID []byte
	Status       string
	PublicKey    []byte
	User         User
}

type CredentialService struct {
	repo CredentialRepository
}

func NewCredentialService(credentialRepo CredentialRepository) *CredentialService {
	return &CredentialService{
		repo: credentialRepo,
	}
}

func (s *CredentialService) GetCredentialByID(ctx context.Context, id string) (*Credential, error) {
	return s.repo.FindById(ctx, id)
}

func (s *CredentialService) GetCredentialByCredentialID(ctx context.Context, id []byte) (*Credential, error) {
	return s.repo.FindByCredentialId(ctx, id)
}

func (s *CredentialService) CreateCredential(ctx context.Context, credential *Credential) error {
	credential.ID = uuid.New()
	credential.Status = "active"

	return s.repo.Save(ctx, credential)
}

func (s *CredentialService) DeleteById(ctx context.Context, credentialId string) error {
	return s.repo.DeleteById(ctx, credentialId)
}
