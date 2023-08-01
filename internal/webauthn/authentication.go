package webauthn

import (
	"context"
	"log"
	"net/http"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/user"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

const rpId = "localhost" // TODO: make customizable

type InitiateAuthenticationRequest struct {
	UserId string `json:"userId"`
}

type InitiateAuthenticationResponse struct {
	OptionId         string                            `json:"optionId"`
	PublicKeyOptions PublicKeyCredentialRequestOptions `json:"publicKey"`
}

type PublicKeyCredentialRequestOptions struct {
	Challenge                 string                      `json:"challenge"`
	RelyingParty              RelyingPartyOptions         `json:"rp"`
	User                      UserOptions                 `json:"user"`
	PublicKeyCredentialParams []PublicKeyCredentialParams `json:"pubKeyCredParams"`
	AuthenticationSelection   AuthenticationSelection     `json:"authenticatorSelection"`
	Timeout                   int                         `json:"timeout"`
	Attestation               string                      `json:"attestation"`
}

type RelyingPartyOptions struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserOptions struct {
	Id          []byte `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type PublicKeyCredentialParams struct {
	Type string `json:"type"`
	Alg  int    `json:"alg"`
}

type AuthenticationSelection struct {
	AuthenticatorAttachment string `json:"authenticatorAttachment"`
	RequireResidentKey      bool   `json:"requireResidentKey"`
	UserVerification        string `json:"userVerification"`
}

type CreateCredentialRequest struct {
	OptionId string                   `json:"optionId"`
	Id       string                   `json:"id"`
	Type     string                   `json:"type"`
	Response CreateCredentialResponse `json:"response"`
}

type UserService interface {
	GetUserByUserID(ctx context.Context, userId []byte) (*user.User, error)
	CreateUser(ctx context.Context, user *user.User) error
}

type CredentialService interface {
	GetCredentialByCredentialID(ctx context.Context, creadentialId []byte) (*user.Credential, error)
	CreateCredential(ctx context.Context, credential *user.Credential) error
}

type AuthenticationService struct {
	initAuthenticationStore core.KeyValueStore[string, InitiateAuthenticationResponse]
	userService             UserService
	credentialService       CredentialService
}

func NewAuthenticationService(initAuthenticationStore core.KeyValueStore[string, InitiateAuthenticationResponse], userService UserService, credentialService CredentialService) *AuthenticationService {
	return &AuthenticationService{
		initAuthenticationStore: initAuthenticationStore,
		userService:             userService,
		credentialService:       credentialService,
	}
}

func (s *AuthenticationService) InitiateAuthentication(request *InitiateAuthenticationRequest) (*InitiateAuthenticationResponse, error) {
	id := uuid.New().String()

	initResponse := &InitiateAuthenticationResponse{
		OptionId: id,
		PublicKeyOptions: PublicKeyCredentialRequestOptions{
			// TODO: randomly generate challenge
			Challenge: "1234567890",
			RelyingParty: RelyingPartyOptions{
				Id:   rpId,
				Name: "Modern Auth",
			},
			User: UserOptions{
				Id:          []byte(request.UserId),
				Name:        request.UserId,
				DisplayName: request.UserId,
			},
			PublicKeyCredentialParams: []PublicKeyCredentialParams{
				{
					Type: "public-key",
					Alg:  -7,
				},
			},
			AuthenticationSelection: AuthenticationSelection{
				AuthenticatorAttachment: "all",
				RequireResidentKey:      false,
				UserVerification:        "preferred",
			},
			Timeout:     60000,
			Attestation: "indirect",
		},
	}

	err := s.initAuthenticationStore.Set(id, initResponse)
	if err != nil {
		return nil, err
	}

	return initResponse, nil
}

func (s *AuthenticationService) Register(ctx context.Context, id string, response *CreateCredentialResponse) error {
	options, err := s.initAuthenticationStore.Get(id)
	if err != nil {
		return err
	}

	credential, err := response.Validate(options)
	if err != nil {
		return err
	}

	// TODO: assess trust of the authenticator

	userInstance := user.User{
		UserID: options.PublicKeyOptions.User.Id,
		Status: "active",
	}

	err = s.userService.CreateUser(ctx, &userInstance)
	if err != nil {
		return err
	}

	credential.User = userInstance

	err = s.credentialService.CreateCredential(ctx, credential)
	if err != nil {
		return err
	}

	return nil
}

type AuthenticationController struct {
	service *AuthenticationService
}

func NewAuthenticationController(service *AuthenticationService) *AuthenticationController {
	return &AuthenticationController{
		service: service,
	}
}

func (c *AuthenticationController) RegisterRoutes(router gin.IRoutes) {
	router.POST("/authentication/initiate", c.initiateAuthentication)
	router.POST("/authentication/create", c.createCredential)
}

func (c *AuthenticationController) initiateAuthentication(ctx *gin.Context) {
	var request InitiateAuthenticationRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	response, err := c.service.InitiateAuthentication(&request)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, response)
}

func (c *AuthenticationController) createCredential(ctx *gin.Context) {
	var request CreateCredentialRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	err = c.service.Register(ctx.Request.Context(), request.OptionId, &request.Response)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}
