package webauthn_test

import (
	"reflect"
	"testing"

	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/Untanky/modern-auth/internal/webauthn"
)

func TestClientDataJSONValidateCreate(t *testing.T) {
	tests := []struct {
		name    string
		json    webauthn.ClientDataJSON
		options *webauthn.InitiateAuthenticationResponse
		want    []byte
		wantErr bool
	}{
		{
			name: "Succeeds",
			json: webauthn.ClientDataJSON(`{"type":"webauthn.create","challenge":"1234567890","origin":"localhost"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialRequestOptions{
					Challenge: []byte("1234567890"),
					RelyingParty: webauthn.RelyingPartyOptions{
						Id: "localhost",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fails when type is not webauthn.create",
			json: webauthn.ClientDataJSON(`{"type":"webauthn.get","challenge":"1234567890","origin":"localhost"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialRequestOptions{
					Challenge: []byte("1234567890"),
					RelyingParty: webauthn.RelyingPartyOptions{
						Id: "localhost",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fails when type is not challenge does not match",
			json: webauthn.ClientDataJSON(`{"type":"webauthn.create","challenge":"abc","origin":"localhost"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialRequestOptions{
					Challenge: []byte("1234567890"),
					RelyingParty: webauthn.RelyingPartyOptions{
						Id: "localhost",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fails when type is not origin does not match",
			json: webauthn.ClientDataJSON(`{"type":"webauthn.create","challenge":"1234567890","origin":"modern-auth.com"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialRequestOptions{
					Challenge: []byte("1234567890"),
					RelyingParty: webauthn.RelyingPartyOptions{
						Id: "localhost",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := tt.json.ValidateCreate(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("ClientDataJSON.ValidateCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				want := utils.HashSHA256([]byte(tt.json))
				if !reflect.DeepEqual(hash, want) {
					t.Errorf("ClientDataJSON.ValidateCreate() = %v, want %v", hash, tt.want)
				}
			}
		})
	}
}
