package oauth2

import (
	"context"
	"fmt"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthorizationRequest struct {
	id            uuid.UUID
	ClientId      string
	RedirectUri   string
	ResponseType  string
	Scope         string
	State         string
	CodeChallenge string
	CodeMethod    string
}

type ResponseUriBuilder interface {
	BuildResponseURI() string
}

type AuthorizationResponse struct {
	RedirectUri string
	State       string
	Issuer      string
	Code        string
}

func (r *AuthorizationResponse) BuildResponseURI() string {
	return fmt.Sprintf("%s?code=%s&state=%s&iss=%s", r.RedirectUri, r.Code, r.State, r.Issuer)
}

type AuthorizationError struct {
	RedirectUri string
	State       string
	Issuer      string
	Error       string
	Description string
}

func (e *AuthorizationError) BuildResponseURI() string {
	return fmt.Sprintf("%s?error=%s&error_description=%s&state=%s&iss=%s", e.RedirectUri, e.Error, e.Description, e.State, e.Issuer)
}

type AuthorizationStore = core.KeyValueStore[string, *AuthorizationRequest]
type CodeStore = core.KeyValueStore[string, *AuthorizationRequest]

type AuthorizationService struct {
	authorizationStore AuthorizationStore
	codeStore          CodeStore
	clientService      *ClientService
	logger             *zap.SugaredLogger
}

func NewAuthorizationService(authorizationStore AuthorizationStore, codeStore CodeStore, clientService *ClientService, logger *zap.SugaredLogger) *AuthorizationService {
	return &AuthorizationService{authorizationStore: authorizationStore, codeStore: codeStore, clientService: clientService, logger: logger}
}

func (s *AuthorizationService) Authorize(request *AuthorizationRequest) (string, *AuthorizationError) {
	s.logger.Debug("Beginning 'authorization_code' flow")
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
	request.id = uuid
	err = s.authorizationStore.Set(stringUuid, request)
	if err != nil {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "server_error",
		}
	}
	s.logger.Infow("Initialized 'authorization_code' flow", "authorizationId", stringUuid)

	return stringUuid, nil
}

func (s *AuthorizationService) Succeed(uuid string) ResponseUriBuilder {
	s.logger.Debugw("Continuing 'authorization_code' flow", "authorizationId", uuid)
	request, err := s.authorizationStore.Get(uuid)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "server_error",
		}
	}

	err = s.authorizationStore.Delete(uuid)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "server_error",
		}
	}

	code := randomString(32) // TODO: generate code
	err = s.codeStore.Set(code, request)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			Error:       "server_error",
		}
	}
	s.logger.Infow("Generated authorization code", "authroizationId", uuid)

	return &AuthorizationResponse{
		RedirectUri: request.RedirectUri,
		Code:        code,
		State:       request.State,
		Issuer:      "https://localhost:8080",
	}
}

type AuthorizationController struct {
	authorizationService *AuthorizationService
}

func NewAuthorizationController(authorizationService *AuthorizationService) *AuthorizationController {
	return &AuthorizationController{authorizationService: authorizationService}
}

func (c *AuthorizationController) RegisterRoutes(router gin.IRouter) {
	router.GET("/authorization", c.authorize)
	router.GET("/authorization/succeed", c.succeed)
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
		ctx.Redirect(302, err.BuildResponseURI())
	}
	ctx.SetCookie("authorization", uuid, 0, "/", "", false, true)
	ctx.Redirect(302, "/index.html")
}

func (c *AuthorizationController) succeed(ctx *gin.Context) {
	uuid, err := ctx.Cookie("authorization")
	if err != nil {
		ctx.Redirect(302, "/index.html")
		return
	}
	response := c.authorizationService.Succeed(uuid)
	ctx.SetCookie("authorization", "", -1, "/", "", false, true)
	ctx.Redirect(302, response.BuildResponseURI())
}
