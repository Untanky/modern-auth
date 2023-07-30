package webauthn

import (
	"fmt"

	"encoding/binary"
	jsonlib "encoding/json"

	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/fxamacker/cbor/v2"
)

type RawClientDataJSON []byte

type clientData struct {
	Origin    string `json:"origin"`
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
}

func (json RawClientDataJSON) VerifyCreate(options *InitiateAuthenticationResponse) (hash []byte, err error) {
	var data clientData
	err = jsonlib.Unmarshal(json, &data)
	if err != nil {
		return nil, err
	}

	if data.Type != "webauthn.create" {
		return nil, fmt.Errorf("invalid type")
	}
	if data.Challenge != string(options.PublicKeyOptions.Challenge) {
		return nil, fmt.Errorf("invalid challenge")
	}
	if data.Origin != options.PublicKeyOptions.RelyingParty.Id {
		return nil, fmt.Errorf("invalid origin")
	}

	return utils.HashSHA256(json), nil
}

type RawAttestationObject []byte

type AttestationStatement interface {
	Verify(authData AuthData, clienDataHash []byte) error
}

type attestationObject struct {
	AuthData       AuthData
	Format         string
	AttestationRaw map[interface{}]interface{}
	Attestation    AttestationStatement
}

func (attestation attestationObject) Verify(options *InitiateAuthenticationResponse, clientDataHash []byte) error {
	if attestation.Format != "packed" {
		return fmt.Errorf("invalid attestation format")
	}

	if err := attestation.AuthData.Verify(options); err != nil {
		return err
	}

	if err := attestation.Attestation.Verify(attestation.AuthData, clientDataHash); err != nil {
		return err
	}

	return nil
}

func (attestation RawAttestationObject) Decode() (*attestationObject, error) {
	var rawAttestationObject map[string]interface{}
	err := cbor.Unmarshal(attestation, &rawAttestationObject)
	if err != nil {
		return nil, err
	}

	var attestationObject attestationObject
	attestationObject.AuthData = decodeAuthData(rawAttestationObject["authData"].([]byte))
	attestationObject.Format = rawAttestationObject["fmt"].(string)
	attestationObject.AttestationRaw = rawAttestationObject["attStmt"].(map[interface{}]interface{})

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

type encodedPublicKey []byte

func (key encodedPublicKey) Algorithm() int {
	// TODO: implement correctly
	return -7
}

func (key encodedPublicKey) Verify(signature []byte, value []byte) bool {
	// TODO: Implement corrected
	return true
}

type AuthData struct {
	RPIDHash               []byte
	Flags                  AuthFlags
	SignCount              []byte
	AAGUID                 []byte
	CredentialID           []byte
	RawCredentialPublicKey []byte
	CredentialPublicKey    PublicKey
}

func decodeAuthData(data []byte) AuthData {
	authData := AuthData{}
	authData.RPIDHash = data[:32]
	authData.Flags = AuthFlags(data[32])
	authData.SignCount = data[33:37]
	authData.AAGUID = data[37:53]
	credentialIDLength := binary.BigEndian.Uint16(data[53:55])
	authData.CredentialID = data[55 : 55+credentialIDLength]
	authData.RawCredentialPublicKey = data[55+credentialIDLength:]
	// TODO: implement decodePublicKey
	return authData
}

func (authData AuthData) Verify(options *InitiateAuthenticationResponse) error {
	if string(authData.RPIDHash) != string(utils.HashSHA256([]byte(options.PublicKeyOptions.RelyingParty.Id))) {
		return fmt.Errorf("invalid rpIdHash")
	}

	if err := authData.Flags.Verify(); err != nil {
		return err
	}

	found := false
	for _, param := range options.PublicKeyOptions.PublicKeyCredentialParams {
		if param.Alg == authData.CredentialPublicKey.Algorithm() {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid algorithm")
	}

	return nil
}

type CreateCredentialResponse struct {
	ClientDataJSON    RawClientDataJSON
	AttestationObject RawAttestationObject
}

func (response *CreateCredentialResponse) Validate(options *InitiateAuthenticationResponse) error {
	clientDataHash, err := response.ClientDataJSON.VerifyCreate(options)
	if err != nil {
		return err
	}

	attestationObject, err := response.AttestationObject.Decode()
	if err != nil {
		return err
	}

	err = attestationObject.Verify(options, clientDataHash)
	if err != nil {
		return err
	}

	return nil
}
