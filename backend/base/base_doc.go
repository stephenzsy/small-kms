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
	GetPersistedSLocator() SLocator
	getDefaultStorageNamespaceID() uuid.UUID
	getDefaultStorageID() uuid.UUID
	GetUpdatedBy() string
	getETag() *azcore.ETag
	setETag(etag azcore.ETag)
	setTimestamp(t time.Time)
	setUpdatedBy(string)
	prepareForWrite(c context.Context, storageNID, storageRID uuid.UUID)

	setRelationsFunc(func(*DocRelations) *DocRelations)
}

type RelName string

type DocRelations struct {
	NamedFrom map[RelName]SLocator `json:"namedFrom,omitempty"`
	NamedTo   map[RelName]SLocator `json:"namedTo,omitempty"`
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

	Relations *DocRelations `json:"@rels,omitempty"`
}

// setRelations implements CRUDDoc.
func (d *BaseDoc) setRelationsFunc(f func(*DocRelations) *DocRelations) {
	d.Relations = f(d.Relations)
}

// GetPersistedSLocator implements CRUDDoc.
func (b *BaseDoc) GetPersistedSLocator() SLocator {
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

	baseDocPatchColumnRelations          = "/@rels"
	baseDocPatchColumnRelationsNamedFrom = "/@rels/namedFrom"
	baseDocPatchColumnRelationsNamedTo   = "/@rels/namedTo"
)

// GetUpdatedBy implements CRUDDoc.
func (d *BaseDoc) GetUpdatedBy() string {
	return d.UpdatedBy
}

func GetDefaultStorageNamespaceIDURL(namespaceKind NamespaceKind, namespaceIdentifier Identifier) string {
	return fmt.Sprintf("https://example.com/v1/r/%s/%s", namespaceKind, namespaceIdentifier.String())
}

func GetDefaultStorageNamespaceID(namespaceKind NamespaceKind, namespaceIdentifier Identifier) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(GetDefaultStorageNamespaceIDURL(namespaceKind, namespaceIdentifier)))
}

func (d *BaseDoc) getDefaultStorageNamespaceID() uuid.UUID {
	return GetDefaultStorageNamespaceID(d.NamespaceKind, d.NamespaceIdentifier)
}

func GetDefaultStorageIDURL(storageNamespaceIDURL string, resourceKind ResourceKind, resourceIdentifier Identifier) string {
	return fmt.Sprintf("%s/%s/%s", storageNamespaceIDURL, resourceKind, resourceIdentifier.String())
}

func GetDefaultStorageID(storageNamespaceIDURL string, resourceKind ResourceKind, resourceIdentifier Identifier) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(GetDefaultStorageIDURL(storageNamespaceIDURL, resourceKind, resourceIdentifier)))
}

func GetDefaultStorageLocator(namespaceKind NamespaceKind, namespaceIdentifier Identifier,
	resourceKind ResourceKind, resourceIdentifier Identifier) SLocator {
	storageNamespaceIDURL := GetDefaultStorageNamespaceIDURL(namespaceKind, namespaceIdentifier)

	return SLocator{uuid.NewSHA1(uuid.NameSpaceURL, []byte(storageNamespaceIDURL)),
		uuid.NewSHA1(uuid.NameSpaceURL, []byte(GetDefaultStorageIDURL(storageNamespaceIDURL, resourceKind, resourceIdentifier)))}
}

func (d *BaseDoc) getDefaultStorageID() uuid.UUID {
	return GetDefaultStorageID(
		GetDefaultStorageNamespaceIDURL(d.NamespaceKind, d.NamespaceIdentifier),
		d.ResourceKind,
		d.ResourceIdentifier)
}

func (d *BaseDoc) getETag() *azcore.ETag {
	return d.ETag
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
	patchByLocator(context.Context, SLocator, azcosmos.PatchOperations, *azcosmos.ItemOptions) error
	NewQueryItemsPager(query string, storageNamespaceID uuid.UUID, o *azcosmos.QueryOptions) *azruntime.Pager[azcosmos.QueryItemsResponse]
	getClient() *azcosmos.ContainerClient
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
	return doc.getDefaultStorageNamespaceID()
}

func resolveStorageID(c context.Context, doc CRUDDoc) uuid.UUID {
	if doc, ok := doc.(CRUDDocHasCustomStorageID); ok {
		return doc.GetStorageID(c)
	}
	return doc.getDefaultStorageID()
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
		return HandleAzCosmosError(err)
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

func (s *azcosmosContainerCRUDDocService) patchByLocator(c context.Context, locator SLocator, ops azcosmos.PatchOperations, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(locator.NID.String())
	nextUpdatedBy := auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	ops.AppendSet(baseDocPatchColumnUpdatedBy, nextUpdatedBy)
	_, err := s.client.PatchItem(c, partitionKey, locator.RID.String(), ops, o)
	if err != nil {
		return err
	}
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
