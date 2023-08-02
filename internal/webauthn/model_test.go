package webauthn_test

import (
	"reflect"
	"testing"

	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/Untanky/modern-auth/internal/webauthn"
)

func TestRawClientDataJSONVerifyCreate(t *testing.T) {
	tests := []struct {
		name    string
		json    webauthn.RawClientDataJSON
		options *webauthn.InitiateAuthenticationResponse
		want    []byte
		wantErr bool
	}{
		{
			name: "Succeeds",
			json: webauthn.RawClientDataJSON(`{"type":"webauthn.create","challenge":"1234567890","origin":"localhost"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialCreationOptions{
					Challenge: "1234567890",
					RelyingParty: webauthn.PublicKeyCredentialRpEntity{
						Id: "localhost",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fails when type is not webauthn.create",
			json: webauthn.RawClientDataJSON(`{"type":"webauthn.get","challenge":"1234567890","origin":"localhost"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialCreationOptions{
					Challenge: "1234567890",
					RelyingParty: webauthn.PublicKeyCredentialRpEntity{
						Id: "localhost",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fails when type is not challenge does not match",
			json: webauthn.RawClientDataJSON(`{"type":"webauthn.create","challenge":"abc","origin":"localhost"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialCreationOptions{
					Challenge: "1234567890",
					RelyingParty: webauthn.PublicKeyCredentialRpEntity{
						Id: "localhost",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fails when type is not origin does not match",
			json: webauthn.RawClientDataJSON(`{"type":"webauthn.create","challenge":"1234567890","origin":"modern-auth.com"}`),
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialCreationOptions{
					Challenge: "1234567890",
					RelyingParty: webauthn.PublicKeyCredentialRpEntity{
						Id: "localhost",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := tt.json.VerifyCreate(tt.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawClientDataJSON.VerifyCreate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				want := utils.HashSHA256([]byte(tt.json))
				if !reflect.DeepEqual(hash, want) {
					t.Errorf("RawClientDataJSON.VerifyCreate() = %v, want %v", hash, tt.want)
				}
			}
		})
	}
}

func TestAuthDataVerify(t *testing.T) {
	tests := []struct {
		name    string
		auth    webauthn.AuthData
		options *webauthn.InitiateAuthenticationResponse
		wantErr bool
	}{
		{
			name: "Succeeds",
			auth: webauthn.AuthData{
				RPIDHash: utils.HashSHA256([]byte("localhost")),
			},
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialCreationOptions{
					RelyingParty: webauthn.PublicKeyCredentialRpEntity{
						Id: "localhost",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Succeeds",
			auth: webauthn.AuthData{
				RPIDHash: []byte("1234567890"),
			},
			options: &webauthn.InitiateAuthenticationResponse{
				PublicKeyOptions: webauthn.PublicKeyCredentialCreationOptions{
					RelyingParty: webauthn.PublicKeyCredentialRpEntity{
						Id: "localhost",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.auth.Verify(tt.options); (err != nil) != tt.wantErr {
				t.Errorf("AuthData.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
