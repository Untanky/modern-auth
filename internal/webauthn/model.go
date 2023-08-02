package webauthn

import (
	"fmt"
	"log"

	"encoding/binary"
	jsonlib "encoding/json"

	"github.com/Untanky/modern-auth/internal/user"
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
	if data.Challenge != string(utils.EncodeBase64([]byte(options.CreationOptions.Challenge))) {
		return nil, fmt.Errorf("invalid challenge")
	}
	// TODO: fix hardcoding
	if data.Origin != "http://localhost:3000" {
		return nil, fmt.Errorf("invalid origin")
	}

	return utils.HashSHA256(json), nil
}

func (json RawClientDataJSON) VerifyGet(options *InitiateAuthenticationResponse) (hash []byte, err error) {
	var data clientData
	err = jsonlib.Unmarshal(json, &data)
	if err != nil {
		return nil, err
	}

	if data.Type != "webauthn.get" {
		return nil, fmt.Errorf("invalid type")
	}
	log.Println(data.Challenge, string(utils.EncodeBase64([]byte(options.RequestOptions.Challenge))))
	if data.Challenge != string(utils.EncodeBase64([]byte(options.RequestOptions.Challenge))) {
		return nil, fmt.Errorf("invalid challenge")
	}
	// TODO: fix hardcoding
	if data.Origin != "http://localhost:3000" {
		return nil, fmt.Errorf("invalid origin")
	}

	return utils.HashSHA256(json), nil
}

type RawAttestationObject []byte

type AttestationStatement interface {
	Verify(authData AuthData, clienDataHash []byte) error
}

type attestationObject struct {
	AuthData    AuthData
	Format      string
	Attestation AttestationStatement
}

func (attestation attestationObject) Verify(options *InitiateAuthenticationResponse, clientDataHash []byte) error {
	if attestation.Format != "packed" {
		return fmt.Errorf("invalid attestation format")
	}

	if err := attestation.AuthData.VerifyCreate(options); err != nil {
		return err
	}

	if err := attestation.Attestation.Verify(attestation.AuthData, clientDataHash); err != nil {
		return err
	}

	return nil
}

func (attestation RawAttestationObject) Decode() (*attestationObject, error) {
	var rawAttestationObject map[string]interface{}
	log.Println(string(attestation))
	err := cbor.Unmarshal(attestation, &rawAttestationObject)
	if err != nil {
		log.Println("Unmarshal", err)
		return nil, err
	}

	log.Println(rawAttestationObject)

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

func (authData AuthData) VerifyCreate(options *InitiateAuthenticationResponse) error {
	// this assumes there is only one relying party
	if string(authData.RPIDHash) != string(utils.HashSHA256([]byte(options.CreationOptions.RelyingParty.Id))) {
		return fmt.Errorf("invalid rpIdHash")
	}

	if err := authData.Flags.Verify(); err != nil {
		return err
	}

	found := false
	for _, param := range options.CreationOptions.PublicKeyCredentialParams {
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

func (authData AuthData) VerifyGet(options *InitiateAuthenticationResponse) error {
	if string(authData.RPIDHash) != string(utils.HashSHA256([]byte(options.RequestOptions.RpID))) {
		return fmt.Errorf("invalid rpIdHash")
	}

	if err := authData.Flags.Verify(); err != nil {
		return err
	}

	return nil
}

type CreateCredentialResponse struct {
	ClientDataJSON    RawClientDataJSON    `json:"clientDataJSON"`
	AttestationObject RawAttestationObject `json:"attestationObject"`
}

func (response *CreateCredentialResponse) Validate(options *InitiateAuthenticationResponse) (*user.Credential, error) {
	clientDataHash, err := response.ClientDataJSON.VerifyCreate(options)
	if err != nil {
		return nil, err
	}

	attestationObject, err := response.AttestationObject.Decode()
	if err != nil {
		return nil, err
	}

	err = attestationObject.Verify(options, clientDataHash)
	if err != nil {
		return nil, err
	}

	return &user.Credential{
		CredentialID: attestationObject.AuthData.CredentialID,
		PublicKey:    attestationObject.AuthData.RawCredentialPublicKey,
	}, nil
}

type RequestCredentialResponse struct {
	ClientDataJSON    RawClientDataJSON `json:"clientDataJSON"`
	AuthenticatorData []byte            `json:"authenticatorData"`
	Signature         []byte            `json:"signature"`
	UserHandle        []byte            `json:"userHandle"`
}

func (response *RequestCredentialResponse) Validate(options *InitiateAuthenticationResponse, credential *user.Credential) error {
	clientDataHash, err := response.ClientDataJSON.VerifyGet(options)
	if err != nil {
		log.Println("VerifyGet", err)
		return err
	}

	authData, err := decodeAuthData(response.AuthenticatorData)
	if err != nil {
		log.Println("Decode", err)
		return err
	}

	err = authData.VerifyGet(options)
	if err != nil {
		log.Println("Verify", err)
		return err
	}

	publicKey, err := decodeKey(credential.PublicKey)
	if err != nil {
		log.Println("decodeKey", err)
		return err
	}

	verificationData := append(authData.Raw, clientDataHash...)
	ok := publicKey.Verify(response.Signature, verificationData)

	if !ok {
		log.Println("Signature", err)
		return fmt.Errorf("invalid signature")
	}

	return nil
}
