package main

import (
	"fmt"
	"github.com/Untanky/modern-auth/apps/oauth2/internal/oauth2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller *controller) issueToken(ctx *gin.Context) {
	tokenRequest, err := controller.parseTokenRequest(ctx)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	tokenResponse, tokenError := controller.tokenService.Token(ctx.Request.Context(), tokenRequest)
	if tokenError != nil {
		ctx.AbortWithError(http.StatusBadRequest, tokenError)
		return
	}

	ctx.JSON(http.StatusOK, tokenResponse)
}

func (controller *controller) parseTokenRequest(ctx *gin.Context) (oauth2.TokenRequest, error) {
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

func (controller *controller) returnGrant(ctx *gin.Context) {
	grant, _ := ctx.Get("grant")

	ctx.JSON(http.StatusOK, grant)
}
