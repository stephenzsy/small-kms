package kmsdoc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
)

type docNsIDType = common.IdentifierWithKind[models.NamespaceKind]
type docIDType = common.IdentifierWithKind[models.ResourceKind]

type KmsDocumentRef interface {
	GetNamespaceID() common.IdentifierWithKind[models.NamespaceKind]
	GetID() common.IdentifierWithKind[models.ResourceKind]
}

var _ KmsDocumentRef = models.ResourceLocator{}

type KmsDocument interface {
	KmsDocumentRef
	GetLocator() models.ResourceLocator
	stampUpdatedWithAuth(context.Context) time.Time
	GetETag() azcore.ETag
	setETag(azcore.ETag)
	setAliasToWithETag(target models.ResourceLocator, etag azcore.ETag)
}

type KmsDocumentSnapshotable[D KmsDocument] interface {
	KmsDocument
	SnapshotWithNewLocator(models.ResourceLocator) D
}

type DocFlag string

const (
	DocFlagAliasSnapshot DocFlag = "alias-snapshot"
)

type BaseDoc struct {
	NamespaceID docNsIDType `json:"namespaceId"`
	ID          docIDType   `json:"id"`

	Updated       time.Time  `json:"updated"`
	Deleted       *time.Time `json:"deleted"`
	UpdatedBy     string     `json:"updatedBy"`
	SchemaVersion int        `json:"schemaVersion"`

	Flags []DocFlag `json:"@flags,omitempty"`

	AliasTo     *models.ResourceLocator  `json:"@alias.to,omitempty"`
	AliasToETag *azcore.ETag             `json:"@alias.to.etag,omitempty"`
	AliasFrom   []models.ResourceLocator `json:"@alias.from,omitempty"`

	ETag azcore.ETag         `json:"-"`    // populated during read
	Kind models.ResourceKind `json:"kind"` // populate during write for index
}

// setAliasToWithETag implements KmsDocument.
func (doc *BaseDoc) setAliasToWithETag(target models.ResourceLocator, etag azcore.ETag) {
	doc.AliasTo = &target
	doc.AliasToETag = &etag
}

func getDefaultQueryColumns() []string {
	return []string{
		"namespaceId",
		"id",
		"updated",
		"deleted",
		"updatedBy",
	}
}

// GetID implements KmsDocument.
func (doc *BaseDoc) GetID() docIDType {
	return doc.ID
}

func (doc *BaseDoc) GetETag() azcore.ETag {
	return doc.ETag
}

// GetNamespaceID implements KmsDocument.
func (doc *BaseDoc) GetNamespaceID() docNsIDType {
	return doc.NamespaceID
}

func (doc *BaseDoc) GetLocator() models.ResourceLocator {
	return models.NewResourceLocator(doc.NamespaceID, doc.ID)
}

func (doc *BaseDoc) setETag(etag azcore.ETag) {
	doc.ETag = etag
}

func (doc *BaseDoc) stampUpdatedWithAuth(c context.Context) time.Time {
	var callerPrincipalIdStr string
	var callerPrincipalName string
	if identity, ok := auth.GetAuthIdentity(c); ok {
		callerPrincipalIdStr = identity.ClientPrincipalID().String()
		callerPrincipalName = identity.ClientPrincipalName()
	}
	doc.Kind = doc.ID.Kind()
	now := time.Now().UTC()
	doc.Updated = now
	doc.UpdatedBy = fmt.Sprintf("%s:%s", callerPrincipalIdStr, callerPrincipalName)
	return now
}

func stampUpdatedWithAuthPatchOps(c context.Context, patchOps *azcosmos.PatchOperations) time.Time {
	var callerPrincipalIdStr string
	var callerPrincipalName string
	if identity, ok := auth.GetAuthIdentity(c); ok {
		callerPrincipalIdStr = identity.ClientPrincipalID().String()
		callerPrincipalName = identity.ClientPrincipalName()
	}
	now := time.Now().UTC()
	patchOps.AppendSet("/updated", now.Format(time.RFC3339))
	patchOps.AppendSet("/updatedBy", fmt.Sprintf("%s:%s", callerPrincipalIdStr, callerPrincipalName))
	return now
}

var _ KmsDocument = (*BaseDoc)(nil)

func Read[D KmsDocument](c common.ServiceContext, locator models.ResourceLocator, target D) error {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(locator.GetNamespaceID().String())
	id := locator.GetID().String()
	resp, err := cc.ReadItem(c, partitionKey, id, nil)
	if err != nil {
		return common.WrapAzRsNotFoundErr(err, fmt.Sprintf("doc:%s/%s", partitionKey, id))
	}
	err = json.Unmarshal(resp.Value, target)
	target.setETag(resp.ETag)
	return err
}

func Create[D KmsDocument](c common.ServiceContext, doc D) error {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	doc.stampUpdatedWithAuth(c)
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String())
	resp, err := cc.CreateItem(c, partitionKey, content, nil)
	if err != nil {
		return err
	}
	doc.setETag(resp.ETag)
	return err
}

func Upsert[D KmsDocument](c common.ServiceContext, doc D) error {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	doc.stampUpdatedWithAuth(c)
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	partitionKey := azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String())
	resp, err := cc.UpsertItem(c, partitionKey, content, nil)
	if err != nil {
		return err
	}
	doc.setETag(resp.ETag)
	return err
}

func Delete[D KmsDocument](c common.ServiceContext, doc D) (err error) {
	return DeleteByRef(c, doc)
}

func DeleteByRef(c common.ServiceContext, locator KmsDocumentRef) (err error) {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(locator.GetNamespaceID().String())
	_, err = cc.DeleteItem(c, partitionKey, locator.GetID().String(), nil)
	return err
}

func Patch[D KmsDocument](c common.ServiceContext, locator models.ResourceLocator, doc D,
	patchOps azcosmos.PatchOperations) error {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(locator.GetNamespaceID().String())
	stampUpdatedWithAuthPatchOps(c, &patchOps)
	resp, err := cc.PatchItem(c, partitionKey, locator.GetID().String(), patchOps, nil)
	if err != nil {
		return err
	}
	doc.setETag(resp.ETag)
	return err
}

func (d *BaseDoc) PopulateResourceRef(r *models.ResourceRef) {
	if d == nil || r == nil {
		return
	}
	r.Id = d.ID.Identifier()
	r.Locator = models.NewResourceLocator(d.NamespaceID, d.ID)
	r.Updated = &d.Updated
	r.UpdatedBy = &d.UpdatedBy
	r.Deleted = d.Deleted
}

func UpsertAliasWithSnapshot[D KmsDocumentSnapshotable[D]](c common.ServiceContext, doc D, aliasLocator models.ResourceLocator) (docClone D, err error) {
	etag := doc.GetETag()
	if etag == "" {
		return docClone, fmt.Errorf("missing etag, target document must be saved to cosmosdb first")
	}
	docClone = doc.SnapshotWithNewLocator(aliasLocator)
	docClone.setAliasToWithETag(doc.GetLocator(), etag)
	err = Upsert(c, docClone)
	return
}
