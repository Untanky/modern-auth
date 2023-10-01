package main

import (
	"github.com/Untanky/modern-auth/apps/oauth2/internal/oauth2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (controller *controller) listClients(ctx *gin.Context) {
	clients, err := controller.clientService.List(ctx)
	if err != nil {
		// TODO: handle Forbidden
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var dtos = make([]oauth2.ClientDTO, 0, len(clients))
	for _, client := range clients {
		dtos = append(dtos, oauth2.ClientDTO{
			ID:           client.ID,
			Scopes:       client.Scopes,
			RedirectURIs: client.RedirectURIs,
		})
	}
	ctx.JSON(http.StatusOK, dtos)
}

func (controller *controller) getClient(ctx *gin.Context) {
	id := ctx.Param("id")
	client, err := controller.clientService.FindById(ctx, id)
	if err != nil {
		// TODO: handle NotFound and Forbidden
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, &oauth2.ClientDTO{
		ID:           client.ID,
		Scopes:       client.Scopes,
		RedirectURIs: client.RedirectURIs,
	})
}

func (controller *controller) createClient(ctx *gin.Context) {
	var dto oauth2.ClientDTO
	err := ctx.BindJSON(&dto)
	if err != nil {
		// TODO: handle Forbidden
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client, err := controller.clientService.Create(ctx, dto)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, &oauth2.ClientDTO{
		ID:           client.ID,
		Scopes:       client.Scopes,
		RedirectURIs: client.RedirectURIs,
	})
}

func (controller *controller) deleteClient(ctx *gin.Context) {
	id := ctx.Param("id")
	err := controller.clientService.Delete(ctx, id)
	if err != nil {
		// TODO: handle NotFound and Forbidden
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}
