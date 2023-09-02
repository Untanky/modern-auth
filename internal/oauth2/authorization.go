package oauth2

import (
	"context"
	"fmt"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthorizationRequest struct {
	id                 uuid.UUID
	authenticationCode []byte
	ClientId           string
	RedirectUri        string
	ResponseType       string
	Scope              string
	State              string
	CodeChallenge      string
	CodeMethod         string
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
type AuthenticationVerifierStore = core.KeyValueStore[string, []byte]

type AuthorizationService struct {
	authorizationStore          AuthorizationStore
	codeStore                   CodeStore
	authenticationVerifierStore AuthenticationVerifierStore
	clientService               *ClientService
	logger                      *zap.SugaredLogger
}

func NewAuthorizationService(authorizationStore AuthorizationStore, codeStore CodeStore, authenticationVerifierStore AuthenticationVerifierStore, clientService *ClientService, logger *zap.SugaredLogger) *AuthorizationService {
	return &AuthorizationService{authorizationStore: authorizationStore, codeStore: codeStore, authenticationVerifierStore: authenticationVerifierStore, clientService: clientService, logger: logger}
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

	request.authenticationCode = utils.HashShake256([]byte(randomString(32)))

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

func (s *AuthorizationService) VerifyAuthentication(uuid string, authenticationVerifier string) ResponseUriBuilder {
	s.logger.Debugw("Continuing 'authorization_code' flow", "authorizationId", uuid)
	hashedAuthenticationVerifier, err := s.authenticationVerifierStore.Get(uuid)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: "",
			State:       "",
			Error:       "server_error",
		}
	}

	decodedVerifier, err := utils.DecodeBase64([]byte(authenticationVerifier))
	if err != nil {
		fmt.Println("FAILED HERE!", authenticationVerifier, err)
		return &AuthorizationError{
			RedirectUri: "",
			State:       "",
			Error:       "bad_request",
		}
	}

	fmt.Println(decodedVerifier, hashedAuthenticationVerifier)

	if string(hashedAuthenticationVerifier) != string(utils.HashShake256(decodedVerifier)) {
		fmt.Println("FAILED HERE!!")
		return &AuthorizationError{
			RedirectUri: "",
			State:       "",
			Error:       "unauthenticated",
		}
	}

	return s.succeed(uuid)
}

func (s *AuthorizationService) succeed(uuid string) ResponseUriBuilder {
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
	ctx.SetCookie("authorization_id", uuid, 0, "/", "", true, true)
	ctx.Redirect(302, "/")
}

func (c *AuthorizationController) succeed(ctx *gin.Context) {
	uuid, err := ctx.Cookie("authorization_id")
	if err != nil {
		ctx.Redirect(302, "/")
		return
	}
	authenticationVerifier, err := ctx.Cookie("authentication_verifier")
	if err != nil {
		ctx.Redirect(302, "/")
		return
	}
	response := c.authorizationService.VerifyAuthentication(uuid, authenticationVerifier)
	ctx.SetCookie("authorization", "", -1, "/", "", false, true)
	ctx.Redirect(302, response.BuildResponseURI())
}
