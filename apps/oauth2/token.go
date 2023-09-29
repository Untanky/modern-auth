package main

import (
	"fmt"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"
)

var tokenService *oauth2.OAuthTokenService

func issueToken(ctx *gin.Context) {
	tokenRequest, err := parseTokenRequest(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tokenResponse, tokenError := tokenService.Token(ctx.Request.Context(), tokenRequest)
	if tokenError != nil {
		ctx.AbortWithError(http.StatusBadRequest, tokenError)
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse)
}

func parseTokenRequest(ctx *gin.Context) (oauth2.TokenRequest, error) {
	var temp struct {
		GrantType string `form:"grant_type" binding:"required"`
	}
	if err := ctx.ShouldBind(&temp); err != nil {
		return nil, err
	}

	switch temp.GrantType {
	case "authorization_code":
		var authorizationCodeTokenRequest *oauth2.AuthorizationCodeTokenRequest
		err := ctx.ShouldBind(authorizationCodeTokenRequest)
		return authorizationCodeTokenRequest, err
	case "refresh_token":
		var refreshTokenRequest *oauth2.RefreshTokenRequest
		err := ctx.ShouldBind(refreshTokenRequest)
		return refreshTokenRequest, err
	default:
		return nil, fmt.Errorf("invalid grant type: %s", temp.GrantType)
	}
}

func validateToken(ctx *gin.Context) {
	authorizationHeader := ctx.GetHeader("Authorization")
	authorizationHeaderParts := strings.Split(authorizationHeader, " ")
	if len(authorizationHeaderParts) != 2 || authorizationHeaderParts[0] != "Bearer" {
		err := fmt.Errorf("invalid authorization header")
		ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	token := authorizationHeaderParts[1]

	grant, err := tokenService.Validate(ctx.Request.Context(), token)
	if err != nil {
		trace.SpanFromContext(ctx.Request.Context()).RecordError(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, grant)
}
