package api

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
)

type APIServer interface {
	RespondRequireAdmin(c echo.Context) error
}

type apiServer struct {
	chCtx                context.Context
	siteURL              string
	docService           base.AzCosmosCRUDDocService
	serviceMsGraphClient *msgraphsdkgo.GraphServiceClient
	legacyClientProvider common.AdminServerClientProvider
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
	case graph.ServiceMsGraphClientContextKey:
		return s.serviceMsGraphClient
	case common.AdminServerClientProviderContextKey:
		return s.legacyClientProvider
	}
	return nil
}

var _ context.Context = (*apiServer)(nil)

func NewApiServer(c context.Context, serverOld *server) *apiServer {

	return &apiServer{
		chCtx:                c,
		siteURL:              common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, "SITE_URL", "https://example.com"),
		docService:           base.NewAzCosmosCRUDDocService(serverOld.clients.azCosmosContainerClientCerts),
		serviceMsGraphClient: serverOld.clients.msGraphClient,
		legacyClientProvider: &serverOld.clients,
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

var _ APIServer = (*apiServer)(nil)
