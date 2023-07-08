package client

import (
	"context"

	"github.com/Untanky/modern-auth/internal/core"
)

type Client struct {
	ID           string
	Scopes       []string
	RedirectURIs []string
}

func (c *Client) RestrictScopes(ctx context.Context, scopes []string) []string {
	var restrictedScopes []string
	for _, scope := range scopes {
		for _, allowedScope := range c.Scopes {
			if scope == allowedScope {
				restrictedScopes = append(restrictedScopes, scope)
			}
		}
	}
	return restrictedScopes
}

func (c *Client) ValidateRedirectURI(ctx context.Context, redirectURI string) bool {
	for _, uri := range c.RedirectURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}

type ClientRepository = core.Repository[string, *Client]

type ClientDTO struct {
	ID           string   `json:"id"`
	Scopes       []string `json:"scopes"`
	RedirectURIs []string `json:"redirect_uris"`
}

type ClientWithSecretDTO struct {
	ClientDTO
	Secret string `json:"secret"`
}

type ClientService struct {
	repo ClientRepository
}

func NewClientService(repo ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) Create(ctx context.Context, dto ClientDTO) (*Client, error) {
	client := &Client{
		ID:           dto.ID,
		Scopes:       dto.Scopes,
		RedirectURIs: dto.RedirectURIs,
	}
	err := s.repo.Save(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (s *ClientService) Update(ctx context.Context, dto ClientDTO) (*Client, error) {
	client, err := s.repo.FindById(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	client.Scopes = dto.Scopes
	client.RedirectURIs = dto.RedirectURIs

	err = s.repo.Save(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}
