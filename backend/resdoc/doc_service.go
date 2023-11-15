package resdoc

import (
	"context"
	"encoding/json"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/rs/zerolog/log"
)

type ContextKey int

const (
	DocServiceContextKey ContextKey = iota
)

type DocService interface {
	Read(c context.Context, identifier DocIdentifier, dst ResourceDocument, opts *azcosmos.ItemOptions) error
	Upsert(context.Context, ResourceDocument, *azcosmos.ItemOptions) error
	NewQueryItemsPager(query string, partitionKey PartitionKey, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse]
}

type azcosmosSingleContainerDocService struct {
	client *azcosmos.ContainerClient
}

func (s *azcosmosSingleContainerDocService) Read(
	c context.Context,
	identifier DocIdentifier,
	dst ResourceDocument, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(identifier.PartitionKey.String())
	resp, err := s.client.ReadItem(c, partitionKey, identifier.ID, nil)
	if err != nil {
		return HandleAzCosmosError(err)
	}
	err = json.Unmarshal(resp.Value, dst)
	dst.setETag(resp.ETag)
	return err
}

func (s *azcosmosSingleContainerDocService) Upsert(
	c context.Context, doc ResourceDocument, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.partitionKey().String())
	doc.prepareForWrite(c)
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	resp, err := s.client.UpsertItem(c, partitionKey, content, o)
	if err != nil {
		return err
	}
	doc.setETag(resp.ETag)
	ts, err := time.Parse(time.RFC1123, resp.RawResponse.Header.Get("Date"))
	if err != nil {
		log.Ctx(c).Warn().Err(err).Str("DateHeader", resp.RawResponse.Header.Get("Date")).Msg("failed to parse timestamp")
	} else {
		doc.setTimestamp(ts)
	}
	return nil
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
