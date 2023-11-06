package base

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
)

type BaseDocument interface {
	GetStorageNamespaceID() DocNamespacePartitionKey
	GetStorageID() ID
	GetStorageFullIdentifier() DocFullIdentifier
	GetUpdatedBy() string
	getETag() *azcore.ETag
	setETag(etag azcore.ETag)
	setTimestamp(t time.Time)
	setUpdatedBy(string)
	prepareForWrite(c context.Context)
}

type BaseDoc struct {
	PartitionKey DocNamespacePartitionKey `json:"namespaceId"`
	ID           ID                       `json:"id"`
	Timestamp    *jwt.NumericDate         `json:"_ts,omitempty"`
	ETag         *azcore.ETag             `json:"_etag,omitempty"`
	Deleted      *time.Time               `json:"deleted,omitempty"`
	UpdatedBy    string                   `json:"updatedBy,omitempty"`
}

func (d *BaseDoc) GetStorageNamespaceID() DocNamespacePartitionKey {
	return d.PartitionKey
}

// GetID implements BaseDocument.
func (d *BaseDoc) GetStorageID() ID {
	return d.ID
}

// GetPersistedSLocator implements CRUDDoc.
func (d *BaseDoc) GetStorageFullIdentifier() DocFullIdentifier {
	return DocFullIdentifier{
		d.PartitionKey,
		d.ID,
	}
}

func (d *BaseDoc) Init(nsKind NamespaceKind, nsID ID, rKind ResourceKind, rID ID) {
	d.PartitionKey = NewDocNamespacePartitionKey(nsKind, nsID, rKind)
	d.ID = rID
}

// setTimestamp implements CRUDDoc.
func (d *BaseDoc) setTimestamp(t time.Time) {
	d.Timestamp = jwt.NewNumericDate(t)
}

// setUpdatedBy implements CRUDDoc.
func (d *BaseDoc) setUpdatedBy(val string) {
	d.UpdatedBy = val
}

const (
	baseDocPatchColumnUpdatedBy = "/updatedBy"
)

// GetUpdatedBy implements CRUDDoc.
func (d *BaseDoc) GetUpdatedBy() string {
	return d.UpdatedBy
}

func (d *BaseDoc) getETag() *azcore.ETag {
	return d.ETag
}

// setETag implements CRUDDoc.
func (d *BaseDoc) setETag(eTag azcore.ETag) {
	d.ETag = &eTag
}

// setUpdated implements CRUDDoc.
func (d *BaseDoc) prepareForWrite(c context.Context) {
	d.UpdatedBy = auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	// clear read-only fields
	d.ETag = nil
	d.Timestamp = nil
}

var _ BaseDocument = (*BaseDoc)(nil)

type AzCosmosCRUDDocService interface {
	Create(context.Context, BaseDocument, *azcosmos.ItemOptions) error
	Upsert(context.Context, BaseDocument, *azcosmos.ItemOptions) error
	Read(c context.Context, docFullIdentifier DocFullIdentifier, dst BaseDocument, opts *azcosmos.ItemOptions) error
	Patch(context.Context, BaseDocument, azcosmos.PatchOperations, *azcosmos.ItemOptions) error
	NewQueryItemsPager(query string, storageNamespaceID DocNamespacePartitionKey, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse]
	getClient() *azcosmos.ContainerClient
	SoftDelete(c context.Context, doc BaseDocument, opts *azcosmos.ItemOptions) error
	Delete(c context.Context, doc BaseDocument, opts *azcosmos.ItemOptions) error
}

func NewAzCosmosCRUDDocService(client *azcosmos.ContainerClient) *azcosmosContainerCRUDDocService {
	return &azcosmosContainerCRUDDocService{
		client: client,
	}
}

func GetAzCosmosCRUDService(c context.Context) AzCosmosCRUDDocService {
	if s, ok := c.Value(AzCosmosCRUDDocServiceContextKey).(AzCosmosCRUDDocService); ok {
		return s
	}
	return nil
}

type azcosmosContainerCRUDDocService struct {
	client *azcosmos.ContainerClient
}

func (s *azcosmosContainerCRUDDocService) Create(c context.Context, doc BaseDocument, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetStorageNamespaceID().String())
	doc.prepareForWrite(c)
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	resp, err := s.client.CreateItem(c, partitionKey, content, o)
	if err != nil {
		return err
	}
	doc.setETag(resp.ETag)
	doc.setTimestamp(time.Now())
	return nil
}

// Read implements CRUDDocService.
func (s *azcosmosContainerCRUDDocService) Read(c context.Context, docFullID DocFullIdentifier, dst BaseDocument, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(docFullID.pKey.String())
	resp, err := s.client.ReadItem(c, partitionKey, string(docFullID.docID), nil)
	if err != nil {
		return HandleAzCosmosError(err)
	}
	err = json.Unmarshal(resp.Value, dst)
	dst.setETag(resp.ETag)
	return err
}

func (s *azcosmosContainerCRUDDocService) NewQueryItemsPager(query string, storageNamespaceID DocNamespacePartitionKey, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse] {
	partitionKey := azcosmos.NewPartitionKeyString(storageNamespaceID.String())
	return s.client.NewQueryItemsPager(query, partitionKey, o)
}

// Upsert implements CRUDDocService.
func (s *azcosmosContainerCRUDDocService) Upsert(c context.Context, doc BaseDocument, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetStorageNamespaceID().String())
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
	doc.setTimestamp(time.Now())
	return nil
}

// Patch implements CRUDDocService.
// this operation does not update fields patched, fields need to be updated manually after call is done
func (s *azcosmosContainerCRUDDocService) Patch(c context.Context, doc BaseDocument, ops azcosmos.PatchOperations, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetStorageNamespaceID().String())
	nextUpdatedBy := auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	if doc.GetUpdatedBy() != nextUpdatedBy {
		ops.AppendSet(baseDocPatchColumnUpdatedBy, nextUpdatedBy)
	}
	resp, err := s.client.PatchItem(c, partitionKey, string(doc.GetStorageID()), ops, o)
	if err != nil {
		return err
	}
	doc.setUpdatedBy(nextUpdatedBy)
	doc.setETag(resp.ETag)
	doc.setTimestamp(time.Now())
	return nil
}

func (s *azcosmosContainerCRUDDocService) Delete(c context.Context, doc BaseDocument, opts *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetStorageNamespaceID().String())
	_, err := s.client.DeleteItem(c, partitionKey, string(doc.GetStorageID()), opts)
	return err
}

// SoftDelete implements AzCosmosCRUDDocService.
func (s *azcosmosContainerCRUDDocService) SoftDelete(c context.Context, doc BaseDocument, opts *azcosmos.ItemOptions) error {
	patchOps := azcosmos.PatchOperations{}
	patchOps.AppendSet("/deleted", time.Now().UTC())
	return s.Patch(c, doc, patchOps, opts)
}

// getClient implements AzCosmosCRUDDocService.
func (s *azcosmosContainerCRUDDocService) getClient() *azcosmos.ContainerClient {
	return s.client
}

var _ AzCosmosCRUDDocService = (*azcosmosContainerCRUDDocService)(nil)

// PopulateModelRef implements ModelRefPopulater.
func (d *BaseDoc) PopulateModelRef(m *ResourceReference) {
	if d == nil || m == nil {
		return
	}
	m.ID = d.ID
	m.Updated = d.Timestamp.Time
	m.Deleted = d.Deleted
	m.UpdatedBy = &d.UpdatedBy
}

var _ ModelRefPopulater[ResourceReference] = (*BaseDoc)(nil)
