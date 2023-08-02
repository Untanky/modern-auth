package webauthn

type CredentialOptions interface {
	GetAuthenticationID() string
	IsCreationOptions() bool
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

func (options *CredentialCreationOptions) GetAuthenticationID() string {
	return options.AuthenticationId
}

func (options *CredentialCreationOptions) IsCreationOptions() bool {
	return true
}

type CredentialRequestOptions struct {
	AuthenticationId string                            `json:"authenticationId"`
	Type             string                            `json:"type"`
	Options          PublicKeyCredentialRequestOptions `json:"publicKey"`
}

func (options *CredentialRequestOptions) GetAuthenticationID() string {
	return options.AuthenticationId
}

func (options *CredentialRequestOptions) IsCreationOptions() bool {
	return false
}

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

type PublicKeyCredentialRequestOptions struct {
	Challenge          []byte                          `json:"challenge"`
	RpID               string                          `json:"rpId"`
	Timeout            uint64                          `json:"timeout"`
	UserVerification   string                          `json:"userVerification"`
	Attestation        string                          `json:"attestation"`
	AttestationFormats []string                        `json:"attestationFormats"`
	AllowCredentials   []PublicKeyCredentialDescriptor `json:"allowCredentials"`
}

type PublicKeyCredentialDescriptor struct {
	Type       string   `json:"type"`
	ID         []byte   `json:"id"`
	Transports []string `json:"transports"`
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

type CreateCredentialRequest struct {
	AuthenticationID string                   `json:"authenticationId"`
	Id               string                   `json:"id"`
	RawID            []byte                   `json:"rawId"`
	Type             string                   `json:"type"`
	Response         CreateCredentialResponse `json:"response"`
}

type RequestCredentialRequest struct {
	AuthenticationID string                    `json:"authenticationId"`
	Id               string                    `json:"id"`
	RawID            []byte                    `json:"rawId"`
	Type             string                    `json:"type"`
	Response         RequestCredentialResponse `json:"response"`
}
