package base

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
)

type CRUDDocHasCustomStorageNamespaceID interface {
	GetStorageNamespaceID(context.Context) uuid.UUID
}

type CRUDDocHasCustomStorageID interface {
	GetStorageID(context.Context) uuid.UUID
}

type CRUDDoc interface {
	// can only be used on a doc that has been read from storage
	GetSLocator() SLocator
	getDefaultStorageNamespaceID(c context.Context) uuid.UUID
	getDefaultStorageID(c context.Context) uuid.UUID
	GetUpdatedBy() string
	setETag(etag azcore.ETag)
	setTimestamp(t time.Time)
	setUpdatedBy(string)
	prepareForWrite(c context.Context, storageNID, storageRID uuid.UUID)
}

type BaseDoc struct {
	StorageNamespaceID uuid.UUID `json:"namespaceId"`
	StorageID          uuid.UUID `json:"id"`

	NamespaceKind       NamespaceKind `json:"namespaceKind"`
	NamespaceIdentifier Identifier    `json:"namespaceIdentifier"`
	ResourceKind        ResourceKind  `json:"resourceKind"`
	ResourceIdentifier  Identifier    `json:"resourceIdentifier"`

	Timestamp *jwt.NumericDate `json:"_ts,omitempty"`
	ETag      *azcore.ETag     `json:"_etag,omitempty"`
	Deleted   *time.Time       `json:"deleted,omitempty"`
	UpdatedBy string           `json:"updatedBy,omitempty"`
}

// GetSLocator implements CRUDDoc.
func (b *BaseDoc) GetSLocator() SLocator {
	return SLocator{
		b.StorageNamespaceID,
		b.StorageID,
	}
}

var queryDefaultColumns = []string{
	"c.id",
	"c.resourceIdentifier",
	"c._ts",
	"c.deleted",
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

func GetDefaultStorageIDURLBase(c context.Context) string {
	if val, ok := c.Value(SiteUrlContextKey).(string); ok {
		return val
	}
	return "https://example.com"
}

func GetDefaultStorageNamespaceIDURL(c context.Context, namespaceKind NamespaceKind, namespaceIdentifier Identifier) string {
	return fmt.Sprintf("%s/v1/r/%s/%s", GetDefaultStorageIDURLBase(c), namespaceKind, namespaceIdentifier.String())
}

func GetDefaultStorageNamespaceID(c context.Context, namespaceKind NamespaceKind, namespaceIdentifier Identifier) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(GetDefaultStorageNamespaceIDURL(c, namespaceKind, namespaceIdentifier)))
}

func (d *BaseDoc) getDefaultStorageNamespaceID(c context.Context) uuid.UUID {
	return GetDefaultStorageNamespaceID(c, d.NamespaceKind, d.NamespaceIdentifier)
}

func GetDefaultStorageIDURL(c context.Context, storageNamespaceIDURL string, resourceKind ResourceKind, resourceIdentifier Identifier) string {
	return fmt.Sprintf("%s/%s/%s", storageNamespaceIDURL, resourceKind, resourceIdentifier.String())
}

func GetDefaultStorageID(c context.Context, storageNamespaceIDURL string, resourceKind ResourceKind, resourceIdentifier Identifier) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(GetDefaultStorageIDURL(c, storageNamespaceIDURL, resourceKind, resourceIdentifier)))
}

func GetDefaultStorageLocator(c context.Context,
	namespaceKind NamespaceKind, namespaceIdentifier Identifier,
	resourceKind ResourceKind, resourceIdentifier Identifier) (uuid.UUID, uuid.UUID) {
	storageNamespaceIDURL := GetDefaultStorageNamespaceIDURL(c, namespaceKind, namespaceIdentifier)

	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(storageNamespaceIDURL)),
		uuid.NewSHA1(uuid.NameSpaceURL, []byte(GetDefaultStorageIDURL(c, storageNamespaceIDURL, resourceKind, resourceIdentifier)))
}

func (d *BaseDoc) getDefaultStorageID(c context.Context) uuid.UUID {
	return GetDefaultStorageID(
		c,
		GetDefaultStorageNamespaceIDURL(c, d.NamespaceKind, d.NamespaceIdentifier),
		d.ResourceKind,
		d.ResourceIdentifier)
}

// setETag implements CRUDDoc.
func (d *BaseDoc) setETag(eTag azcore.ETag) {
	d.ETag = &eTag
}

// setUpdated implements CRUDDoc.
func (d *BaseDoc) prepareForWrite(c context.Context, sNID, sRID uuid.UUID) {
	d.StorageNamespaceID = sNID
	d.StorageID = sRID
	d.UpdatedBy = auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	// clear read-only fields
	d.ETag = nil
	d.Timestamp = nil
}

var _ CRUDDoc = (*BaseDoc)(nil)

type AzCosmosCRUDDocService interface {
	Create(context.Context, CRUDDoc, *azcosmos.ItemOptions) error
	Upsert(context.Context, CRUDDoc, *azcosmos.ItemOptions) error
	Read(c context.Context, storageNamespaceID, storageID uuid.UUID, dst CRUDDoc, opts *azcosmos.ItemOptions) error
	Patch(context.Context, CRUDDoc, azcosmos.PatchOperations, *azcosmos.ItemOptions) error
	NewQueryItemsPager(query string, storageNamespaceID uuid.UUID, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse]
	// TODO: SoftDelete(context.Context)
	// TODO: Purge(context.Context)
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

func resolveStorageNamespaceID(c context.Context, doc CRUDDoc) uuid.UUID {
	if doc, ok := doc.(CRUDDocHasCustomStorageNamespaceID); ok {
		return doc.GetStorageNamespaceID(c)
	}
	return doc.getDefaultStorageNamespaceID(c)
}

func resolveStorageID(c context.Context, doc CRUDDoc) uuid.UUID {
	if doc, ok := doc.(CRUDDocHasCustomStorageID); ok {
		return doc.GetStorageID(c)
	}
	return doc.getDefaultStorageID(c)
}

func (s *azcosmosContainerCRUDDocService) Create(c context.Context, doc CRUDDoc, o *azcosmos.ItemOptions) error {
	sNID := resolveStorageNamespaceID(c, doc)
	partitionKey := azcosmos.NewPartitionKeyString(sNID.String())
	doc.prepareForWrite(c, sNID, resolveStorageID(c, doc))
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
func (s *azcosmosContainerCRUDDocService) Read(c context.Context, storageNamespaceID, storageID uuid.UUID, dst CRUDDoc, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(storageNamespaceID.String())
	resp, err := s.client.ReadItem(c, partitionKey, storageID.String(), nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Value, dst)
	dst.setETag(resp.ETag)
	return err
}

func (s *azcosmosContainerCRUDDocService) NewQueryItemsPager(query string, storageNamespaceID uuid.UUID, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse] {
	partitionKey := azcosmos.NewPartitionKeyString(storageNamespaceID.String())
	return s.client.NewQueryItemsPager(query, partitionKey, o)
}

// Upsert implements CRUDDocService.
func (s *azcosmosContainerCRUDDocService) Upsert(c context.Context, doc CRUDDoc, o *azcosmos.ItemOptions) error {
	sNID := resolveStorageNamespaceID(c, doc)
	partitionKey := azcosmos.NewPartitionKeyString(sNID.String())
	doc.prepareForWrite(c, sNID, resolveStorageID(c, doc))
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
func (s *azcosmosContainerCRUDDocService) Patch(c context.Context, doc CRUDDoc, ops azcosmos.PatchOperations, o *azcosmos.ItemOptions) error {
	sNID := resolveStorageNamespaceID(c, doc)
	partitionKey := azcosmos.NewPartitionKeyString(sNID.String())
	nextUpdatedBy := auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	if doc.GetUpdatedBy() != nextUpdatedBy {
		ops.AppendSet(baseDocPatchColumnUpdatedBy, nextUpdatedBy)
	}
	resp, err := s.client.PatchItem(c, partitionKey, resolveStorageID(c, doc).String(), ops, o)
	if err != nil {
		return err
	}
	doc.setUpdatedBy(nextUpdatedBy)
	doc.setETag(resp.ETag)
	doc.setTimestamp(time.Now())
	return nil
}

// Purge implements AzCosmosCRUDDocService.
func (*azcosmosContainerCRUDDocService) Purge(context.Context) {
	panic("unimplemented")
}

// SoftDelete implements AzCosmosCRUDDocService.
func (*azcosmosContainerCRUDDocService) SoftDelete(context.Context) {
	panic("unimplemented")
}

var _ AzCosmosCRUDDocService = (*azcosmosContainerCRUDDocService)(nil)

// PopulateModelRef implements ModelRefPopulater.
func (d *BaseDoc) PopulateModelRef(m *ResourceReference) {
	if d == nil || m == nil {
		return
	}
	m.Id = SLocator{d.StorageNamespaceID, d.StorageID}
	m.Updated = d.Timestamp.Time
	m.Deleted = d.Deleted
	m.UpdatedBy = d.UpdatedBy
	m.NamespaceKind = d.NamespaceKind
	m.NamespaceIdentifier = d.NamespaceIdentifier
	m.ResourceKind = d.ResourceKind
	m.ResourceIdentifier = d.ResourceIdentifier
}

var _ ModelRefPopulater[ResourceReference] = (*BaseDoc)(nil)
