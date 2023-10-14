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
	"github.com/stephenzsy/small-kms/backend/shared"
)

type RequestContext = common.RequestContext

type KmsDocumentRef interface {
	GetNamespaceID() shared.NamespaceIdentifier
	GetID() shared.ResourceIdentifier
}

var _ KmsDocumentRef = shared.ResourceLocator{}

type KmsDocument interface {
	KmsDocumentRef
	GetLocator() shared.ResourceLocator
	stampUpdatedWithAuth(context.Context) time.Time
	GetETag() azcore.ETag
	setETag(azcore.ETag)
	setAliasToWithETag(target shared.ResourceLocator, etag azcore.ETag)
}

type KmsDocumentSnapshotable[D KmsDocument] interface {
	KmsDocument
	SnapshotWithNewLocator(shared.ResourceLocator) D
}

type DocFlag string

type BaseDoc struct {
	NamespaceID shared.NamespaceIdentifier `json:"namespaceId"`
	ID          shared.ResourceIdentifier  `json:"id"`

	Updated       time.Time  `json:"updated"`
	Deleted       *time.Time `json:"deleted"`
	UpdatedBy     string     `json:"updatedBy"`
	SchemaVersion int        `json:"schemaVersion"`

	AliasTo     *shared.ResourceLocator           `json:"@alias.to,omitempty"`
	AliasToETag *azcore.ETag                      `json:"@alias.to.etag,omitempty"`
	Owner       *shared.ResourceLocator           `json:"@owner,omitempty"`
	Owns        map[string]shared.ResourceLocator `json:"@owns,omitempty"`

	ETag azcore.ETag         `json:"-"`    // populated during read
	Kind shared.ResourceKind `json:"kind"` // populate during write for index
}

// setAliasToWithETag implements KmsDocument.
func (doc *BaseDoc) setAliasToWithETag(target shared.ResourceLocator, etag azcore.ETag) {
	doc.AliasTo = &target
	doc.AliasToETag = &etag
}

const QueryColumnNameAliasTo = "c[\"@alias.to\"]"
const QueryColumnNameOwner = "c[\"@owner\"]"

var queryDefaultColumns = []string{
	"c.namespaceId",
	"c.id",
	"c.updated",
	"c.deleted",
	"c.updatedBy",
}

// GetID implements KmsDocument.
func (doc *BaseDoc) GetID() shared.ResourceIdentifier {
	return doc.ID
}

func (doc *BaseDoc) GetETag() azcore.ETag {
	return doc.ETag
}

// GetNamespaceID implements KmsDocument.
func (doc *BaseDoc) GetNamespaceID() shared.NamespaceIdentifier {
	return doc.NamespaceID
}

func (doc *BaseDoc) GetLocator() shared.ResourceLocator {
	if doc == nil {
		return shared.ResourceLocator{}
	}
	if doc.AliasTo != nil {
		return *doc.AliasTo
	}
	return shared.NewResourceLocator(doc.NamespaceID, doc.ID)
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

func Read[D KmsDocument](c context.Context, locator shared.ResourceLocator, target D) error {
	cc := common.GetAdminServerClientProvider(c).AzCosmosContainerClient()
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

func Create[D KmsDocument](c context.Context, doc D) error {
	cc := common.GetAdminServerClientProvider(c).AzCosmosContainerClient()
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
	return nil
}

func Upsert[D KmsDocument](c context.Context, doc D) error {
	cc := common.GetAdminServerClientProvider(c).AzCosmosContainerClient()
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
	return nil
}

func Delete[D KmsDocument](c context.Context, doc D) (err error) {
	return DeleteByRef(c, doc)
}

func DeleteByRef(c context.Context, locator KmsDocumentRef) (err error) {
	cc := common.GetAdminServerClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(locator.GetNamespaceID().String())
	_, err = cc.DeleteItem(c, partitionKey, locator.GetID().String(), nil)
	return err
}

func PatchWithWriteBack[D KmsDocument](c context.Context,
	locator shared.ResourceLocator,
	dstDoc D,
	patchOps azcosmos.PatchOperations) error {
	cc := common.GetAdminServerClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(locator.GetNamespaceID().String())
	stampUpdatedWithAuthPatchOps(c, &patchOps)
	resp, err := cc.PatchItem(c, partitionKey, locator.GetID().String(), patchOps, &azcosmos.ItemOptions{EnableContentResponseOnWrite: true})
	if err != nil {
		return err
	}
	dstDoc.setETag(resp.ETag)
	return json.Unmarshal(resp.Value, dstDoc)
}

func Patch[D KmsDocument](c context.Context,
	doc D,
	patchOps azcosmos.PatchOperations,
	opts *azcosmos.ItemOptions) error {
	cc := common.GetAdminServerClientProvider(c).AzCosmosContainerClient()
	locator := doc.GetLocator()
	partitionKey := azcosmos.NewPartitionKeyString(locator.GetNamespaceID().String())
	stampUpdatedWithAuthPatchOps(c, &patchOps)
	resp, err := cc.PatchItem(c, partitionKey, locator.GetID().String(), patchOps, opts)
	if err != nil {
		return err
	}
	doc.setETag(resp.ETag)
	return nil
}

func (d *BaseDoc) PopulateResourceRef(r *shared.ResourceRef) {
	if d == nil || r == nil {
		return
	}
	r.Id = d.ID.Identifier()
	r.Locator = shared.NewResourceLocator(d.NamespaceID, d.ID)
	r.Updated = &d.Updated
	r.UpdatedBy = &d.UpdatedBy
	r.Deleted = d.Deleted
}

func UpsertAliasWithSnapshot[D KmsDocumentSnapshotable[D]](c context.Context, doc D, aliasLocator shared.ResourceLocator) (docClone D, err error) {
	etag := doc.GetETag()
	if etag == "" {
		return docClone, fmt.Errorf("missing etag, target document must be saved to cosmosdb first")
	}
	docClone = doc.SnapshotWithNewLocator(aliasLocator)
	docClone.setAliasToWithETag(doc.GetLocator(), etag)
	err = Upsert(c, docClone)
	return
}
