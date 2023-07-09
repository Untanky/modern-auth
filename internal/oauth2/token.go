package oauth2

import (
	"crypto"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type authoritzationGrant struct {
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
	AccessToken  string
	TokenType    string
	ExpiresIn    int
	Scope        string
	RefreshToken string
}

type TokenError struct {
	Error            string
	ErrorDescription string
}

type OAuthTokenService struct {
	codeStore CodeStore
}

func NewOAuthTokenService(codeStore CodeStore) *OAuthTokenService {
	return &OAuthTokenService{codeStore: codeStore}
}

func (s *OAuthTokenService) Token(request TokenRequest) (*TokenResponse, *TokenError) {
	// validate request
	var grant *authoritzationGrant
	var err *TokenError
	switch actualRequest := request.(type) {
	case *AuthorizationCodeTokenRequest:
		grant, err = s.authorizationCodeToken(actualRequest)
	default:
		log.Println(request)
		return nil, &TokenError{
			Error:            "unsupported_grant_type",
			ErrorDescription: "grant type not supported",
		}
	}

	if err != nil {
		return nil, err
	}

	// generate token
	accessToken := "abc"
	refreshToken := "def"

	// store grant with token

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(-time.Since(grant.ExpiresAt).Seconds()),
		RefreshToken: refreshToken,
		Scope:        grant.Scope,
	}, nil
}

func (s *OAuthTokenService) authorizationCodeToken(tokenRequest *AuthorizationCodeTokenRequest) (*authoritzationGrant, *TokenError) {
	if tokenRequest.GrantType != "authorization_code" {
		return nil, &TokenError{
			Error:            "unsupported_grant_type",
			ErrorDescription: "grant type not supported",
		}
	}

	authorizationRequest, err := s.codeStore.Get(tokenRequest.Code)
	if err != nil {
		return nil, &TokenError{
			Error:            "invalid_grant",
			ErrorDescription: "authorization code not found",
		}
	}

	if authorizationRequest.ClientId != tokenRequest.ClientId {
		return nil, &TokenError{
			Error:            "invalid_grant",
			ErrorDescription: "client id does not match",
		}
	}

	if authorizationRequest.RedirectUri != tokenRequest.RedirectUri {
		return nil, &TokenError{
			Error:            "invalid_grant",
			ErrorDescription: "redirect uri does not match",
		}
	}

	if authorizationRequest.CodeChallenge != "" && authorizationRequest.CodeChallenge != hash(tokenRequest.CodeVerifier) {
		return nil, &TokenError{
			Error:            "invalid_grant",
			ErrorDescription: "code verifier invalid",
		}
	}

	return &authoritzationGrant{
		IssueRefreshToken: true,
		ID:                uuid.New(),
		Scope:             authorizationRequest.Scope,
		Client:            &Client{},
		Subject:           "xyz",
		IssuedAt:          time.Now(),
		ExpiresAt:         time.Now().Add(time.Hour),
	}, nil
}

type TokenController struct {
	tokenService *OAuthTokenService
}

func NewTokenController(tokenService *OAuthTokenService) *TokenController {
	return &TokenController{tokenService: tokenService}
}

func (c *TokenController) RegisterRoutes(router gin.IRouter) {
	router.POST("/token", c.token)
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

// function to hash a string
func hash(s string) string {
	hash := crypto.SHA256.New()
	hash.Write([]byte(s))
	return string(hash.Sum(nil))
}
