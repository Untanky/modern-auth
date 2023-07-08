package oauth2

import (
	"context"
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
	repo ClientRepository
}

func NewClientService(repo ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) FindById(ctx context.Context, id string) (*Client, error) {
	client, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &Client{
		ID:           client.ID,
		Scopes:       strings.Split(client.Scopes, ","),
		RedirectURIs: strings.Split(client.RedirectURIs, ","),
	}, nil

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

	err = s.repo.Save(ctx, client)
	if err != nil {
		return nil, err
	}

	return &Client{
		ID:           dto.ID,
		Scopes:       dto.Scopes,
		RedirectURIs: dto.RedirectURIs,
	}, nil
}

type ClientController struct {
	service *ClientService
	repo    ClientRepository
}

func NewClientController(service *ClientService) *ClientController {
	return &ClientController{service: service, repo: service.repo}
}

func (c *ClientController) RegisterRoutes(router gin.IRouter) {
	router.GET("", c.list)
	router.GET("/:id", c.get)
	router.POST("", c.create)
	router.DELETE("/:id", c.delete)
}

func (c *ClientController) list(ctx *gin.Context) {
	clients, err := c.repo.FindAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var dtos []ClientDTO = make([]ClientDTO, 0, len(clients))
	for _, client := range clients {
		dtos = append(dtos, ClientDTO{
			ID:           client.ID,
			Scopes:       strings.Split(client.Scopes, ","),
			RedirectURIs: strings.Split(client.RedirectURIs, ","),
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
	err := c.repo.DeleteById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusOK)
}
