package oauth2

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type TokenRequest interface {
	GetGrantType() string
}

type tokenRequest struct {
	GrantType string `form:"grant_type" binding:"required"`
	ClientId  string `form:"client_id" binding:"required"`
}

func (r *tokenRequest) GetGrantType() string {
	return r.GrantType
}

type AuthorizationCodeTokenRequest struct {
	tokenRequest
	RedirectUri  string `form:"redirect_uri"`
	Code         string `form:"code" binding:"required"`
	CodeVerifier string `form:"code_verifier"`
}

type RefreshTokenRequest struct {
	tokenRequest
	RefreshToken string `form:"refresh_token" binding:"required"`
}

type timestamp time.Time

func (t *timestamp) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprint(time.Time(*t).Unix())), nil
}

type AuthorizationGrant struct {
	IssueRefreshToken bool      `json:"-"`
	ID                uuid.UUID `json:"-"`
	Scope             string    `json:"scope"`
	ClientId          string    `json:"client_id"`
	SubjectId         string    `json:"sub"`
	IssuedAt          timestamp `json:"iat"`
	ExpiresAt         timestamp `json:"exp"`
	NotBefore         timestamp `json:"nbf,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	TokenType    string `json:"token_type" binding:"required"`
	ExpiresIn    int    `json:"expires_in" binding:"required"`
	Scope        string `json:"scope" binding:"required"`
	RefreshToken string `json:"refresh_token"`
}

type TokenError struct {
	ErrorType        string `json:"error" binding:"required"`
	ErrorDescription string `json:"error_description" binding:"required"`
}

func (e *TokenError) Error() string {
	return e.ErrorDescription
}

type OAuthTokenService struct {
	codeStore              CodeStore
	accessTokenHandler     TokenHandler
	refreshTokenHandler    TokenHandler
	logger                 *slog.Logger
	tokenRequestInstrument metric.Int64Counter
}

func NewOAuthTokenService(codeStore CodeStore, accessTokenHandler TokenHandler, refreshTokenHandler TokenHandler, tokenRequestInstrument metric.Int64Counter) *OAuthTokenService {
	logger := slog.Default().With(slog.String("service", "oauth-token"))

	return &OAuthTokenService{
		codeStore:              codeStore,
		accessTokenHandler:     accessTokenHandler,
		refreshTokenHandler:    refreshTokenHandler,
		logger:                 logger,
		tokenRequestInstrument: tokenRequestInstrument,
	}
}

func (s *OAuthTokenService) Token(ctx context.Context, request TokenRequest) (*TokenResponse, *TokenError) {
	// validate request
	var grant *AuthorizationGrant
	var err *TokenError
	s.logger.Debug("Handling token request")

	switch actualRequest := request.(type) {
	case *AuthorizationCodeTokenRequest:
		grant, err = s.authorizationCodeToken(ctx, actualRequest)
	case *RefreshTokenRequest:
		grant, err = s.refreshToken(ctx, actualRequest)
	default:
		return nil, &TokenError{
			ErrorType:        "unsupported_grant_type",
			ErrorDescription: "grant type not supported",
		}
	}

	if err != nil {
		s.logger.Warn("Token request failed", "err", err)
		return nil, err
	}

	accessToken, e := s.accessTokenHandler.GenerateToken(ctx, grant)
	if e != nil {
		return nil, &TokenError{
			ErrorType:        "server_error",
			ErrorDescription: "failed to generate access token",
		}
	}
	s.logger.Debug("Generated access token", "authorization_id", grant.ID)

	var refreshToken string
	if grant.IssueRefreshToken {
		refreshToken, e = s.refreshTokenHandler.GenerateToken(ctx, grant)
		if e != nil {
			s.logger.Warn("Refresh token generation failed", "err", err, "authorization_id", grant.ID)
			refreshToken = ""
		} else {
			s.logger.Debug("Generated refresh token", "authorization_id", grant.ID)
		}
	}
	s.tokenRequestInstrument.Add(ctx, 1, metric.WithAttributes(attribute.Key("client_id").String(grant.ClientId), attribute.Key("grant_type").String(request.GetGrantType())))

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(-time.Since(time.Time(grant.ExpiresAt)).Seconds()),
		RefreshToken: refreshToken,
		Scope:        grant.Scope,
	}, nil
}

func (s *OAuthTokenService) authorizationCodeToken(ctx context.Context, tokenRequest *AuthorizationCodeTokenRequest) (*AuthorizationGrant, *TokenError) {
	if tokenRequest.GrantType != "authorization_code" {
		return nil, &TokenError{
			ErrorType:        "unsupported_grant_type",
			ErrorDescription: "grant type not supported",
		}
	}

	authorizationRequest, err := s.codeStore.WithContext(ctx).Get(tokenRequest.Code)
	if err != nil {
		return nil, &TokenError{
			ErrorType:        "invalid_grant",
			ErrorDescription: "authorization code not found",
		}
	}

	if authorizationRequest.ClientId != tokenRequest.ClientId {
		return nil, &TokenError{
			ErrorType:        "invalid_grant",
			ErrorDescription: "client id does not match",
		}
	}

	if authorizationRequest.RedirectUri != tokenRequest.RedirectUri {
		return nil, &TokenError{
			ErrorType:        "invalid_grant",
			ErrorDescription: "redirect uri does not match",
		}
	}

	if authorizationRequest.CodeChallenge != "" && authorizationRequest.CodeChallenge != hash(tokenRequest.CodeVerifier) {
		return nil, &TokenError{
			ErrorType:        "invalid_grant",
			ErrorDescription: "code verifier invalid",
		}
	}

	s.logger.Info("Successfully validated 'authorization_code' token request",
		"client_id", tokenRequest.ClientId,
		"grant_type", "authorization_code",
		"authorization_id", authorizationRequest.id,
	)

	return &AuthorizationGrant{
		IssueRefreshToken: true,
		ID:                authorizationRequest.id,
		Scope:             authorizationRequest.Scope,
		ClientId:          authorizationRequest.ClientId,
		SubjectId:         "xyz",
		IssuedAt:          timestamp(time.Now()),
		ExpiresAt:         timestamp(time.Now().Add(time.Hour)),
	}, nil
}

func (s *OAuthTokenService) refreshToken(ctx context.Context, tokenRequest *RefreshTokenRequest) (*AuthorizationGrant, *TokenError) {
	if tokenRequest.GrantType != "refresh_token" {
		return nil, &TokenError{
			ErrorType:        "unsupported_grant_type",
			ErrorDescription: "grant type not supported",
		}
	}

	grant, err := s.refreshTokenHandler.Validate(ctx, tokenRequest.RefreshToken)
	if err != nil {
		return nil, &TokenError{
			ErrorType:        "invalid_grant",
			ErrorDescription: "refresh token not found",
		}
	}

	if grant.ClientId != tokenRequest.ClientId {
		return nil, &TokenError{
			ErrorType:        "invalid_grant",
			ErrorDescription: "client id does not match",
		}
	}

	s.logger.Info("Successfully validated 'refresh_token' token request",
		"client_id", tokenRequest.ClientId,
		"grant_type", "refresh_token",
		"authorization_id", grant.ID,
	)

	return &AuthorizationGrant{
		IssueRefreshToken: false,
		ID:                grant.ID,
		Scope:             grant.Scope,
		ClientId:          grant.ClientId,
		SubjectId:         grant.SubjectId,
		IssuedAt:          timestamp(time.Now()),
		ExpiresAt:         timestamp(time.Now().Add(time.Hour)),
		NotBefore:         timestamp(time.Now()),
	}, nil
}

func (s *OAuthTokenService) Validate(ctx context.Context, token string) (*AuthorizationGrant, error) {
	return s.accessTokenHandler.Validate(ctx, token)
}

type TokenHandler interface {
	GenerateToken(ctx context.Context, grant *AuthorizationGrant) (string, error)
	Validate(ctx context.Context, token string) (*AuthorizationGrant, error)
}

type TokenStore = core.KeyValueStore[string, *AuthorizationGrant]

type RandomTokenHandler struct {
	tokenSize       int
	store           TokenStore
	logger          *slog.Logger
	tokensGenerated metric.Int64Counter
}

func NewRandomTokenHandler(tokenType string, tokenSize int, store TokenStore, tokensGenerated metric.Int64Counter) *RandomTokenHandler {
	logger := slog.Default().With(slog.String("service", "token-handler"), slog.String("type", tokenType))

	return &RandomTokenHandler{tokenSize: tokenSize, store: store, logger: logger, tokensGenerated: tokensGenerated}
}

func (h *RandomTokenHandler) GenerateToken(ctx context.Context, grant *AuthorizationGrant) (string, error) {
	token := randomString(h.tokenSize)
	secret := core.NewSecretValue(token)
	err := h.store.WithContext(ctx).Set(secret.String(), grant)
	if err != nil {
		return "", err
	}
	h.logger.Debug("Generated randomized token", "authorization_id", grant.ID)
	h.tokensGenerated.Add(context.Background(), 1, metric.WithAttributes(attribute.Key("client_id").String(grant.ClientId)))
	return token, nil
}

func (h *RandomTokenHandler) Validate(ctx context.Context, token string) (*AuthorizationGrant, error) {
	secret := core.NewSecretValue(token)
	grant, err := h.store.WithContext(ctx).Get(secret.String())
	if err != nil {
		return nil, err
	}
	h.logger.Debug("Successfully validated token", "authorization_id", grant.ID)
	return grant, err
}

type TokenController struct {
	tokenService *OAuthTokenService
}

func NewTokenController(tokenService *OAuthTokenService) *TokenController {
	return &TokenController{tokenService: tokenService}
}

func (c *TokenController) RegisterRoutes(router gin.IRouter) {
	router.POST("/token", c.token)
	router.POST("/token/validate", c.validate)
}

func (c *TokenController) token(ctx *gin.Context) {
	var temp tokenRequest
	var tokenRequest TokenRequest
	var err error
	if err = ctx.ShouldBind(&temp); err != nil {
		trace.SpanFromContext(ctx.Request.Context()).RecordError(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	switch temp.GrantType {
	case "authorization_code":
		var authorizationCodeTokenRequest AuthorizationCodeTokenRequest
		err = ctx.ShouldBind(&authorizationCodeTokenRequest)
		tokenRequest = &authorizationCodeTokenRequest
	case "refresh_token":
		var refreshTokenRequest RefreshTokenRequest
		err = ctx.ShouldBind(&refreshTokenRequest)
		tokenRequest = &refreshTokenRequest
	default:
		err = fmt.Errorf("invalid grant type: %s", temp.GrantType)
	}
	if err != nil {
		trace.SpanFromContext(ctx.Request.Context()).RecordError(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	tokenResponse, tokenError := c.tokenService.Token(ctx.Request.Context(), tokenRequest)
	if tokenError != nil {
		trace.SpanFromContext(ctx.Request.Context()).RecordError(tokenError)
		ctx.JSON(http.StatusBadRequest, tokenError)
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse)
}

func (c *TokenController) validate(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderParts) != 2 || authorizationHeaderParts[0] != "Bearer" {
		err := fmt.Errorf("invalid authorization header")
		trace.SpanFromContext(ctx.Request.Context()).RecordError(err)
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	token := authorizationHeaderParts[1]

	grant, err := c.tokenService.Validate(ctx.Request.Context(), token)
	if err != nil {
		trace.SpanFromContext(ctx.Request.Context()).RecordError(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, grant)
}
