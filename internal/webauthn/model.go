package webauthn

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"

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
	if data.Challenge != string(utils.EncodeBase64([]byte(options.PublicKeyOptions.Challenge))) {
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

func decodeAuthData(data []byte) AuthData {
	authData := AuthData{}
	authData.Raw = data
	authData.RPIDHash = data[:32]
	authData.Flags = AuthFlags(data[32])
	authData.SignCount = data[33:37]
	authData.AAGUID = data[37:53]
	credentialIDLength := binary.BigEndian.Uint16(data[53:55])
	authData.CredentialID = data[55 : 55+credentialIDLength]
	authData.RawCredentialPublicKey = data[55+credentialIDLength:]
	publicKey, err := decodeKey(authData.RawCredentialPublicKey)
	if err != nil {
		panic(err)
	}
	authData.CredentialPublicKey = publicKey
	// TODO: implement decodePublicKey
	return authData
}

func decodeKey(data []byte) (PublicKey, error) {
	raw := make(map[interface{}]interface{}, 0)
	err := cbor.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	alg := raw[uint64(3)].(int64)
	switch alg {
	case -7:
		return &es256PublicKey{
			key: &ecdsa.PublicKey{
				Curve: elliptic.P256(),
				X:     big.NewInt(0).SetBytes(raw[int64(-2)].([]byte)),
				Y:     big.NewInt(0).SetBytes(raw[int64(-3)].([]byte)),
			},
		}, nil
	default:
		return nil, fmt.Errorf("invalid algorithm")
	}
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
	ClientDataJSON    RawClientDataJSON    `json:"clientDataJSON"`
	AttestationObject RawAttestationObject `json:"attestationObject"`
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
