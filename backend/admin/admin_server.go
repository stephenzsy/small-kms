package admin

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

type adminServer struct {
	azKeysClient                    *azkeys.Client
	azBlobClient                    *azblob.Client
	azBlobContainerClient           *azblobcontainer.Client
	azCosmosClient                  *azcosmos.Client
	azCosmosDatabaseClient          *azcosmos.DatabaseClient
	azCosmosContainerClientCerts    *azcosmos.ContainerClient
	azCosmosContainerClientPolicies *azcosmos.ContainerClient
}

type AdminServerInternal interface {
	ReadCertEnrollPolicyDBItem(ctx context.Context, namespaceID uuid.UUID) (result CertificateEnrollmentPolicyDTO, err error)
	ReadCertDBItem(c context.Context, namespaceID uuid.UUID, id uuid.UUID) (result CertDBItem, err error)
	FetchCertificatePEMBlob(ctx context.Context, blobName string) ([]byte, error)
}

func NewAdminServer() *adminServer {
	common.MustGetenv(common.DefualtEnvVarAzCosmosResourceEndpoint)
	common.MustGetenv(common.DefualtEnvVarAzKeyvaultResourceEndpoint)
	common.MustGetenv(common.DefualtEnvVarAzStroageBlobResourceEndpoint)

	s := adminServer{}
	var err error
	s.azKeysClient, err = common.GetAzKeysClient()
	if err != nil {
		log.Panic("Failed to get az keys client", err.Error())
	}
	s.azBlobClient, err = common.GetAzStorageBlobClient()
	if err != nil {
		log.Panicf("Failed to get az blob client: %s", err.Error())
	}
	s.azBlobContainerClient = s.azBlobClient.ServiceClient().NewContainerClient(common.GetEnvWithDefault("AZURE_STORAGEBLOB_CONTAINERNAME_CERTS", "certs"))

	s.azCosmosClient, err = common.GetAzCosmosClient()
	if err != nil {
		log.Panicf("Failed to get az cosmos client: %s", err.Error())
	}
	s.azCosmosDatabaseClient, err = s.azCosmosClient.NewDatabase(common.GetEnvWithDefault("AZURE_COSMOS_DATABASE_ID", "kms"))
	if err != nil {
		log.Panicf("Failed to get az cosmos database client: %s", err.Error())
	}
	s.azCosmosContainerClientCerts, err = s.azCosmosDatabaseClient.NewContainer(common.GetEnvWithDefault("AZURE_COSMOS_CONTAINERNAME_CERTS", "Certs"))
	if err != nil {
		log.Panicf("Failed to get az cosmos container client for Certs: %s", err.Error())
	}
	s.azCosmosContainerClientPolicies, err = s.azCosmosDatabaseClient.NewContainer(common.GetEnvWithDefault("AZURE_COSMOS_CONTAINERNAME_POLICIES", "Policies"))
	if err != nil {
		log.Panicf("Failed to get az cosmos container client for Policies: %s", err.Error())
	}
	return &s
}
