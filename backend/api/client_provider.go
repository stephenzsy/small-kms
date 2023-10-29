package api

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
)

const (
	DefualtEnvVarAzKeyvaultResourceEndpoint    = "AZURE_KEYVAULT_RESOURCEENDPOINT"
	DefualtEnvVarAzCosmosResourceEndpoint      = "AZURE_COSMOS_RESOURCEENDPOINT"
	DefualtEnvVarAzStroageBlobResourceEndpoint = "AZURE_STORAGEBLOB_RESOURCEENDPOINT"
)

type clientProvider struct {
	azCosmosEndpoint             string
	azCosmosClient               *azcosmos.Client
	azCosmosDatabaseID           string
	azCosmosDatabaseClient       *azcosmos.DatabaseClient
	azCosmosContainerID          string
	azCosmosContainerClientCerts *azcosmos.ContainerClient

	azBlobEndpoint        string
	azBlobClient          *azblob.Client
	azBlobContainerID     string
	azBlobContainerClient *azblobcontainer.Client

	msGraphClient *msgraphsdkgo.GraphServiceClient
}

// MsGraphClient implements common.AdminServerClientProvider.
func (p *clientProvider) MsGraphClient() *msgraphsdkgo.GraphServiceClient {
	return p.msGraphClient
}

// CertsAzBlobContainerClient implements common.AdminServerClientProvider.
func (p *clientProvider) CertsAzBlobContainerClient() *azblobcontainer.Client {
	return p.azBlobContainerClient
}

func (p *clientProvider) AzCosmosContainerClient() *azcosmos.ContainerClient {
	return p.azCosmosContainerClientCerts
}

var _ common.AdminServerClientProvider = (*clientProvider)(nil)

func extractKeyVaultName(keyvaultEndpoing string) string {
	if parsed, err := url.Parse(keyvaultEndpoing); err == nil {
		return strings.Split(parsed.Host, ".")[0]
	}
	return ""
}

func newServerClientProvider(s *server) (p clientProvider, err error) {

	creds := s.ServiceIdentity().TokenCredential()

	if cosmosConnStr := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, "AZURE_COSMOS_CONNECTION_STRING", ""); cosmosConnStr != "" {
		p.azCosmosClient, err = azcosmos.NewClientFromConnectionString(cosmosConnStr, nil)
		if err != nil {
			log.Panicf("Failed to create az cosmos client from connection string: %s", err.Error())
		}
	} else {
		if p.azCosmosEndpoint = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, DefualtEnvVarAzCosmosResourceEndpoint, ""); p.azCosmosEndpoint == "" {
			err = fmt.Errorf("%w: %s", common.ErrMissingEnvVar, DefualtEnvVarAzCosmosResourceEndpoint)
			return
		}
		if p.azCosmosClient, err = azcosmos.NewClient(p.azCosmosEndpoint, creds, nil); err != nil {
			return
		}
	}

	p.azCosmosDatabaseID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, "AZURE_COSMOS_DATABASE_ID", "kms")
	if p.azCosmosDatabaseClient, err = p.azCosmosClient.NewDatabase(p.azCosmosDatabaseID); err != nil {
		return
	}
	p.azCosmosContainerID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, "AZURE_COSMOS_CONTAINERNAME_CERTS", "Certs")
	if p.azCosmosContainerClientCerts, err = p.azCosmosDatabaseClient.NewContainer(p.azCosmosContainerID); err != nil {
		return
	}

	if p.azBlobEndpoint = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, DefualtEnvVarAzStroageBlobResourceEndpoint, ""); p.azBlobEndpoint == "" {
		err = fmt.Errorf("%w: %s", common.ErrMissingEnvVar, DefualtEnvVarAzStroageBlobResourceEndpoint)
		return
	}
	if p.azBlobClient, err = azblob.NewClient(p.azBlobEndpoint, creds, nil); err != nil {
		return
	}
	p.azBlobContainerID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, "AZURE_STORAGEBLOB_CONTAINERNAME_CERTS", "certs")
	p.azBlobContainerClient = p.azBlobClient.ServiceClient().NewContainerClient(p.azBlobContainerID)

	if p.msGraphClient, err = msgraphsdkgo.NewGraphServiceClientWithCredentials(s.appIdentity.TokenCredential(), nil); err != nil {
		return
	}
	return
}

type requestClientProvider struct {
	parent                         *server
	credentialContext              RequestContext
	onBehalfOfCreds                azcore.TokenCredential
	cachedDelegatedMsGraphClient   *msgraphsdkgo.GraphServiceClient
	cachedArmRoleAssignmentsClient *armauthorization.RoleAssignmentsClient
}

func (p *requestClientProvider) getOnbehalfOfCreds() (azcore.TokenCredential, error) {
	var err error
	if p.onBehalfOfCreds == nil {
		authIdentity := auth.GetAuthIdentity(p.credentialContext)
		if p.onBehalfOfCreds, err = authIdentity.GetOnBehalfOfTokenCredential(p.credentialContext, nil); err != nil {
			return nil, err
		}
	}
	return p.onBehalfOfCreds, err
}

// MsGraphClient implements common.AdminServerRequestClientProvider.
func (p *requestClientProvider) MsGraphClient() (*msgraphsdkgo.GraphServiceClient, error) {
	var err error
	if p.cachedDelegatedMsGraphClient == nil {
		var creds azcore.TokenCredential
		if creds, err = p.getOnbehalfOfCreds(); err != nil {
			return nil, err
		}
		p.cachedDelegatedMsGraphClient, err = msgraphsdkgo.NewGraphServiceClientWithCredentials(creds, nil)
	}
	return p.cachedDelegatedMsGraphClient, err
}

var _ common.AdminServerRequestClientProvider = (*requestClientProvider)(nil)
