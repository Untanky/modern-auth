package webauthn

import (
	"fmt"

	"github.com/Untanky/modern-auth/internal/domain"
	"github.com/Untanky/modern-auth/internal/utils"
)

type CredentialResponse interface {
	Validate(options CredentialOptions, credential *domain.Credential) error
}

type PublicKeyCredentialOptions interface {
	ValidateClientData(clientData clientData) error
	ValidateAttestationObject(attestationObject attestationObject) error
	ValidateAuthenticatorData(authenticatorData AuthData) error
}

type PublicKeyCredentialRequestOptions struct {
	UserId             []byte                          `json:"-"`
	Challenge          []byte                          `json:"challenge"`
	RpID               string                          `json:"rpId"`
	Timeout            uint64                          `json:"timeout"`
	UserVerification   string                          `json:"userVerification"`
	Attestation        string                          `json:"attestation"`
	AttestationFormats []string                        `json:"attestationFormats"`
	AllowCredentials   []PublicKeyCredentialDescriptor `json:"allowCredentials"`
}

func (options *PublicKeyCredentialRequestOptions) ValidateClientData(clientData clientData) error {
	if clientData.Type != "webauthn.get" {
		return fmt.Errorf("invalid type")
	}
	if clientData.Challenge != string(utils.EncodeBase64([]byte(options.Challenge))) {
		return fmt.Errorf("invalid challenge")
	}
	// TODO: fix hardcoding
	if clientData.Origin != "http://localhost:3000" {
		return fmt.Errorf("invalid origin")
	}

	return nil
}

func (options *PublicKeyCredentialRequestOptions) ValidateAttestationObject(attestationObject attestationObject) error {
	return fmt.Errorf("not implemented")
}

func (options *PublicKeyCredentialRequestOptions) ValidateAuthenticatorData(authenticatorData AuthData) error {
	if string(authenticatorData.RPIDHash) != string(utils.HashSHA256([]byte(options.RpID))) {
		return fmt.Errorf("invalid rpIdHash")
	}

	if err := authenticatorData.Flags.Verify(); err != nil {
		return err
	}

	return nil
}

type PublicKeyCredentialDescriptor struct {
	Type       string   `json:"type"`
	ID         []byte   `json:"id"`
	Transports []string `json:"transports"`
}

type RequestCredentialResponse struct {
	ClientData        clientData
	AuthenticatorData AuthData
	Signature         []byte
	UserHandle        []byte
}

func (response *RequestCredentialResponse) Validate(options PublicKeyCredentialOptions, credential *domain.Credential) error {
	err := options.ValidateClientData(response.ClientData)
	if err != nil {
		return err
	}

	clientDataHash := utils.HashSHA256(response.ClientData.Raw)

	err = options.ValidateAuthenticatorData(response.AuthenticatorData)
	if err != nil {
		return err
	}

	publicKey, err := decodeKey(credential.PublicKey)
	if err != nil {
		return err
	}

	verificationData := append(response.AuthenticatorData.Raw, clientDataHash...)
	ok := publicKey.Verify(response.Signature, verificationData)
	if !ok {
		return fmt.Errorf("invalid signature")
	}

	return nil
}
