package base

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
)

type BaseDoc struct {
	StorageNamespaceID uuid.UUID `json:"namespaceId"`
	StroageID          uuid.UUID `json:"id"`

	NamespaceKind       NamespaceKind `json:"namespaceKind"`
	NamespaceIdentifier Identifier    `json:"namespaceIdentifier"`
	ResourceKind        ResourceKind  `json:"resourceKind"`
	ResourceIdentifier  Identifier    `json:"resourceIdentifier"`

	Timestamp *jwt.NumericDate `json:"_ts,omitempty"`
	ETag      *azcore.ETag     `json:"_etag,omitempty"`
	Deleted   *time.Time       `json:"deleted,omitempty"`
	UpdatedBy string           `json:"updatedBy,omitempty"`
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

func (d *BaseDoc) GetStorageIdBaseUrl(c context.Context) string {
	if val, ok := c.Value(SiteUrlContextKey).(string); ok {
		return val
	}
	return "https://example.com"
}

func (d *BaseDoc) GetStorageNamespaceIdUrl(c context.Context) string {
	return fmt.Sprintf("%s/v1/r/%s/%s", d.GetStorageIdBaseUrl(c), d.NamespaceKind, d.NamespaceIdentifier.String())
}

// default implementation get storage ID
func (d *BaseDoc) GetStorageNamespaceID(c context.Context) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(d.GetStorageNamespaceIdUrl(c)))
}

func (d *BaseDoc) GetStorageIdUrl(c context.Context) string {
	return fmt.Sprintf("%s/%s/%s", d.GetStorageNamespaceID(c), d.ResourceKind, d.ResourceIdentifier.String())
}

func (d *BaseDoc) GetStorageID(c context.Context) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(d.GetStorageIdUrl(c)))
}

// setETag implements CRUDDoc.
func (d *BaseDoc) setETag(eTag azcore.ETag) {
	d.ETag = &eTag
}

// setUpdated implements CRUDDoc.
func (d *BaseDoc) prepareForWrite(c context.Context) {
	d.StorageNamespaceID = d.GetStorageNamespaceID(c)
	d.StroageID = d.GetStorageID(c)
	d.UpdatedBy = auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	// clear read-only fields
	d.ETag = nil
	d.Timestamp = nil
}

type CRUDDocHasCustomStorageID interface {
	GetStorageID(context.Context) uuid.UUID
}

type CRUDDoc interface {
	GetStorageNamespaceID(context.Context) uuid.UUID
	CRUDDocHasCustomStorageID
	GetUpdatedBy() string
	setETag(etag azcore.ETag)
	setTimestamp(t time.Time)
	setUpdatedBy(string)
	prepareForWrite(c context.Context)
}

var _ CRUDDoc = (*BaseDoc)(nil)

type AzCosmosCRUDDocService interface {
	Create(context.Context, CRUDDoc, *azcosmos.ItemOptions) error
	Upsert(context.Context, CRUDDoc, *azcosmos.ItemOptions) error
	Read(context.Context, uuid.UUID, uuid.UUID, CRUDDoc, *azcosmos.ItemOptions) error
	Patch(context.Context, CRUDDoc, azcosmos.PatchOperations, *azcosmos.ItemOptions) error
	SoftDelete(context.Context)
	Purge(context.Context)
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

func (s *azcosmosContainerCRUDDocService) Create(c context.Context, doc CRUDDoc, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetStorageNamespaceID(c).String())
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
func (s *azcosmosContainerCRUDDocService) Read(c context.Context, namespaceStorageID, storageID uuid.UUID, dst CRUDDoc, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(namespaceStorageID.String())
	resp, err := s.client.ReadItem(c, partitionKey, storageID.String(), nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Value, dst)
	dst.setETag(resp.ETag)
	return err
}

// Upsert implements CRUDDocService.
func (s *azcosmosContainerCRUDDocService) Upsert(c context.Context, doc CRUDDoc, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetStorageNamespaceID(c).String())
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
func (s *azcosmosContainerCRUDDocService) Patch(c context.Context, doc CRUDDoc, ops azcosmos.PatchOperations, o *azcosmos.ItemOptions) error {
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetStorageNamespaceID(c).String())
	nextUpdatedBy := auth.GetAuthIdentity(c).ClientPrincipalDisplayName()
	if doc.GetUpdatedBy() != nextUpdatedBy {
		ops.AppendSet(baseDocPatchColumnUpdatedBy, nextUpdatedBy)
	}
	resp, err := s.client.PatchItem(c, partitionKey, doc.GetStorageID(c).String(), ops, o)
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
