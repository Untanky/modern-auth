package oauth2

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/gin-gonic/gin"
)

type ClientModel struct {
	ID           string `gorm:"primaryKey"`
	Scopes       string
	RedirectURIs string
}

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

type ClientRepository = core.Repository[string, *ClientModel]

type ClientDTO struct {
	ID           string   `json:"id"`
	Scopes       []string `json:"scopes"`
	RedirectURIs []string `json:"redirectURIs"`
}

type ClientWithSecretDTO struct {
	ClientDTO
	Secret string `json:"secret"`
}

type ClientService struct {
	repo   ClientRepository
	logger *slog.Logger
}

func NewClientService(repo ClientRepository) *ClientService {
	logger := slog.Default().With(slog.String("service", "client"))

	return &ClientService{repo: repo, logger: logger}
}

func (s *ClientService) FindById(ctx context.Context, id string) (*Client, error) {
	client, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	s.logger.Info("Found client", "client_id", client.ID)

	return &Client{
		ID:           client.ID,
		Scopes:       strings.Split(client.Scopes, ","),
		RedirectURIs: strings.Split(client.RedirectURIs, ","),
	}, nil
}

func (s *ClientService) List(ctx context.Context) ([]*Client, error) {
	clients, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	s.logger.Info("List all clients", "count", len(clients))

	var results []*Client
	for _, client := range clients {
		results = append(results, &Client{
			ID:           client.ID,
			Scopes:       strings.Split(client.Scopes, ","),
			RedirectURIs: strings.Split(client.RedirectURIs, ","),
		})
	}
	return results, nil
}

func (s *ClientService) Create(ctx context.Context, dto ClientDTO) (*Client, error) {
	clientModel := &ClientModel{
		ID:           dto.ID,
		Scopes:       strings.Join(dto.Scopes, ","),
		RedirectURIs: strings.Join(dto.RedirectURIs, ","),
	}
	err := s.repo.Save(ctx, clientModel)
	if err != nil {
		return nil, err
	}
	s.logger.Info("Created client", "client_id", dto.ID)

	client := &Client{
		ID:           dto.ID,
		Scopes:       dto.Scopes,
		RedirectURIs: dto.RedirectURIs,
	}
	return client, nil
}

func (s *ClientService) Update(ctx context.Context, dto ClientDTO) (*Client, error) {
	client, err := s.repo.FindById(ctx, dto.ID)
	if err != nil {
		return nil, err
	}

	client.Scopes = strings.Join(dto.Scopes, ",")
	client.RedirectURIs = strings.Join(dto.RedirectURIs, ",")

	err = s.repo.Update(ctx, client)
	if err != nil {
		return nil, err
	}
	s.logger.Info("Updated client", "client_id", dto.ID)

	return &Client{
		ID:           dto.ID,
		Scopes:       dto.Scopes,
		RedirectURIs: dto.RedirectURIs,
	}, nil
}

func (s *ClientService) Delete(ctx context.Context, id string) error {
	err := s.repo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	s.logger.Info("Deleted client", "client_id", id)
	return nil
}

type ClientController struct {
	service *ClientService
}

func NewClientController(service *ClientService) *ClientController {
	return &ClientController{service: service}
}

func (c *ClientController) RegisterRoutes(router gin.IRouter) {
	router.GET("", c.list)
	router.GET("/:id", c.get)
	router.POST("", c.create)
	router.DELETE("/:id", c.delete)
}

func (c *ClientController) list(ctx *gin.Context) {
	clients, err := c.service.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var dtos []ClientDTO = make([]ClientDTO, 0, len(clients))
	for _, client := range clients {
		dtos = append(dtos, ClientDTO{
			ID:           client.ID,
			Scopes:       client.Scopes,
			RedirectURIs: client.RedirectURIs,
		})
	}
	ctx.JSON(http.StatusOK, dtos)
}

func (c *ClientController) get(ctx *gin.Context) {
	id := ctx.Param("id")
	client, err := c.service.FindById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, &ClientDTO{
		ID:           client.ID,
		Scopes:       client.Scopes,
		RedirectURIs: client.RedirectURIs,
	})
}

func (c *ClientController) create(ctx *gin.Context) {
	var dto ClientDTO
	err := ctx.BindJSON(&dto)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := c.service.Create(ctx, dto)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, &ClientDTO{
		ID:           client.ID,
		Scopes:       client.Scopes,
		RedirectURIs: client.RedirectURIs,
	})
}

func (c *ClientController) delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}
