package common

import (
	ctx "context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type serviceContextKey string

const (
	azCosmosContainerClientContextKey serviceContextKey = "azCosmosContainerClient"
)

type ServiceContext ctx.Context

func CreateServiceContext(parent ctx.Context, client *azcosmos.ContainerClient) ServiceContext {
	return ctx.WithValue(parent, azCosmosContainerClientContextKey, client)
}

func GetAzCosmosContainerClient(s ServiceContext) *azcosmos.ContainerClient {
	return s.Value(azCosmosContainerClientContextKey).(*azcosmos.ContainerClient)
}
