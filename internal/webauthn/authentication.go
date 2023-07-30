package webauthn

import (
	"fmt"
	"net/http"

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
	Challenge                 []byte                      `json:"challenge"`
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
	Id       string               `json:"id"`
	Type     string               `json:"type"`
	Response CreateCredentialData `json:"response"`
}

type CreateCredentialData struct {
	AttestationObject []byte `json:"attestationObject"`
	ClientDataJSON    []byte `json:"clientDataJSON"`
}

type UserService interface {
	IsUserIdAvailable(userId string) bool
	GetUser(userId string) (interface{}, error)
	CreateUser(user interface{}) error
}

type Credential interface{}

type CredentialService interface {
	GetByCredentialId(id string) (Credential, error)
	Create(credential Credential) error
}

type AuthenticationService struct {
}

func NewAuthenticationService() *AuthenticationService {
	return &AuthenticationService{}
}

func (s *AuthenticationService) InitiateAuthentication(request *InitiateAuthenticationRequest) *InitiateAuthenticationResponse {
	return &InitiateAuthenticationResponse{
		PublicKeyOptions: PublicKeyCredentialRequestOptions{
			// TODO: randomly generate challenge
			Challenge: []byte("1234567890"),
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

	// TODO: create user
	// TODO: link credential to user

	return nil
}

type AuthenticationController struct {
}

func NewAuthenticationController() *AuthenticationController {
	return &AuthenticationController{}
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

	response := NewAuthenticationService().InitiateAuthentication(&request)

	ctx.JSON(200, response)
}

func (c *AuthenticationController) createCredential(ctx *gin.Context) {
	var request CreateCredentialRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	fmt.Println(request)

	ctx.JSON(200, gin.H{
		"success": true,
	})
}
