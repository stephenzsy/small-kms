package common

import (
	ctx "context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type serviceContextKey string

const (
	clientProviderContextKey serviceContextKey = "serviceProvider"
)

type ClientProvider interface {
	AzCosmosContainerClient() *azcosmos.ContainerClient
	MsGraphDelegatedClient(ctx.Context) (*msgraphsdkgo.GraphServiceClient, error)
}

type ServiceContext ctx.Context

func WithClientProvider(parent ctx.Context, provider ClientProvider) ServiceContext {
	return ctx.WithValue(parent, clientProviderContextKey, provider)
}

func GetClientProvider(s ServiceContext) ClientProvider {
	return s.Value(clientProviderContextKey).(ClientProvider)
}
