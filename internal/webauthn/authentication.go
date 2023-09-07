package webauthn

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Untanky/modern-auth/internal/core"
	"github.com/Untanky/modern-auth/internal/domain"
	"github.com/Untanky/modern-auth/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const rpId = "localhost" // TODO: make customizable

type AuthenticationService struct {
	initAuthenticationStore     core.KeyValueStore[string, CredentialOptions]
	authenticationVerifierStore core.KeyValueStore[string, []byte]
	userService                 *domain.UserService
	credentialService           *domain.CredentialService
	logger                      *slog.Logger
}

func NewAuthenticationService(
	initAuthenticationStore core.KeyValueStore[string, CredentialOptions],
	authenticationVerifierStore core.KeyValueStore[string, []byte],
	userService *domain.UserService,
	credentialService *domain.CredentialService,
) *AuthenticationService {
	logger := slog.Default().With("service", "web-authentication")

	return &AuthenticationService{
		initAuthenticationStore:     initAuthenticationStore,
		authenticationVerifierStore: authenticationVerifierStore,
		userService:                 userService,
		credentialService:           credentialService,
		logger:                      logger,
	}
}

func (s *AuthenticationService) InitiateAuthentication(ctx context.Context, request *InitiateAuthenticationRequest) (CredentialOptions, error) {
	id := uuid.New().String()
	userIdBytes := []byte(request.UserId)

	grouped := s.logger.With("userId", utils.EncodeBase64(utils.HashShake256(userIdBytes))).WithGroup("authentication").With("id", id)

	grouped.DebugContext(ctx, "Starting authentication")

	user, err := s.userService.GetUserByUserID(context.TODO(), userIdBytes)
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
					Id:          userIdBytes,
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

		grouped = grouped.With("type", "create")

		grouped.InfoContext(ctx, "Requesting new credential")
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
				UserId: userIdBytes,
				// TODO: randomly generate challenge
				Challenge:        []byte("1234567890"),
				RpID:             rpId,
				UserVerification: "preferred",
				Attestation:      "direct",
				AllowCredentials: allowCredentials,
				Timeout:          60000,
			},
		}

		grouped = grouped.With("type", "get")

		grouped.InfoContext(ctx, "Requesting existing credential")
	}

	err = s.initAuthenticationStore.Set(id, initResponse)
	if err != nil {
		return nil, err
	}

	grouped.InfoContext(ctx, "Initialized authentication")

	return initResponse, nil
}

func (s *AuthenticationService) Register(ctx context.Context, request *CreateCredentialRequest) (*Success, error) {
	grouped := s.logger.WithGroup("authentication").With(slog.String("id", request.AuthenticationID), slog.String("type", "create"))

	grouped.DebugContext(ctx, "Received credential request")

	options, err := s.initAuthenticationStore.Get(request.AuthenticationID)
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

	grouped.DebugContext(ctx, "Parsed credential request")

	credential := &domain.Credential{}

	err = response.Validate(options.GetOptions(), credential)
	if err != nil {
		return nil, err
	}

	grouped.DebugContext(ctx, "Validated credential request")

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

	result, err := s.IssueGrant(ctx, credential.User)
	if err != nil {
		return nil, err
	}

	grouped.InfoContext(ctx, "Registration successful")

	return result, nil
}

type Success struct {
	AccessToken  *domain.AccessToken  `json:"accessToken"`
	RefreshToken *domain.RefreshToken `json:"refreshToken"`
}

func (s *AuthenticationService) Login(ctx context.Context, request *RequestCredentialRequest) (*Success, error) {
	grouped := s.logger.WithGroup("authentication").With(slog.String("id", request.AuthenticationID), slog.String("type", "get"))

	grouped.DebugContext(ctx, "Received credential request")

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

	grouped.DebugContext(ctx, "Parsed credential request")

	err = response.Validate(options.GetOptions(), credential)
	if err != nil {
		return nil, err
	}

	grouped.DebugContext(ctx, "Validated credential request")

	result, err := s.IssueGrant(ctx, credential.User)
	if err != nil {
		return nil, err
	}

	grouped.InfoContext(ctx, "Login success")

	return result, nil
}

func (s *AuthenticationService) IssueGrant(ctx context.Context, user *domain.User) (*Success, error) {
	grant := domain.NewGrant(user.ID)
	grant.AllowRefreshToken = true
	grant.ExpiresAt = grant.IssuedAt.Add(time.Hour * 24 * 30)
	grant.NotBefore = grant.IssuedAt
	grant.Scope = []string{"openid", "profile", "email", "authorization"}
	grant.ClientID = "central"
	grant.SubjectID = user.ID
	accessToken, refreshToken, err := domain.RegisterGrant(ctx, grant)
	if err != nil {
		return nil, err
	}

	s.logger.DebugContext(ctx, "Issued authentication grant", "userUid", user.ID, "grantId", grant.ID)

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

	response, err := c.service.InitiateAuthentication(ctx.Request.Context(), &request)
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

	result, err := c.service.Register(ctx.Request.Context(), &request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request",
		})
		return
	}

	cookie, err := ctx.Cookie("authorization_id")
	if err != nil || cookie == "" {
		ctx.JSON(200, &result)
		return
	}

	authVerifier, err := c.service.continueAuthorization(ctx.Request.Context(), cookie)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal_server_error",
		})
	}
	ctx.SetCookie("authentication_verifier", string(utils.EncodeBase64(authVerifier)), 300, "", "localhost", true, true)

	// TODO: maybe redirect
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

	cookie, err := ctx.Cookie("authorization_id")
	if err != nil || cookie == "" {
		ctx.JSON(200, &result)
		return
	}

	authVerifier, err := c.service.continueAuthorization(ctx.Request.Context(), cookie)
	ctx.SetCookie("authentication_verifier", string(utils.EncodeBase64(authVerifier)), 300, "", "localhost", true, true)

	// TODO: maybe redirect
	ctx.JSON(200, &result)
}

func (s *AuthenticationService) continueAuthorization(ctx context.Context, authorizationId string) ([]byte, error) {
	rand := make([]byte, 64)
	utils.RandomBytes(rand)
	firstHash := utils.HashShake256(rand)
	secondHash := utils.HashShake256(firstHash)

	fmt.Println(firstHash, secondHash)

	err := s.authenticationVerifierStore.Set(authorizationId, secondHash)
	if err != nil {
		return nil, err
	}

	return firstHash, nil
}
