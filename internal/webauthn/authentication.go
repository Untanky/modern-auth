package webauthn

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/user"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const rpId = "localhost" // TODO: make customizable

type AuthenticationService struct {
	initAuthenticationStore core.KeyValueStore[string, CredentialOptions]
	userService             *user.UserService
	credentialService       *user.CredentialService
}

func NewAuthenticationService(
	initAuthenticationStore core.KeyValueStore[string, CredentialOptions],
	userService *user.UserService,
	credentialService *user.CredentialService,
) *AuthenticationService {
	return &AuthenticationService{
		initAuthenticationStore: initAuthenticationStore,
		userService:             userService,
		credentialService:       credentialService,
	}
}

func (s *AuthenticationService) InitiateAuthentication(request *InitiateAuthenticationRequest) (CredentialOptions, error) {
	id := uuid.New().String()

	user, err := s.userService.GetUserByUserID(context.TODO(), []byte(request.UserId))
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	var initResponse CredentialOptions

	if user == nil {
		initResponse = &CredentialCreationOptions{
			AuthenticationId: id,
			Type:             "create",
			Options: PublicKeyCredentialCreationOptions{
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

		allowCredentials := []PublicKeyCredentialDescriptor{}
		for _, credential := range credentials {
			allowCredentials = append(allowCredentials, PublicKeyCredentialDescriptor{
				Type: "public-key",
				ID:   credential.CredentialID,
			})
		}

		initResponse = &CredentialRequestOptions{
			AuthenticationId: id,
			Type:             "get",
			Options: PublicKeyCredentialRequestOptions{
				UserId: []byte(request.UserId),
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

func (s *AuthenticationService) Register(ctx context.Context, id string, request *CreateCredentialRequest) error {
	options, err := s.initAuthenticationStore.Get(id)
	if err != nil {
		return err
	}

	// NOTE: maybe move this to the authenticator controller
	clientData := &clientData{
		Raw: request.Response.ClientDataJSON,
	}
	err = json.Unmarshal(request.Response.ClientDataJSON, clientData)
	if err != nil {
		return err
	}

	attestationObject, err := request.Response.AttestationObject.Decode()
	if err != nil {
		return err
	}

	response := &CreationCredentialResponse{
		ClientData:        *clientData,
		AttestationObject: *attestationObject,
	}

	credential := &user.Credential{}

	err = response.Validate(options.GetOptions(), credential)
	if err != nil {
		return err
	}

	userInstance := &user.User{
		UserID: options.GetUserID(),
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

func (s *AuthenticationService) Login(ctx context.Context, request *RequestCredentialRequest) error {
	options, err := s.initAuthenticationStore.Get(request.AuthenticationID)
	if err != nil {
		return err
	}

	credential, err := s.credentialService.GetCredentialByCredentialID(ctx, request.RawID)
	if err != nil {
		return err
	}

	// NOTE: maybe move this to the authenticator controller
	clientData := &clientData{
		Raw: request.Response.ClientDataJSON,
	}
	err = json.Unmarshal(request.Response.ClientDataJSON, clientData)
	if err != nil {
		return err
	}

	authenticatorData, err := decodeAuthData(request.Response.AuthenticatorData)
	if err != nil {
		return err
	}

	response := &RequestCredentialResponse{
		ClientData:        *clientData,
		AuthenticatorData: authenticatorData,
		Signature:         request.Response.Signature,
		UserHandle:        request.Response.UserHandle,
	}

	err = response.Validate(options.GetOptions(), credential)
	if err != nil {
		return err
	}

	// TODO: assess trust of the authenticator

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
	router.POST("/authentication/validate", c.getCredential)
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	err = c.service.Register(ctx.Request.Context(), request.AuthenticationID, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}

func (c *AuthenticationController) getCredential(ctx *gin.Context) {
	var request RequestCredentialRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	err = c.service.Login(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success": true,
	})
}
