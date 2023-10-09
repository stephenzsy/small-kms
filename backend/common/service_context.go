package common

import (
	"context"
	ctx "context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type ServerContext interface {
	ctx.Context
	ServiceClientProvider
}

type ServiceClientProvider interface {
	AzCosmosContainerClient() *azcosmos.ContainerClient
	AzKeysClient() *azkeys.Client
	AzCertificatesClient() *azcertificates.Client
	MsGraphServerClient() *msgraphsdkgo.GraphServiceClient
	MsGraphDelegatedClient(context.Context) (*msgraphsdkgo.GraphServiceClient, error)
	AzBlobContainerClient() *azblobcontainer.Client
}
