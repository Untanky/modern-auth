package oauth2

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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
	ErrorType   string
	Description string
}

func (e *AuthorizationError) BuildResponseURI() string {
	return fmt.Sprintf("%s?error=%s&error_description=%s&state=%s&iss=%s", e.RedirectUri, e.ErrorType, e.Description, e.State, e.Issuer)
}

func (e *AuthorizationError) Error() string {
	return fmt.Sprintf("AuthorizationError: %s", e.Description)
}

type AuthorizationStore = core.KeyValueStore[string, *AuthorizationRequest]
type CodeStore = core.KeyValueStore[string, *AuthorizationRequest]
type AuthenticationVerifierStore = core.KeyValueStore[string, []byte]

type AuthorizationService struct {
	authorizationStore          AuthorizationStore
	codeStore                   CodeStore
	authenticationVerifierStore AuthenticationVerifierStore
	clientService               *ClientService
	logger                      *slog.Logger
	authorizationCodeInit       metric.Int64Counter
	authorizationCodeSuccess    metric.Int64Counter
}

func NewAuthorizationService(authorizationStore AuthorizationStore, codeStore CodeStore, authenticationVerifierStore AuthenticationVerifierStore, clientService *ClientService, authorizationCodeInit metric.Int64Counter, authorizationCodeSuccess metric.Int64Counter) *AuthorizationService {
	logger := slog.Default().With(slog.String("service", "authorization"))

	return &AuthorizationService{
		authorizationStore:          authorizationStore,
		codeStore:                   codeStore,
		authenticationVerifierStore: authenticationVerifierStore,
		clientService:               clientService,
		logger:                      logger,
		authorizationCodeInit:       authorizationCodeInit,
		authorizationCodeSuccess:    authorizationCodeSuccess,
	}
}

func (s *AuthorizationService) Authorize(ctx context.Context, request *AuthorizationRequest) (string, *AuthorizationError) {
	s.logger.Debug("Beginning 'authorization_code' flow")
	client, err := s.clientService.FindById(ctx, request.ClientId)
	if err != nil {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			ErrorType:   "invalid_client",
			Description: "client not found",
		}
	}

	if request.RedirectUri == "" {
		request.RedirectUri = client.RedirectURIs[0]
	} else if !client.ValidateRedirectURI(ctx, request.RedirectUri) {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			ErrorType:   "invalid_request",
			Description: "redirect_uri not allowed",
		}
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			ErrorType:   "server_error",
		}
	}

	request.authenticationCode = utils.HashShake256([]byte(randomString(32)))

	stringUuid := uuid.String()
	request.id = uuid
	err = s.authorizationStore.WithContext(ctx).Set(stringUuid, request)
	if err != nil {
		return "", &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			ErrorType:   "server_error",
		}
	}
	s.logger.Info("Initialized 'authorization_code' flow", "authorizationId", stringUuid)
	s.authorizationCodeInit.Add(context.Background(), 1, metric.WithAttributes(attribute.Key("client_id").String(request.ClientId)))

	return stringUuid, nil
}

func (s *AuthorizationService) VerifyAuthentication(ctx context.Context, uuid string, authenticationVerifier string) ResponseUriBuilder {
	s.logger.Debug("Continuing 'authorization_code' flow", "authorizationId", uuid)
	hashedAuthenticationVerifier, err := s.authenticationVerifierStore.Get(uuid)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: "",
			State:       "",
			ErrorType:   "server_error",
		}
	}

	decodedVerifier, err := utils.DecodeBase64([]byte(authenticationVerifier))
	if err != nil {
		return &AuthorizationError{
			RedirectUri: "",
			State:       "",
			ErrorType:   "bad_request",
		}
	}

	if string(hashedAuthenticationVerifier) != string(utils.HashShake256(decodedVerifier)) {
		return &AuthorizationError{
			RedirectUri: "",
			State:       "",
			ErrorType:   "unauthenticated",
		}
	}

	return s.succeed(ctx, uuid)
}

func (s *AuthorizationService) succeed(ctx context.Context, uuid string) ResponseUriBuilder {
	s.logger.Debug("Continuing 'authorization_code' flow", "authorizationId", uuid)
	store := s.authorizationStore.WithContext(ctx)
	request, err := store.Get(uuid)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			ErrorType:   "server_error",
		}
	}

	err = store.Delete(uuid)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			ErrorType:   "server_error",
		}
	}

	code := randomString(32)
	err = s.codeStore.WithContext(ctx).Set(code, request)
	if err != nil {
		return &AuthorizationError{
			RedirectUri: request.RedirectUri,
			State:       request.State,
			ErrorType:   "server_error",
		}
	}
	s.logger.Info("Generated authorization code", "authroizationId", uuid)
	s.authorizationCodeSuccess.Add(ctx, 1, metric.WithAttributes(attribute.Key("client_id").String(request.ClientId)))

	return &AuthorizationResponse{
		RedirectUri: request.RedirectUri,
		Code:        code,
		State:       request.State,
		Issuer:      "https://localhost:8080",
	}
}
