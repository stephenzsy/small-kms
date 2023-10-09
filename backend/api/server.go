package api

import (
	ctx "context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type server struct {
	common.CommonConfig
	serverContext         ctx.Context
	azBlobClient          *azblob.Client
	azBlobContainerClient *azblobcontainer.Client
	serverMsGraphClient   *msgraphsdkgo.GraphServiceClient
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

// MsGraphServerClient implements common.ServerContext.
func (s *server) MsGraphServerClient() *msgraphsdkgo.GraphServiceClient {
	return s.serverMsGraphClient
}

// Value implements common.ServerContext.
func (s *server) Value(key any) any {
	return s.serverContext.Value(key)
}

// AzBlobContainerClient implements common.ClientProvider.
func (s *server) AzBlobContainerClient() *azblobcontainer.Client {
	return s.azBlobContainerClient
}

type H = map[string]string
type RequestContext = common.RequestContext

func wrapResponse[T interface{}](c echo.Context, defaultStatus int, data T, err error) error {
	switch {
	case err == nil:
		return c.JSON(defaultStatus, data)
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

// MsGraphDelegatedClient implements common.ClientProvider.
func (s *server) MsGraphDelegatedClient(c ctx.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	if authIdentity, ok := auth.GetAuthIdentity(c); ok {
		if creds, err := authIdentity.GetOnBehalfOfTokenCredential(s, nil); err != nil {
			return nil, err
		} else {
			return msgraphsdkgo.NewGraphServiceClientWithCredentials(creds, nil)
		}
	}
	return nil, fmt.Errorf("%w: no auth header to authenticate to graph service", common.ErrStatusUnauthorized)
}

func NewServer(c ctx.Context) (models.ServerInterface, echo.MiddlewareFunc) {
	commonConfig, err := common.NewCommonConfig()
	if err != nil {
		log.Panic().Err(err).Msg("failed to create common config")
	}
	s := server{
		CommonConfig:  &commonConfig,
		serverContext: c,
	}
	storageBlobEndpoint := common.MustGetenv(common.DefualtEnvVarAzStroageBlobResourceEndpoint)
	s.azBlobClient, err = azblob.NewClient(storageBlobEndpoint, s.DefaultAzCredential(), nil)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to get az blob client")
	}
	s.azBlobContainerClient = s.azBlobClient.ServiceClient().NewContainerClient(common.GetEnvWithDefault("AZURE_STORAGEBLOB_CONTAINERNAME_CERTS", "certs"))
	s.serverMsGraphClient, err = msgraphsdkgo.NewGraphServiceClientWithCredentials(s.ConfidentialAppCredential(), nil)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to get ms graph client")
	}

	return &s, s.InjectServerContext()
}

func (s *server) InjectServerContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(common.EchoContextWithServerContext(c, s))
		}
	}
}
