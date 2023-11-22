package resdoc

import (
	"context"
	"encoding/json"
	"time"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
)

type ContextKey int

const (
	DocServiceContextKey ContextKey = iota
)

type DocService interface {
	Read(c context.Context, identifier DocIdentifier, dst ResourceDocument, opts *azcosmos.ItemOptions) error
	Create(context.Context, ResourceDocument, *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	Upsert(context.Context, ResourceDocument, *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	NewQueryItemsPager(query string, partitionKey PartitionKey, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse]
	Delete(context.Context, DocIdentifier, *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	Patch(context.Context, ResourceDocument, azcosmos.PatchOperations, *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
	PatchByIdentifier(context.Context, DocIdentifier, azcosmos.PatchOperations, *azcosmos.ItemOptions) (azcosmos.ItemResponse, error)
}

type azcosmosSingleContainerDocService struct {
	client *azcosmos.ContainerClient
}

// Patch implements DocService.
func (s *azcosmosSingleContainerDocService) Patch(c context.Context, doc ResourceDocument, patchOps azcosmos.PatchOperations, opts *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	partitionKey := azcosmos.NewPartitionKeyString(doc.partitionKey().String())
	nextUpdatedBy := auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	if doc.getUpdatedBy() != nextUpdatedBy {
		patchOps.AppendSet("/updatedBy", nextUpdatedBy)
	}
	resp, err := s.client.PatchItem(c, partitionKey, string(doc.getID()), patchOps, opts)
	if err != nil {
		return resp, err
	}
	doc.setUpdatedBy(nextUpdatedBy)
	doc.setETag(resp.ETag)
	ts, err := time.Parse(time.RFC1123, resp.RawResponse.Header.Get("Date"))
	if err != nil {
		log.Ctx(c).Warn().Err(err).Str("DateHeader", resp.RawResponse.Header.Get("Date")).Msg("failed to parse timestamp")
	} else {
		doc.setTimestamp(ts)
	}
	return resp, nil
}

func (s *azcosmosSingleContainerDocService) PatchByIdentifier(
	c context.Context, identifier DocIdentifier, patchOps azcosmos.PatchOperations, opts *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	partitionKey := azcosmos.NewPartitionKeyString(identifier.PartitionKey.String())
	nextUpdatedBy := auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	patchOps.AppendSet("/updatedBy", nextUpdatedBy)
	return s.client.PatchItem(c, partitionKey, identifier.ID, patchOps, opts)
}

// Delete implements DocService.
func (s *azcosmosSingleContainerDocService) Delete(c context.Context, docIdentifier DocIdentifier, opts *azcosmos.ItemOptions) (azcosmos.ItemResponse, error) {
	partitionKey := azcosmos.NewPartitionKeyString(docIdentifier.PartitionKey.String())
	return s.client.DeleteItem(c, partitionKey, docIdentifier.ID, opts)
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

func (s *azcosmosSingleContainerDocService) Create(
	c context.Context, doc ResourceDocument, o *azcosmos.ItemOptions) (resp azcosmos.ItemResponse, err error) {
	partitionKey := azcosmos.NewPartitionKeyString(doc.partitionKey().String())
	doc.prepareForWrite(c)
	content, err := json.Marshal(doc)
	if err != nil {
		return resp, err
	}
	resp, err = s.client.CreateItem(c, partitionKey, content, o)
	if err != nil {
		return resp, err
	}
	doc.setETag(resp.ETag)
	ts, err := time.Parse(time.RFC1123, resp.RawResponse.Header.Get("Date"))
	if err != nil {
		log.Ctx(c).Warn().Err(err).Str("DateHeader", resp.RawResponse.Header.Get("Date")).Msg("failed to parse timestamp")
	} else {
		doc.setTimestamp(ts)
	}
	return resp, nil
}

func (s *azcosmosSingleContainerDocService) Upsert(
	c context.Context, doc ResourceDocument, o *azcosmos.ItemOptions) (resp azcosmos.ItemResponse, err error) {
	partitionKey := azcosmos.NewPartitionKeyString(doc.partitionKey().String())
	doc.prepareForWrite(c)
	content, err := json.Marshal(doc)
	if err != nil {
		return resp, err
	}
	resp, err = s.client.UpsertItem(c, partitionKey, content, o)
	if err != nil {
		return resp, err
	}
	doc.setETag(resp.ETag)
	ts, err := time.Parse(time.RFC1123, resp.RawResponse.Header.Get("Date"))
	if err != nil {
		log.Ctx(c).Warn().Err(err).Str("DateHeader", resp.RawResponse.Header.Get("Date")).Msg("failed to parse timestamp")
	} else {
		doc.setTimestamp(ts)
	}
	return resp, nil
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
