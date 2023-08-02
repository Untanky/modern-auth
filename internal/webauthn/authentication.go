package webauthn

import (
	"context"
	"log"
	"net/http"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/user"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const rpId = "localhost" // TODO: make customizable

type InitiateAuthenticationRequest struct {
	UserId string `json:"userId"`
}

type InitiateAuthenticationResponse struct {
	OptionId        string                             `json:"optionId"`
	CreationOptions PublicKeyCredentialCreationOptions `json:"publicKey"`
	RequestOptions  PublicKeyCredentialRequestOptions  `json:"publicKeyFoo"`
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
	OptionId string                   `json:"optionId"`
	Id       string                   `json:"id"`
	Type     string                   `json:"type"`
	Response CreateCredentialResponse `json:"response"`
}

type UserService interface {
	GetUserByUserID(ctx context.Context, userId []byte) (*user.User, error)
	CreateUser(ctx context.Context, user *user.User) error
}

type CredentialService interface {
	GetCredentialByCredentialID(ctx context.Context, creadentialId []byte) (*user.Credential, error)
	GetCredentialsByUserID(ctx context.Context, userId uuid.UUID) ([]*user.Credential, error)
	CreateCredential(ctx context.Context, credential *user.Credential) error
}

type AuthenticationService struct {
	initAuthenticationStore core.KeyValueStore[string, InitiateAuthenticationResponse]
	userService             UserService
	credentialService       CredentialService
}

func NewAuthenticationService(initAuthenticationStore core.KeyValueStore[string, InitiateAuthenticationResponse], userService UserService, credentialService CredentialService) *AuthenticationService {
	return &AuthenticationService{
		initAuthenticationStore: initAuthenticationStore,
		userService:             userService,
		credentialService:       credentialService,
	}
}

func (s *AuthenticationService) InitiateAuthentication(request *InitiateAuthenticationRequest) (*InitiateAuthenticationResponse, error) {
	id := uuid.New().String()

	user, err := s.userService.GetUserByUserID(context.TODO(), []byte(request.UserId))
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var initResponse *InitiateAuthenticationResponse

	if user == nil {
		initResponse = &InitiateAuthenticationResponse{
			OptionId: id,
			CreationOptions: PublicKeyCredentialCreationOptions{
				// TODO: randomly generate challenge
				Challenge: []byte("1234567890"),
				RelyingParty: PublicKeyCredentialRpEntity{
					Id:   rpId,
					Name: "Modern Auth",
				},
				User: PublicKeyCredentialUserEntity{
					Id:          []byte(request.UserId),
					Name:        request.UserId,
					DisplayName: request.UserId,
				},
				PublicKeyCredentialParams: []PublicKeyCredentialParameters{
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
	} else {
		credentials, err := s.credentialService.GetCredentialsByUserID(context.TODO(), user.ID)
		if err != nil {
			return nil, err
		}
		log.Println(credentials)

		allowCredentials := []PublicKeyCredentialDescriptor{}
		for _, credential := range credentials {
			allowCredentials = append(allowCredentials, PublicKeyCredentialDescriptor{
				Type: "public-key",
				ID:   credential.CredentialID,
			})
		}

		initResponse = &InitiateAuthenticationResponse{
			OptionId: id,
			RequestOptions: PublicKeyCredentialRequestOptions{
				// TODO: randomly generate challenge
				Challenge:        []byte("1234567890"),
				RpID:             rpId,
				UserVerification: "preferred",
				Attestation:      "indirect",
				AllowCredentials: allowCredentials,
				Timeout:          60000,
			},
		}
	}

	err = s.initAuthenticationStore.Set(id, initResponse)
	if err != nil {
		return nil, err
	}

	return initResponse, nil
}

func (s *AuthenticationService) Register(ctx context.Context, id string, response *CreateCredentialResponse) error {
	options, err := s.initAuthenticationStore.Get(id)
	if err != nil {
		return err
	}

	credential, err := response.Validate(options)
	if err != nil {
		return err
	}

	// TODO: assess trust of the authenticator

	userInstance := &user.User{
		UserID: options.CreationOptions.User.Id,
		Status: "active",
	}

	err = s.userService.CreateUser(ctx, userInstance)
	if err != nil {
		return err
	}

	credential.User = userInstance

	err = s.credentialService.CreateCredential(ctx, credential)
	if err != nil {
		return err
	}

	return nil
}

type AuthenticationController struct {
	service *AuthenticationService
}

func NewAuthenticationController(service *AuthenticationService) *AuthenticationController {
	return &AuthenticationController{
		service: service,
	}
}

func (c *AuthenticationController) RegisterRoutes(router gin.IRoutes) {
	router.POST("/authentication/initiate", c.initiateAuthentication)
	router.POST("/authentication/create", c.createCredential)
}

func (c *AuthenticationController) initiateAuthentication(ctx *gin.Context) {
	var request InitiateAuthenticationRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	response, err := c.service.InitiateAuthentication(&request)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, response)
}

func (c *AuthenticationController) createCredential(ctx *gin.Context) {
	var request CreateCredentialRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	err = c.service.Register(ctx.Request.Context(), request.OptionId, &request.Response)
	if err != nil {
		log.Printf("ERROR: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}
