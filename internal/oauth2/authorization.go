package oauth2

import (
	"context"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthorizationRequest struct {
	ClientId      string
	RedirectUri   string
	ResponseType  string
	Scope         string
	State         string
	CodeChallenge string
	CodeMethod    string
}

type AuthorizationResponse struct {
	RedirectUri string
	State       string
	Code        string
}

type AuthorizationError struct {
	RedirectUri string
	State       string
	Error       string
	Description string
}

type AuthorizationStore = core.KeyValueStore[string, AuthorizationRequest]

type AuthorizationService struct {
	authorizationStore AuthorizationStore
	clientService      *ClientService
}

func NewAuthorizationService(authorizationStore AuthorizationStore, clientService *ClientService) *AuthorizationService {
	return &AuthorizationService{authorizationStore: authorizationStore, clientService: clientService}
}

// Note: maybe the function argument should be a dto instead of a request...
func (s *AuthorizationService) Authorize(request *AuthorizationRequest) (string, *AuthorizationError) {
	client, err := s.clientService.FindById(context.TODO(), request.ClientId)
	if err != nil {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "invalid_client",
			Description: "client not found",
		}
	}
	if request.RedirectUri == "" {
		request.RedirectUri = client.RedirectURIs[0]
	} else if !client.ValidateRedirectURI(context.TODO(), request.RedirectUri) {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "invalid_request",
			Description: "redirect_uri not allowed",
		}
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "server_error",
		}
	}

	stringUuid := uuid.String()
	err = s.authorizationStore.Set(stringUuid, request)
	if err != nil {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "server_error",
		}
	}

	return stringUuid, nil
}

type AuthorizationController struct {
	authorizationService *AuthorizationService
}

func NewAuthorizationController(authorizationService *AuthorizationService) *AuthorizationController {
	return &AuthorizationController{authorizationService: authorizationService}
}

func (c *AuthorizationController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/authorization", c.authorize)
}

func (c *AuthorizationController) authorize(ctx *gin.Context) {
	request := &AuthorizationRequest{
		ClientId:      ctx.Query("client_id"),
		CodeChallenge: ctx.Query("code_challenge"),
		CodeMethod:    ctx.Query("code_method"),
		RedirectUri:   ctx.Query("redirect_uri"),
		ResponseType:  ctx.Query("response_type"),
		Scope:         ctx.Query("scope"),
		State:         ctx.Query("state"),
	}

	uuid, err := c.authorizationService.Authorize(request)
	if err != nil {
		ctx.Redirect(302, err.RedirectUri)
	}
	ctx.SetCookie("authorization", uuid, 0, "/", "", false, true)
	ctx.Redirect(302, "/index.html")
}
