package webauthn

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/fxamacker/cbor/v2"
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
	AttestationObject AttestationObject `json:"attestationObject"`
	ClientDataJSON    ClientData        `json:"clientDataJSON"`
}

type ClientData struct {
	Hash      []byte
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
	Origin    string `json:"origin"`
}

func (c *ClientData) UnmarshalJSON(base64Data []byte) error {
	data, err := utils.DecodeBase64(base64Data[1 : len(base64Data)-1])
	if err != nil {
		fmt.Println("ERROR", err)
		return err
	}

	var rawClientData map[string]interface{}
	err = json.Unmarshal(data[:len(data)-1], &rawClientData)
	if err != nil {
		return err
	}

	c.Type = rawClientData["type"].(string)
	c.Challenge = rawClientData["challenge"].(string)
	c.Origin = rawClientData["origin"].(string)
	c.Hash = utils.HashSHA256(data)
	return nil
}

type AttestationObject struct {
	AuthData       AuthData                    `json:"authData"`
	Format         string                      `json:"fmt"`
	AttestationRaw map[interface{}]interface{} `json:"attStmt"`
}

func (a *AttestationObject) UnmarshalJSON(base64Data []byte) error {
	data, err := utils.DecodeBase64(base64Data[1 : len(base64Data)-1])
	if err != nil {
		return err
	}

	var rawAttestationObject map[string]interface{}
	err = cbor.Unmarshal(data, &rawAttestationObject)
	if err != nil {
		return err
	}

	a.AuthData = a.decodeAuthData(rawAttestationObject["authData"].([]byte))
	a.Format = rawAttestationObject["fmt"].(string)
	a.AttestationRaw = rawAttestationObject["attStmt"].(map[interface{}]interface{})

	return nil
}

func (a *AttestationObject) decodeAuthData(data []byte) AuthData {
	authData := AuthData{}
	authData.RPIDHash = data[:32]
	authData.Flags = data[32]
	authData.SignCount = data[33:37]
	authData.AAGUID = data[37:53]
	credentialIDLength := binary.BigEndian.Uint16(data[53:55])
	authData.CredentialID = data[55 : 55+credentialIDLength]
	authData.CredentialPublicKey = data[55+credentialIDLength:]
	return authData
}

type AuthData struct {
	RPIDHash            []byte
	Flags               byte
	SignCount           []byte
	AAGUID              []byte
	CredentialID        []byte
	CredentialPublicKey []byte
}

type UserService interface {
	IsUserIdAvailable(userId string) bool
	GetUser(userId string) (interface{}, error)
	CreateUser(user interface{}) error
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
		fmt.Println("ERROR", err)
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
