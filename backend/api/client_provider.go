package api

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
)

const (
	DefualtEnvVarAzKeyvaultResourceEndpoint    = "AZURE_KEYVAULT_RESOURCEENDPOINT"
	DefualtEnvVarAzCosmosResourceEndpoint      = "AZURE_COSMOS_RESOURCEENDPOINT"
	DefualtEnvVarAzStroageBlobResourceEndpoint = "AZURE_STORAGEBLOB_RESOURCEENDPOINT"
)

type clientProvider struct {
	azKeyvaultEndpoint   string
	cachedKeyvaultName   string
	azKeysClient         *azkeys.Client
	azCertificatesClient *azcertificates.Client

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

// AzCertificatesClient implements common.AdminServerClientProvider.
func (p *clientProvider) AzCertificatesClient() *azcertificates.Client {
	return p.azCertificatesClient
}

// AzKeysClient implements common.AdminServerClientProvider.
func (p *clientProvider) AzKeysClient() *azkeys.Client {
	return p.azKeysClient
}

func (p *clientProvider) AzCosmosContainerClient() *azcosmos.ContainerClient {
	return p.azCosmosContainerClientCerts
}

func (p *clientProvider) keyvaultName() string {
	return p.cachedKeyvaultName
}

var _ common.AdminServerClientProvider = (*clientProvider)(nil)

func newServerClientProvider(s *server) (p clientProvider, err error) {

	creds := s.ServiceIdentity().TokenCredential()
	if p.azKeyvaultEndpoint = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixService, DefualtEnvVarAzKeyvaultResourceEndpoint, ""); p.azKeyvaultEndpoint == "" {
		err = fmt.Errorf("%w: %s", common.ErrMissingEnvVar, DefualtEnvVarAzKeyvaultResourceEndpoint)
		return
	}

	if parsed, parseErr := url.Parse(p.azKeyvaultEndpoint); parseErr == nil {
		p.cachedKeyvaultName = strings.Split(parsed.Host, ".")[0]
	} else {
		err = fmt.Errorf("%w: %s=%s", common.ErrInvalidEnvVar, DefualtEnvVarAzKeyvaultResourceEndpoint, p.azKeyvaultEndpoint)
		return
	}

	if p.azKeysClient, err = azkeys.NewClient(p.azKeyvaultEndpoint, creds, nil); err != nil {
		return
	}

	if p.azCertificatesClient, err = azcertificates.NewClient(p.azKeyvaultEndpoint, creds, nil); err != nil {
		return
	}

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

// GetKeyvaultCertificateResourceScopeID implements common.AdminServerRequestClientProvider.
func (p *requestClientProvider) GetKeyvaultCertificateResourceScopeID(certificateName string, category string) string {
	if category != "secrets" {
		category = "certificates"
	}
	return fmt.Sprintf("subscriptions/%s/resourceGroups/%s/providers/Microsoft.KeyVault/vaults/%s/%s/%s",
		p.parent.subscriptionId,
		p.parent.resourceGroupName,
		p.parent.clients.keyvaultName(),
		category,
		certificateName)
}

func (p *requestClientProvider) getOnbehalfOfCreds() (azcore.TokenCredential, error) {
	var err error
	if p.onBehalfOfCreds == nil {
		if authIdentity, ok := auth.GetAuthIdentity(p.credentialContext); ok {
			if p.onBehalfOfCreds, err = authIdentity.GetOnBehalfOfTokenCredential(p.parent.appIdentity, nil); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("%w: only authorized request can get delegated client", common.ErrStatusUnauthorized)
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

// ArmRoleAssignmentsClient implements common.AdminServerRequestClientProvider.
func (p *requestClientProvider) ArmRoleAssignmentsClient() (*armauthorization.RoleAssignmentsClient, error) {
	var err error
	if p.cachedArmRoleAssignmentsClient == nil {
		var creds azcore.TokenCredential
		if creds, err = p.getOnbehalfOfCreds(); err != nil {
			return nil, err
		}
		p.cachedArmRoleAssignmentsClient, err = armauthorization.NewRoleAssignmentsClient(p.parent.subscriptionId, creds, nil)
	}
	return p.cachedArmRoleAssignmentsClient, err
}

var _ common.AdminServerRequestClientProvider = (*requestClientProvider)(nil)
