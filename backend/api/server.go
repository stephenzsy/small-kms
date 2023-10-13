package api

import (
	ctx "context"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type server struct {
	common.CommonServer
	serverContext     ctx.Context
	clients           clientProvider
	appIdentity       common.AzureAppConfidentialIdentity
	subscriptionId    string
	resourceGroupName string
}

// ConfidentialAppIdentity implements common.ConfidentialAppIdentityProvider.
func (s *server) ConfidentialAppIdentity() common.AzureAppConfidentialIdentity {
	return s.appIdentity
}

// Deadline implements common.ServerContext.
func (s *server) Deadline() (deadline time.Time, ok bool) {
	return s.serverContext.Deadline()
}

// Done implements common.ServerContext.
func (s *server) Done() <-chan struct{} {
	return s.serverContext.Done()
}

// Err implements common.ServerContext.
func (s *server) Err() error {
	return s.serverContext.Err()
}

// Value implements common.ServerContext.
func (s *server) Value(key any) any {
	return s.serverContext.Value(key)
}

type H = map[string]string
type RequestContext = common.RequestContext

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

// GetOnBehalfOfTokenCredential implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) GetOnBehalfOfTokenCredential(userAssertion string, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error) {
	return azidentity.NewOnBehalfOfCredentialWithSecret(i.tenantID, i.clientID, userAssertion, i.clientSecret, opts)
}

// TokenCredential implements common.AzureAppConfidentialIdentity.
func (i *appConfidentialIdentity) TokenCredential() azcore.TokenCredential {
	return i.clientSecretCredential
}

var _ common.AzureAppConfidentialIdentity = (*appConfidentialIdentity)(nil)
var _ models.ServerInterface = (*server)(nil)

func NewServer(c ctx.Context) *server {

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
	s.serverContext = common.WithAdminServerClientProvider(c, &s.clients)

	s.subscriptionId = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, common.IdentityEnvVarNameAzSubscriptionID, "")
	s.resourceGroupName = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, common.IdentityEnvVarNameAzResourceGroupName, "")

	return &s
}

func (s *server) GetPreAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(common.WrapEchoContext(c, s.serverContext))

		}
	}
}

func (s *server) GetAfterAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if rc, ok := c.(RequestContext); ok {
				rc = common.WithAdminServerRequestClientProvider(rc, &requestClientProvider{
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
