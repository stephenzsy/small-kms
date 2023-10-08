package common

import (
	ctx "context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	azblobcontainer "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type serviceContextKey string

const (
	clientProviderContextKey serviceContextKey = "serviceProvider"
)

type ClientProvider interface {
	AzCosmosContainerClient() *azcosmos.ContainerClient
	AzKeysClient() *azkeys.Client
	MsGraphDelegatedClient(ctx.Context) (*msgraphsdkgo.GraphServiceClient, error)
	AzBlobContainerClient() *azblobcontainer.Client
}

type ServiceContext ctx.Context

func WithClientProvider(parent ctx.Context, provider ClientProvider) ServiceContext {
	return ctx.WithValue(parent, clientProviderContextKey, provider)
}

func GetClientProvider(s ServiceContext) ClientProvider {
	return s.Value(clientProviderContextKey).(ClientProvider)
}
