package main

import (
	"fmt"
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/gin-gonic/gin"
	"net/http"
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

func returnGrant(ctx *gin.Context) {
	grant, _ := ctx.Get("grant")

	ctx.JSON(http.StatusOK, grant)
}
