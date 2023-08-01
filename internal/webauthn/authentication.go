package webauthn

import (
	"context"
	"log"
	"net/http"

	"github.com/Untanky/modern-auth/internal/user"

	"github.com/gin-gonic/gin"
)

const rpId = "localhost" // TODO: make customizable

type InitiateAuthenticationRequest struct {
	UserId string `json:"userId"`
}

type InitiateAuthenticationResponse struct {
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
	Id       string                   `json:"id"`
	Type     string                   `json:"type"`
	Response CreateCredentialResponse `json:"response"`
}

type UserService interface {
	GetUserByUserId(ctx context.Context, userId string) (*user.User, error)
	CreateUser(ctx context.Context, user *user.User) error
}

type Credential interface{}

type CredentialService interface {
	GetByCredentialId(id string) (Credential, error)
	Create(credential Credential) error
}

type AuthenticationService struct {
	userService UserService
}

func NewAuthenticationService(userService UserService) *AuthenticationService {
	return &AuthenticationService{
		userService: userService,
	}
}

func (s *AuthenticationService) InitiateAuthentication(request *InitiateAuthenticationRequest) *InitiateAuthenticationResponse {
	return &InitiateAuthenticationResponse{
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
}

func (s *AuthenticationService) Register(response *CreateCredentialResponse) error {
	options := s.InitiateAuthentication(&InitiateAuthenticationRequest{})

	err := response.Validate(options)
	if err != nil {
		return err
	}

	// TODO: assess trust of the authenticator

	user := user.User{}

	// TODO: create user
	err = s.userService.CreateUser(context.TODO(), &user)
	if err != nil {
		return err
	}

	// TODO: link credential to user

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

	response := c.service.InitiateAuthentication(&request)

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

	err = c.service.Register(&request.Response)
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
