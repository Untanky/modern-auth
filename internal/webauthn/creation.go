package webauthn

import (
	"fmt"

	"github.com/Untanky/modern-auth/internal/user"
	"github.com/Untanky/modern-auth/internal/utils"
)

type PublicKeyCredentialCreationOptions struct {
	Challenge                 []byte                          `json:"challenge"`
	RelyingParty              PublicKeyCredentialRpEntity     `json:"rp"`
	User                      PublicKeyCredentialUserEntity   `json:"user"`
	PublicKeyCredentialParams []PublicKeyCredentialParameters `json:"pubKeyCredParams"`
	AuthenticationSelection   AuthenticationSelection         `json:"authenticatorSelection"`
	Timeout                   uint64                          `json:"timeout"`
	Attestation               string                          `json:"attestation"`
	AttestationFormats        []string                        `json:"attestationFormats"`
}

type PublicKeyCredentialRpEntity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PublicKeyCredentialUserEntity struct {
	Id          []byte `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type PublicKeyCredentialParameters struct {
	Type string `json:"type"`
	Alg  int    `json:"alg"`
}

type AuthenticationSelection struct {
	AuthenticatorAttachment string `json:"authenticatorAttachment"`
	RequireResidentKey      bool   `json:"requireResidentKey"`
	UserVerification        string `json:"userVerification"`
}

func (options *PublicKeyCredentialCreationOptions) ValidateClientData(clientData clientData) error {
	if clientData.Type != "webauthn.create" {
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

func (options *PublicKeyCredentialCreationOptions) ValidateAttestationObject(attestationObject attestationObject) error {
	if attestationObject.Format != "packed" {
		return fmt.Errorf("invalid attestation format")
	}

	if err := options.ValidateAuthenticatorData(attestationObject.AuthData); err != nil {
		return err
	}

	return nil
}

func (options *PublicKeyCredentialCreationOptions) ValidateAuthenticatorData(authenticatorData AuthData) error {
	// this assumes there is only one relying party
	if string(authenticatorData.RPIDHash) != string(utils.HashSHA256([]byte(options.RelyingParty.Id))) {
		return fmt.Errorf("invalid rpIdHash")
	}

	if err := authenticatorData.Flags.Verify(); err != nil {
		return err
	}

	found := false
	for _, param := range options.PublicKeyCredentialParams {
		if param.Alg == authenticatorData.CredentialPublicKey.Algorithm() {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("invalid algorithm")
	}

	return nil
}

type CreationCredentialResponse struct {
	ClientData        clientData
	AttestationObject attestationObject
}

func (response *CreationCredentialResponse) Validate(options PublicKeyCredentialOptions, credential *user.Credential) error {
	err := options.ValidateClientData(response.ClientData)
	if err != nil {
		return err
	}

	clientDataHash := utils.HashSHA256(response.ClientData.Raw)

	err = options.ValidateAttestationObject(response.AttestationObject)
	if err != nil {
		return err
	}

	err = response.AttestationObject.Attestation.Verify(response.AttestationObject.AuthData, clientDataHash)
	if err != nil {
		return err
	}

	credential.CredentialID = response.AttestationObject.AuthData.CredentialID
	credential.PublicKey = response.AttestationObject.AuthData.RawCredentialPublicKey

	return nil
}
