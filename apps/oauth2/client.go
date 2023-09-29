package main

import (
	"github.com/Untanky/modern-auth/internal/oauth2"
	"github.com/gin-gonic/gin"
	"net/http"
)

var clientService *oauth2.ClientService

func listClients(ctx *gin.Context) {
	clients, err := clientService.List(ctx)
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

func getClient(ctx *gin.Context) {
	id := ctx.Param("id")
	client, err := clientService.FindById(ctx, id)
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

func createClient(ctx *gin.Context) {
	var dto oauth2.ClientDTO
	err := ctx.BindJSON(&dto)
	if err != nil {
		// TODO: handle Forbidden
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client, err := clientService.Create(ctx, dto)
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

func deleteClient(ctx *gin.Context) {
	id := ctx.Param("id")
	err := clientService.Delete(ctx, id)
	if err != nil {
		// TODO: handle NotFound and Forbidden
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Status(http.StatusOK)
}
