package webauthn

import (
	"fmt"

	jsonlib "encoding/json"

	"github.com/Untanky/modern-auth/internal/utils"
)

type ClientDataJSON []byte

type clientData struct {
	Origin    string `json:"origin"`
	Type      string `json:"type"`
	Challenge string `json:"challenge"`
}

func (json ClientDataJSON) ValidateCreate(options *InitiateAuthenticationResponse) (hash []byte, err error) {
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

type CreateCredentialResponse struct {
	ClientDataJSON ClientDataJSON `json:"clientDataJSON"`
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

func (s *AuthenticationService) Register(request *CreateCredentialRequest) error {
	if request.Response.ClientDataJSON.Type != "webauthn.create" {
		return fmt.Errorf("invalid type")
	}

	// TODO: verify against randomly generated challenge
	if request.Response.ClientDataJSON.Challenge != "1234567890" {
		return fmt.Errorf("invalid challenge")
	}

	if request.Response.ClientDataJSON.Origin != rpId {
		return fmt.Errorf("invalid origin")
	}

	if string(request.Response.AttestationObject.AuthData.RPIDHash) == string(utils.HashSHA256([]byte(rpId))) {
		return fmt.Errorf("invalid rpIdHash")
	}

	// TODO: verify that the UP bit of the flags is set

	if request.Response.AttestationObject.Format != "packed" {
		return fmt.Errorf("unsupported format")
	}

	if request.Response.AttestationObject.AttestationRaw["alg"] != -7 {
		return fmt.Errorf("unsupported algorithm")
	}

	return nil
}
