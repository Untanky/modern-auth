package oauth2

import (
	"net/http"
	"strings"
	"time"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type TokenRequest interface {
}

type tokenRequest struct {
	GrantType string `form:"grant_type" binding:"required"`
	ClientId  string `form:"client_id" binding:"required"`
}

type AuthorizationCodeTokenRequest struct {
	tokenRequest
	RedirectUri  string `form:"redirect_uri"`
	Code         string `form:"code" binding:"required"`
	CodeVerifier string `form:"code_verifier"`
}

type AuthorizationGrant struct {
	IssueRefreshToken bool
	ID                uuid.UUID
	Scope             string
	Client            *Client
	Subject           string
	IssuedAt          time.Time
	ExpiresAt         time.Time
	NotBefore         time.Time
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" binding:"required"`
	TokenType    string `json:"token_type" binding:"required"`
	ExpiresIn    int    `json:"expires_in" binding:"required"`
	Scope        string `json:"scope" binding:"required"`
	RefreshToken string `json:"refresh_token"`
}

type TokenError struct {
	ErrorTitle       string `json:"error" binding:"required"`
	ErrorDescription string `json:"error_description" binding:"required"`
}

func (e *TokenError) Error() string {
	return e.ErrorDescription
}

type OAuthTokenService struct {
	codeStore           CodeStore
	accessTokenHandler  TokenHandler
	refreshTokenHandler TokenHandler
	logger              *zap.SugaredLogger
}

func NewOAuthTokenService(codeStore CodeStore, accessTokenHandler TokenHandler, refreshTokenHandler TokenHandler, logger *zap.SugaredLogger) *OAuthTokenService {
	return &OAuthTokenService{
		codeStore:           codeStore,
		accessTokenHandler:  accessTokenHandler,
		refreshTokenHandler: refreshTokenHandler,
		logger:              logger,
	}
}

func (s *OAuthTokenService) Token(request TokenRequest) (*TokenResponse, *TokenError) {
	// validate request
	var grant *AuthorizationGrant
	var err *TokenError
	s.logger.Debugw("Handling token request")

	switch actualRequest := request.(type) {
	case *AuthorizationCodeTokenRequest:
		grant, err = s.authorizationCodeToken(actualRequest)
	default:
		return nil, &TokenError{
			ErrorTitle:       "unsupported_grant_type",
			ErrorDescription: "grant type not supported",
		}
	}

	if err != nil {
		return nil, err
	}

	accessToken, e := s.accessTokenHandler.GenerateToken(grant)
	if e != nil {
		return nil, &TokenError{
			ErrorTitle:       "server_error",
			ErrorDescription: "failed to generate access token",
		}
	}
	s.logger.Debugw("Generated access token", "authorization_id", grant.ID)

	var refreshToken string
	if grant.IssueRefreshToken {
		refreshToken, e = s.refreshTokenHandler.GenerateToken(grant)
		if e != nil {
			s.logger.Warnw("Refresh token generation failed", "err", err, "authorization_id", grant.ID)
			refreshToken = ""
		} else {
			s.logger.Debugw("Generated refresh token", "authorization_id", grant.ID)
		}
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(-time.Since(grant.ExpiresAt).Seconds()),
		RefreshToken: refreshToken,
		Scope:        grant.Scope,
	}, nil
}

func (s *OAuthTokenService) authorizationCodeToken(tokenRequest *AuthorizationCodeTokenRequest) (*AuthorizationGrant, *TokenError) {
	if tokenRequest.GrantType != "authorization_code" {
		return nil, &TokenError{
			ErrorTitle:       "unsupported_grant_type",
			ErrorDescription: "grant type not supported",
		}
	}

	authorizationRequest, err := s.codeStore.Get(tokenRequest.Code)
	if err != nil {
		return nil, &TokenError{
			ErrorTitle:       "invalid_grant",
			ErrorDescription: "authorization code not found",
		}
	}

	if authorizationRequest.ClientId != tokenRequest.ClientId {
		return nil, &TokenError{
			ErrorTitle:       "invalid_grant",
			ErrorDescription: "client id does not match",
		}
	}

	if authorizationRequest.RedirectUri != tokenRequest.RedirectUri {
		return nil, &TokenError{
			ErrorTitle:       "invalid_grant",
			ErrorDescription: "redirect uri does not match",
		}
	}

	if authorizationRequest.CodeChallenge != "" && authorizationRequest.CodeChallenge != hash(tokenRequest.CodeVerifier) {
		return nil, &TokenError{
			ErrorTitle:       "invalid_grant",
			ErrorDescription: "code verifier invalid",
		}
	}

	s.logger.Infow("Successfully validated 'authorization_code' token request",
		"client_id", tokenRequest.ClientId,
		"grant_type", "authorization_code",
		"authorization_id", authorizationRequest.id,
	)

	return &AuthorizationGrant{
		IssueRefreshToken: true,
		ID:                authorizationRequest.id,
		Scope:             authorizationRequest.Scope,
		Client:            &Client{},
		Subject:           "xyz",
		IssuedAt:          time.Now(),
		ExpiresAt:         time.Now().Add(time.Hour),
	}, nil
}

func (s *OAuthTokenService) Validate(token string) (*AuthorizationGrant, error) {
	return s.accessTokenHandler.Validate(token)
}

type TokenHandler interface {
	GenerateToken(grant *AuthorizationGrant) (string, error)
	Validate(token string) (*AuthorizationGrant, error)
}

type TokenStore = core.KeyValueStore[string, AuthorizationGrant]

type RandomTokenHandler struct {
	tokenSize int
	store     TokenStore
	logger    *zap.SugaredLogger
}

func NewRandomTokenHandler(tokenSize int, store TokenStore, logger *zap.SugaredLogger) *RandomTokenHandler {
	return &RandomTokenHandler{tokenSize: tokenSize, store: store, logger: logger}
}

func (h *RandomTokenHandler) GenerateToken(grant *AuthorizationGrant) (string, error) {
	token := randomString(h.tokenSize)
	err := h.store.Set(token, grant)
	if err != nil {
		return "", err
	}
	h.logger.Debugw("Generated randomized token", "authorization_id", grant.ID)
	return token, nil
}

func (h *RandomTokenHandler) Validate(token string) (*AuthorizationGrant, error) {
	grant, err := h.store.Get(token)
	if err != nil {
		return nil, err
	}
	h.logger.Debugw("Successfully validated token", "authorization_id", grant.ID)
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
	if err := ctx.ShouldBind(&temp); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}
	if temp.GrantType == "authorization_code" {
		var authorizationCodeTokenRequest AuthorizationCodeTokenRequest
		if err := ctx.ShouldBind(&authorizationCodeTokenRequest); err != nil {
			ctx.JSON(http.StatusBadRequest, err)
			return
		}
		tokenRequest = &authorizationCodeTokenRequest
	}

	tokenResponse, tokenError := c.tokenService.Token(tokenRequest)
	if tokenError != nil {
		ctx.JSON(http.StatusBadRequest, tokenError)
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse)
}

func (c *TokenController) validate(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderParts) != 2 || authorizationHeaderParts[0] != "Bearer" {
		ctx.JSON(http.StatusBadRequest, "invalid authorization header")
		return
	}
	token := authorizationHeaderParts[1]

	grant, err := c.tokenService.Validate(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, grant)
}
