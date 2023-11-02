package api

import (
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type server struct {
	common.CommonServer
	clients     clientProvider
	appIdentity common.AzureAppConfidentialIdentity
}

// ConfidentialAppIdentity implements common.ConfidentialAppIdentityProvider.
func (s *server) ConfidentialAppIdentity() common.AzureAppConfidentialIdentity {
	return s.appIdentity
}

type H = map[string]string
type RequestContext = ctx.RequestContext

func wrapResponse[T interface{}](c echo.Context, defaultStatus int, data T, err error) error {
	switch {
	case err == nil:
		return c.JSON(defaultStatus, data)
	case errors.Is(err, common.ErrStatus2xxCreated):
		return c.JSON(http.StatusCreated, data)
	case errors.Is(err, common.ErrStatusBadRequest):
		return c.JSON(http.StatusBadRequest, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusUnauthorized):
		return c.JSON(http.StatusUnauthorized, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusForbidden):
		return c.JSON(http.StatusForbidden, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusNotFound):
		return c.JSON(http.StatusNotFound, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusConflict):
		return c.JSON(http.StatusNotFound, H{"error": err.Error()})
	default:
		c.Logger().Error("internal error", err)
		return c.JSON(http.StatusInternalServerError, H{"error": "internal error"})
	}
}

func respondRequireAdmin(c echo.Context) error {
	return c.JSON(http.StatusForbidden, map[string]string{"message": "admin access required"})
}

func wrapEchoResponse(c echo.Context, err error) error {
	if err == nil {
		return err
	}
	switch {
	case errors.Is(err, common.ErrStatusBadRequest):
		return c.JSON(http.StatusBadRequest, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusUnauthorized):
		return c.JSON(http.StatusUnauthorized, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusForbidden):
		return c.JSON(http.StatusForbidden, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusNotFound):
		return c.JSON(http.StatusNotFound, H{"error": err.Error()})
	case errors.Is(err, common.ErrStatusConflict):
		return c.JSON(http.StatusNotFound, H{"error": err.Error()})
	}
	log.Error().Err(err).Msg("internal error")
	return c.JSON(http.StatusInternalServerError, H{"error": "internal error"})
}

type appConfidentialIdentity struct {
	tenantID               string
	clientID               string
	clientSecret           string
	clientSecretCredential *azidentity.ClientSecretCredential
}

// ClientID implements auth.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) ClientID() string {
	return i.clientID
}

// TenantID implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) TenantID() string {
	return i.tenantID
}

// GetOnBehalfOfTokenCredential implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) NewOnBehalfOfTokenCredential(userAssertion string, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error) {
	return azidentity.NewOnBehalfOfCredentialWithSecret(i.tenantID, i.clientID, userAssertion, i.clientSecret, opts)
}

// TokenCredential implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) TokenCredential() azcore.TokenCredential {
	return i.clientSecretCredential
}

var _ common.AzureAppConfidentialIdentity = (*appConfidentialIdentity)(nil)

func NewServer() *server {

	commonConfig, err := common.NewCommonConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to create common server")
	}

	s := server{
		CommonServer: &commonConfig,
	}

	appId := appConfidentialIdentity{}
	if appId.tenantID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzTenantID, ""); appId.tenantID == "" {
		log.Panic().Msg("No app tenant ID found in environment variable")
	}
	if appId.clientID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzClientID, ""); appId.clientID == "" {
		log.Panic().Msg("No app client ID found in environment variable")
	}
	if appId.clientSecret = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzClientSecret, ""); appId.clientSecret == "" {
		log.Panic().Msg("No app client secret found in environment variable")
	}
	if appId.clientSecretCredential, err = azidentity.NewClientSecretCredential(
		appId.tenantID, appId.clientID, appId.clientSecret, nil); err != nil {
		log.Panic().Err(err).Msg("Failed to create app client secret credential")
	}
	s.appIdentity = &appId

	if s.clients, err = newServerClientProvider(&s); err != nil {
		log.Panic().Err(err).Msg("failed to create client provider")
	}

	return &s
}

func (s *server) GetAfterAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rc, ok := c.(RequestContext); ok {
				rc = rc.WithValue(common.AdminServerRequestClientProvierContextKey, &requestClientProvider{
					parent:            s,
					credentialContext: rc,
				})
				return next(rc)
			}
			return next(c)
		}
	}
}

var _ common.ConfidentialAppIdentityProvider = (*server)(nil)
