package gorm

import (
	"context"

	domain "github.com/Untanky/modern-auth/internal/user"
	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Credential struct {
	gorm.Model
	ID           uuid.UUID `gorm:"primaryKey;type:uuid"`
	CredentialID []byte    `gorm:"type:bytea;unique;index;not null"`
	PublicKey    []byte    `gorm:"type:bytea;not null"`
	Status       string    `gorm:"not null"`
}

type GormCredentialRepo struct {
	GormRepository[string, *Credential, *domain.Credential]
}

func NewGormCredentialRepo(db *gorm.DB) *GormCredentialRepo {
	return &GormCredentialRepo{
		GormRepository: GormRepository[string, *Credential, *domain.Credential]{
			db: db,
			toGormModel: func(credential *domain.Credential) *Credential {
				return &Credential{
					ID:           credential.ID,
					CredentialID: credential.CredentialID,
					PublicKey:    credential.PublicKey,
					Status:       credential.Status,
				}
			},
			toModel: func(gormCredential *Credential) *domain.Credential {
				return &domain.Credential{
					ID:           gormCredential.ID,
					CredentialID: gormCredential.CredentialID,
					PublicKey:    gormCredential.PublicKey,
					Status:       gormCredential.Status,
				}
			},
		},
	}
}

func (r *GormCredentialRepo) FindByCredentialId(ctx context.Context, credentialId []byte) (*domain.Credential, error) {
	var gormCredential Credential
	credentialId = utils.EncodeBase64(credentialId)
	err := r.db.WithContext(ctx).Where("credential_id = ?", credentialId).First(&gormCredential).Error
	if err != nil {
		return nil, err
	}

	return r.toModel(&gormCredential), nil
}
