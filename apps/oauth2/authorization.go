package main

import (
	"github.com/Untanky/modern-auth/apps/oauth2/internal/oauth2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller *controller) startAuthorization(ctx *gin.Context) {
	request := &oauth2.AuthorizationRequest{}
	err := ctx.ShouldBindQuery(request)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	uuid, authorizationErr := controller.authorizationService.Authorize(ctx.Request.Context(), request)
	if authorizationErr != nil {
		ctx.Error(authorizationErr)
		ctx.Redirect(302, authorizationErr.BuildResponseURI())
		return
	}
	ctx.SetCookie("authorization_id", uuid, 0, "/", "", true, true)
	ctx.Redirect(302, "/")
}

func (controller *controller) succeedAuthorization(ctx *gin.Context) {
	uuid, err := ctx.Cookie("authorization_id")
	if err != nil {
		ctx.Error(err)
		ctx.Redirect(302, "/")
		return
	}
	authenticationVerifier, err := ctx.Cookie("authentication_verifier")
	if err != nil {
		ctx.Error(err)
		ctx.Redirect(302, "/")
		return
	}
	response := controller.authorizationService.VerifyAuthentication(ctx.Request.Context(), uuid, authenticationVerifier)
	ctx.SetCookie("authorization", "", -1, "/", "", false, true)
	ctx.Redirect(302, response.BuildResponseURI())
}
