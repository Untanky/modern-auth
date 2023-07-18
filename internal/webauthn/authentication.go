package webauthn

import "github.com/gin-gonic/gin"

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

type AuthenticationController struct {
}

func NewAuthenticationController() *AuthenticationController {
	return &AuthenticationController{}
}

func (c *AuthenticationController) RegisterRoutes(router gin.IRoutes) {
	router.POST("/authentication/initiate", c.initiateAuthentication)
}

func (c *AuthenticationController) initiateAuthentication(ctx *gin.Context) {
	ctx.JSON(200, InitiateAuthenticationResponse{
		PublicKeyOptions: PublicKeyCredentialRequestOptions{
			Challenge: []byte("1234567890"),
			RelyingParty: RelyingPartyOptions{
				Id:   "localhost",
				Name: "localhost",
			},
			User: UserOptions{
				Id:          []byte("1234567890"),
				Name:        "test",
				DisplayName: "test",
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
			Attestation: "none",
		},
	})
}
