package webauthn

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const rpId = "localhost" // TODO: make customizable

type AuthenticationService struct {
	initAuthenticationStore core.KeyValueStore[string, CredentialOptions]
	userService             *domain.UserService
	credentialService       *domain.CredentialService
}

func NewAuthenticationService(
	initAuthenticationStore core.KeyValueStore[string, CredentialOptions],
	userService *domain.UserService,
	credentialService *domain.CredentialService,
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
				Attestation:      "direct",
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

func (s *AuthenticationService) Register(ctx context.Context, id string, request *CreateCredentialRequest) (*Success, error) {
	options, err := s.initAuthenticationStore.Get(id)
	if err != nil {
		return nil, err
	}

	// NOTE: maybe move this to the authenticator controller
	clientData := &clientData{
		Raw: request.Response.ClientDataJSON,
	}
	err = json.Unmarshal(request.Response.ClientDataJSON, clientData)
	if err != nil {
		return nil, err
	}

	attestationObject, err := request.Response.AttestationObject.Decode()
	if err != nil {
		return nil, err
	}

	response := &CreationCredentialResponse{
		ClientData:        *clientData,
		AttestationObject: *attestationObject,
	}

	credential := &domain.Credential{}

	err = response.Validate(options.GetOptions(), credential)
	if err != nil {
		return nil, err
	}

	userInstance := &domain.User{
		UserID: options.GetUserID(),
		Status: "active",
	}

	err = s.userService.CreateUser(ctx, userInstance)
	if err != nil {
		return nil, err
	}

	credential.User = userInstance

	err = s.credentialService.CreateCredential(ctx, credential)
	if err != nil {
		return nil, err
	}

	result, err := s.IssueGrant(ctx, credential)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type Success struct {
	AccessToken  *domain.AccessToken  `json:"accessToken"`
	RefreshToken *domain.RefreshToken `json:"refreshToken"`
}

func (s *AuthenticationService) Login(ctx context.Context, request *RequestCredentialRequest) (*Success, error) {
	options, err := s.initAuthenticationStore.Get(request.AuthenticationID)
	if err != nil {
		return nil, err
	}

	credential, err := s.credentialService.GetCredentialByCredentialID(ctx, request.RawID)
	if err != nil {
		return nil, err
	}

	// NOTE: maybe move this to the authenticator controller
	clientData := &clientData{
		Raw: request.Response.ClientDataJSON,
	}
	err = json.Unmarshal(request.Response.ClientDataJSON, clientData)
	if err != nil {
		return nil, err
	}

	authenticatorData, err := decodeAuthData(request.Response.AuthenticatorData)
	if err != nil {
		return nil, err
	}

	response := &RequestCredentialResponse{
		ClientData:        *clientData,
		AuthenticatorData: authenticatorData,
		Signature:         request.Response.Signature,
		UserHandle:        request.Response.UserHandle,
	}

	err = response.Validate(options.GetOptions(), credential)
	if err != nil {
		return nil, err
	}

	result, err := s.IssueGrant(ctx, credential)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *AuthenticationService) IssueGrant(ctx context.Context, credential *domain.Credential) (*Success, error) {
	grant := domain.NewGrant(credential.User.ID)
	grant.AllowRefreshToken = true
	grant.ExpiresAt = grant.IssuedAt.Add(time.Hour * 24 * 30)
	grant.NotBefore = grant.IssuedAt
	grant.Scope = []string{"openid", "profile", "email", "authorization"}
	grant.ClientID = "central"
	grant.SubjectID = credential.User.ID
	accessToken, refreshToken, err := domain.RegisterGrant(ctx, grant)
	if err != nil {
		return nil, err
	}

	log.Println("SUCCESS", accessToken, refreshToken)

	return &Success{AccessToken: accessToken, RefreshToken: refreshToken}, nil
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

	result, err := c.service.Register(ctx.Request.Context(), request.AuthenticationID, &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, &result)
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

	result, err := c.service.Login(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	ctx.JSON(200, &result)
}
