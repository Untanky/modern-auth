package webauthn

import (
	"encoding/binary"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

type CredentialOptions interface {
	GetUserID() []byte
	GetAuthenticationID() string
	IsCreationOptions() bool
	GetOptions() PublicKeyCredentialOptions
}

// Request object to initiate an authentication flow,
type InitiateAuthenticationRequest struct {
	// The unique identifier selected by the user
	//
	// Never print this value in plain text
	UserId string `json:"userId"`
}

type CredentialCreationOptions struct {
	AuthenticationId string                             `json:"authenticationId"`
	Type             string                             `json:"type"`
	Options          PublicKeyCredentialCreationOptions `json:"publicKey"`
}

func (options *CredentialCreationOptions) GetUserID() []byte {
	return options.Options.User.Id
}

func (options *CredentialCreationOptions) GetAuthenticationID() string {
	return options.AuthenticationId
}

func (options *CredentialCreationOptions) IsCreationOptions() bool {
	return true
}

func (options *CredentialCreationOptions) GetOptions() PublicKeyCredentialOptions {
	return &options.Options
}

type CredentialRequestOptions struct {
	AuthenticationId string                            `json:"authenticationId"`
	Type             string                            `json:"type"`
	Options          PublicKeyCredentialRequestOptions `json:"publicKey"`
}

func (options *CredentialRequestOptions) GetUserID() []byte {
	return options.Options.UserId
}

func (options *CredentialRequestOptions) GetAuthenticationID() string {
	return options.AuthenticationId
}

func (options *CredentialRequestOptions) IsCreationOptions() bool {
	return false
}

func (options *CredentialRequestOptions) GetOptions() PublicKeyCredentialOptions {
	return &options.Options
}

type CreateCredentialRequest struct {
	AuthenticationID string                      `json:"authenticationId"`
	Id               string                      `json:"id"`
	RawID            []byte                      `json:"rawId"`
	Type             string                      `json:"type"`
	Response         RawCreateCredentialResponse `json:"response"`
}

type RequestCredentialRequest struct {
	AuthenticationID string                       `json:"authenticationId"`
	Id               string                       `json:"id"`
	RawID            []byte                       `json:"rawId"`
	Type             string                       `json:"type"`
	Response         RawRequestCredentialResponse `json:"response"`
}

type RawClientDataJSON []byte

type clientData struct {
	Raw       []byte `json:"-"`
	Origin    string `json:"origin"`
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
}

type RawAttestationObject []byte

type attestationObject struct {
	AuthData    AuthData
	Format      string
	Attestation AttestationStatement
}

type AttestationStatement interface {
	Verify(authData AuthData, clienDataHash []byte) error
}

func (attestation RawAttestationObject) Decode() (*attestationObject, error) {
	var rawAttestationObject map[string]interface{}
	err := cbor.Unmarshal(attestation, &rawAttestationObject)
	if err != nil {
		return nil, err
	}

	var attestationObject attestationObject
	attestationObject.AuthData, err = decodeAuthData(rawAttestationObject["authData"].([]byte))
	if err != nil {
		return nil, err
	}
	attestationObject.Format = rawAttestationObject["fmt"].(string)

	attestationStatement := rawAttestationObject["attStmt"].(map[interface{}]interface{})
	switch attestationObject.Format {
	case "packed":
		certificates := [][]byte{}
		chain := attestationStatement["x5c"]
		if chain != nil {
			for _, cert := range chain.([]interface{}) {
				certificates = append(certificates, cert.([]byte))
			}
		}
		attestationObject.Attestation = &packedAttestationStatemment{
			algorithm:        int(attestationStatement["alg"].(int64)),
			signature:        attestationStatement["sig"].([]byte),
			certificateChain: certificates,
		}
	default:
		return nil, fmt.Errorf("invalid attestation format")
	}

	return &attestationObject, nil
}

type AuthFlags byte

func (flags AuthFlags) Verify() error {
	// TODO: implement
	return nil
}

type PublicKey interface {
	Algorithm() int
	Verify(signature []byte, value []byte) bool
}

type AuthData struct {
	Raw                    []byte
	RPIDHash               []byte
	Flags                  AuthFlags
	SignCount              []byte
	AAGUID                 []byte
	CredentialID           []byte
	RawCredentialPublicKey []byte
	CredentialPublicKey    PublicKey
}

func decodeAuthData(data []byte) (AuthData, error) {
	authData := AuthData{}
	authData.Raw = data
	authData.RPIDHash = data[:32]
	authData.Flags = AuthFlags(data[32])
	authData.SignCount = data[33:37]
	if (len(data)) == 37 {
		return authData, nil
	}

	authData.AAGUID = data[37:53]
	credentialIDLength := binary.BigEndian.Uint16(data[53:55])
	authData.CredentialID = data[55 : 55+credentialIDLength]
	authData.RawCredentialPublicKey = data[55+credentialIDLength:]
	publicKey, err := decodeKey(authData.RawCredentialPublicKey)
	if err != nil {
		return authData, err
	}
	authData.CredentialPublicKey = publicKey
	return authData, nil
}

type RawCreateCredentialResponse struct {
	ClientDataJSON    RawClientDataJSON    `json:"clientDataJSON"`
	AttestationObject RawAttestationObject `json:"attestationObject"`
}

type RawRequestCredentialResponse struct {
	ClientDataJSON    RawClientDataJSON `json:"clientDataJSON"`
	AuthenticatorData []byte            `json:"authenticatorData"`
	Signature         []byte            `json:"signature"`
	UserHandle        []byte            `json:"userHandle"`
}
