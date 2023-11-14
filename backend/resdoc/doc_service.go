package resdoc

import (
	"context"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type ContextKey int

const (
	DocServiceContextKey ContextKey = iota
)

type DocService interface {
	NewQueryItemsPager(query string, partitionKey PartitionKey, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse]
}

type azcosmosSingleContainerDocService struct {
	client *azcosmos.ContainerClient
}

func (s *azcosmosSingleContainerDocService) NewQueryItemsPager(
	query string,
	partitionKey PartitionKey,
	o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse] {

	return s.client.NewQueryItemsPager(query, azcosmos.NewPartitionKeyString(partitionKey.String()), o)
}

var _ DocService = (*azcosmosSingleContainerDocService)(nil)

func NewAzCosmosSingleContainerDocService(client *azcosmos.ContainerClient) *azcosmosSingleContainerDocService {
	return &azcosmosSingleContainerDocService{
		client: client,
	}
}

func GetDocService(c context.Context) DocService {
	return c.Value(DocServiceContextKey).(DocService)
}
