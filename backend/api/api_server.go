package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/labstack/echo/v4"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/graph"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type APIServer interface {
	EnvService() common.EnvService
	RespondRequireAdmin(c echo.Context) error
	GetAzKeyVaultEndpoint() string
	GetBuildID() string
	WithDelegatedARMAuthRoleAssignmentsClient(c ctx.RequestContext) (ctx.RequestContext, *armauthorization.RoleAssignmentsClient, error)
	GetAzSubscriptionID() string
	GetResourceGroupName() string
	GetKeyVaultName() string
}

type apiServer struct {
	common.CommonServer
	parentCtx               context.Context
	docService              base.AzCosmosCRUDDocService
	docServiceNew           resdoc.DocService
	serviceMsGraphClient    *msgraphsdkgo.GraphServiceClient
	azKeyVaultEndpoint      string
	azCertificatesClient    *azcertificates.Client
	azKeysClient            *azkeys.Client
	azSecretsClient         *azsecrets.Client
	appConfidentialIdentity auth.AzureAppConfidentialIdentity
	buildID                 string
	azSubscriptionID        string
	resourceGroupName       string
	extractedKeyVaultName   string

	azCosmosEndpoint        string
	azCosmosClient          *azcosmos.Client
	azCosmosDatabaseID      string
	azCosmosDatabaseClient  *azcosmos.DatabaseClient
	azCosmosContainerID     string
	azCosmosContainerClient *azcosmos.ContainerClient
}

// AzSecretsClient implements kv.AzKeyVaultService.
func (s *apiServer) AzSecretsClient() *azsecrets.Client {
	return s.azSecretsClient
}

// GetBuildID implements APIServer.
func (s *apiServer) GetBuildID() string {
	return s.buildID
}

// GetAzKeyVaultEndpoint implements APIServer.
func (s *apiServer) GetAzKeyVaultEndpoint() string {
	return s.azKeyVaultEndpoint
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
	return c.JSON(http.StatusForbidden, map[string]string{"message": "admin access required"})
}

func RespondPagerList[T any](c ctx.RequestContext, pager *utils.SerializableItemsPager[T]) error {
	jsonBlob, err := json.Marshal(pager)
	if err != nil {
		log.Ctx(c).Error().Err(err).Send()
		return c.String(http.StatusInternalServerError, "internal error")
	}
	return c.JSONBlob(http.StatusOK, jsonBlob)
}

// Deadline implements context.Context.
func (s *apiServer) Deadline() (deadline time.Time, ok bool) {
	return s.parentCtx.Deadline()
}

// Done implements context.Context.
func (s *apiServer) Done() <-chan struct{} {
	return s.parentCtx.Done()
}

// Err implements context.Context.
func (s *apiServer) Err() error {
	return s.parentCtx.Err()
}

// Value implements context.Context.
func (s *apiServer) Value(key any) any {
	switch key {
	case base.AzCosmosCRUDDocServiceContextKey:
		return s.docService
	case resdoc.DocServiceContextKey:
		return s.docServiceNew
	case kv.AzKeyVaultServiceContextKey:
		return s
	case graph.ServiceClientIDContextKey:
		return s.ServiceIdentity().ClientID()
	case graph.ServiceMsGraphClientContextKey:
		return s.serviceMsGraphClient
	case graph.ServiceMsGraphClientClientIDContextKey:
		return s.appConfidentialIdentity.ClientID()
	case auth.AppConfidentialIdentityContextKey:
		return s.appConfidentialIdentity
	}
	return s.parentCtx.Value(key)
}

func NewApiServer(c context.Context, buildID string) (*apiServer, error) {
	commonConfig, err := common.NewCommonConfig(common.NewEnvService(), buildID)
	if err != nil {
		return nil, err
	}
	s := &apiServer{
		CommonServer: commonConfig,
		parentCtx:    c,
		buildID:      buildID,
	}
	var ok bool

	s.appConfidentialIdentity, err = getAppConfidentialIdentity(s.EnvService())
	if err != nil {
		return nil, err
	}

	// cosmos
	if cosmosConnStr := s.EnvService().Default(envKeyAzCosmosConnectionString, "", common.IdentityEnvVarPrefixService); cosmosConnStr != "" {
		s.azCosmosClient, err = azcosmos.NewClientFromConnectionString(cosmosConnStr, nil)
		if err != nil {
			return nil, err
		}
	} else if s.azCosmosEndpoint, ok = s.EnvService().RequireNonWhitespace(envKeyAzCosmosResourceEndpoint, common.IdentityEnvVarPrefixService); !ok {
		return nil, s.EnvService().ErrMissing(envKeyAzCosmosResourceEndpoint)
	} else if s.azCosmosClient, err = azcosmos.NewClient(s.azCosmosEndpoint, s.ServiceIdentity().TokenCredential(), nil); err != nil {
		return nil, err
	}

	s.azCosmosDatabaseID = s.EnvService().Default(envKeyAzCosmosDatabaseID, "kms", common.IdentityEnvVarPrefixService)
	if s.azCosmosDatabaseClient, err = s.azCosmosClient.NewDatabase(s.azCosmosDatabaseID); err != nil {
		return nil, err
	}
	s.azCosmosContainerID = s.EnvService().Default(envKeyAzCosmosContainerName, "Certs", common.IdentityEnvVarPrefixService)
	if s.azCosmosContainerClient, err = s.azCosmosDatabaseClient.NewContainer(s.azCosmosContainerID); err != nil {
		return nil, err
	}
	s.docService = base.NewAzCosmosCRUDDocService(s.azCosmosContainerClient)
	s.docServiceNew = resdoc.NewAzCosmosSingleContainerDocService(s.azCosmosContainerClient)

	// keyvault
	if s.azKeyVaultEndpoint, ok = s.EnvService().RequireNonWhitespace(common.EnvKeyAzKeyvaultResourceEndpoint, common.IdentityEnvVarPrefixService); !ok {
		return s, s.EnvService().ErrMissing(common.EnvKeyAzKeyvaultResourceEndpoint)
	}
	s.extractedKeyVaultName = cloudkeyaz.ExtractKeyVaultName(s.azKeyVaultEndpoint)
	if s.azKeysClient, err = azkeys.NewClient(s.azKeyVaultEndpoint, s.ServiceIdentity().TokenCredential(), nil); err != nil {
		return s, err
	}
	if s.azCertificatesClient, err = azcertificates.NewClient(s.azKeyVaultEndpoint, s.ServiceIdentity().TokenCredential(), nil); err != nil {
		return s, err
	}
	if s.azSecretsClient, err = azsecrets.NewClient(s.azKeyVaultEndpoint, s.ServiceIdentity().TokenCredential(), nil); err != nil {
		return s, err
	}

	s.azSubscriptionID = s.EnvService().Default(common.EnvKeyAzSubscriptionID, "", common.IdentityEnvVarPrefixService)
	s.resourceGroupName = s.EnvService().Default(common.EnvKeyAzResourceGroupName, "", common.IdentityEnvVarPrefixService)

	s.serviceMsGraphClient, err = msgraphsdkgo.NewGraphServiceClientWithCredentials(s.appConfidentialIdentity.TokenCredential(), nil)
	if err != nil {
		return s, err
	}

	return s, nil
}

var _ kv.AzKeyVaultService = (*apiServer)(nil)
var _ APIServer = (*apiServer)(nil)
var _ context.Context = (*apiServer)(nil)

func (s *apiServer) GetAzSubscriptionID() string {
	return s.azSubscriptionID
}

func (s *apiServer) GetResourceGroupName() string {
	return s.resourceGroupName
}

func (s *apiServer) GetKeyVaultName() string {
	return s.extractedKeyVaultName
}
