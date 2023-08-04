package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/google/uuid"
)

type GrantRepo = core.Repository[uuid.UUID, *grant]
type GrantStore = core.KeyValueStore[string, *grant]

var grantRepo GrantRepo
var grantStore GrantStore

func init() {
	grantStore = core.NewInMemoryKeyValueStore[*grant]()
}

type grant struct {
	ID                uuid.UUID
	SubjectID         uuid.UUID
	ClientID          string
	Scope             []string
	IssuedAt          time.Time
	ExpiresAt         time.Time
	NotBefore         time.Time
	AllowRefreshToken bool
}

func NewGrant(subjectID uuid.UUID) *grant {
	return &grant{
		ID:       uuid.New(),
		IssuedAt: time.Now(),
	}
}

type Token interface {
	Key() string
	MarshalJSON() ([]byte, error)
}

type AccessToken [32]byte

func (token *AccessToken) Key() string {
	return fmt.Sprintf("access_token:%s", utils.EncodeBase64(utils.HashShake256(token[:])))
}

func (token *AccessToken) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", utils.EncodeBase64(token[:]))), nil
}

type RefreshToken [48]byte

func (token RefreshToken) Key() string {
	return fmt.Sprintf("refresh_token:%s", utils.EncodeBase64(utils.HashShake256(token[:])))
}

func (token *RefreshToken) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", utils.EncodeBase64(token[:]))), nil
}

func FindGrantByGrantID(grantID uuid.UUID) (*grant, error) {
	g, err := grantStore.Get(grantID.String())
	if err != nil {
		return nil, err
	}
	return g, nil
}

func FindGrantByToken(ctx context.Context, accessToken Token) (*grant, error) {
	g, err := grantStore.Get(accessToken.Key())
	if err != nil {
		return nil, err
	}
	return g, nil
}

func RegisterGrant(ctx context.Context, g *grant) (*AccessToken, *RefreshToken, error) {
	// err := grantRepo.Save(ctx, g)
	// if err != nil {
	// 	return nil, nil, err
	// }

	accessToken := createAccessToken()
	err := storeToken(accessToken, g)
	if err != nil {
		return nil, nil, err
	}

	if g.AllowRefreshToken {
		refreshToken := createRefreshToken()
		err = storeToken(refreshToken, g)
		if err != nil {
			return nil, nil, err
		}
		return accessToken, refreshToken, nil
	}

	return accessToken, nil, nil
}

func createAccessToken() *AccessToken {
	tokenBytes := make([]byte, 32)
	utils.RandomBytes(tokenBytes)
	accessToken := AccessToken(tokenBytes)
	return &accessToken
}

func createRefreshToken() *RefreshToken {
	tokenBytes := make([]byte, 48)
	utils.RandomBytes(tokenBytes)
	refreshToken := RefreshToken(tokenBytes)
	return &refreshToken
}

func storeToken(token Token, g *grant) error {
	return grantStore.Set(token.Key(), g)
}

func LeaseGrant(ctx context.Context, refreshToken *RefreshToken) (*AccessToken, *grant, error) {
	g, err := FindGrantByToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	accessToken := createAccessToken()
	err = storeToken(accessToken, g)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, g, nil
}
