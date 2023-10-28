package api

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
)

type APIServer interface {
	RespondRequireAdmin(c echo.Context) error
}

type apiServer struct {
	chCtx                   context.Context
	serviceIdentity         auth.AzureIdentity
	siteURL                 string
	docService              base.AzCosmosCRUDDocService
	serviceMsGraphClient    *msgraphsdkgo.GraphServiceClient
	azCertificatesClient    *azcertificates.Client
	azKeysClient            *azkeys.Client
	legacyClientProvider    common.AdminServerClientProvider
	appConfidentialIdentity auth.AzureAppConfidentialIdentity
}

// AzCertificatesClient implements kv.AzKeyVaultService.
func (s *apiServer) AzCertificatesClient() *azcertificates.Client {
	return s.azCertificatesClient
}

// AzKeysClient implements kv.AzKeyVaultService.
func (s *apiServer) AzKeysClient() *azkeys.Client {
	return s.azKeysClient
}

// respondRequireAdmin implements APIServer.
func (*apiServer) RespondRequireAdmin(c echo.Context) error {
	return respondRequireAdmin(c)
}

// Deadline implements context.Context.
func (s *apiServer) Deadline() (deadline time.Time, ok bool) {
	return s.chCtx.Deadline()
}

// Done implements context.Context.
func (s *apiServer) Done() <-chan struct{} {
	return s.chCtx.Done()
}

// Err implements context.Context.
func (s *apiServer) Err() error {
	return s.chCtx.Err()
}

// Value implements context.Context.
func (s *apiServer) Value(key any) any {
	switch key {
	case base.SiteUrlContextKey:
		return s.siteURL
	case base.AzCosmosCRUDDocServiceContextKey:
		return s.docService
	case kv.AzKeyVaultServiceContextKey:
		return s
	case graph.ServiceClientIDContextKey:
		return s.serviceIdentity.ClientID()
	case graph.ServiceMsGraphClientContextKey:
		return s.serviceMsGraphClient
	case graph.ServiceMsGraphClientClientIDContextKey:
		return s.appConfidentialIdentity.ClientID()
	case auth.AppConfidentialIdentityContextKey:
		return s.appConfidentialIdentity
	case common.AdminServerClientProviderContextKey:
		return s.legacyClientProvider
	}
	return nil
}

var _ context.Context = (*apiServer)(nil)

func NewApiServer(c context.Context, serverOld *server) *apiServer {

	return &apiServer{
		chCtx:                   c,
		serviceIdentity:         serverOld.ServiceIdentity(),
		siteURL:                 common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, "SITE_URL", "https://example.com"),
		docService:              base.NewAzCosmosCRUDDocService(serverOld.clients.azCosmosContainerClientCerts),
		serviceMsGraphClient:    serverOld.clients.msGraphClient,
		azCertificatesClient:    serverOld.clients.azCertificatesClient,
		azKeysClient:            serverOld.clients.azKeysClient,
		appConfidentialIdentity: serverOld.appIdentity,
		legacyClientProvider:    &serverOld.clients,
	}
}

func (s *apiServer) InjectServiceContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c = ctx.NewInjectedRequestContext(c, s)
			return next(c)
		}
	}
}

var _ kv.AzKeyVaultService = (*apiServer)(nil)
var _ APIServer = (*apiServer)(nil)
