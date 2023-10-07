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
type docLocatorType = common.Locator[models.NamespaceKind, models.ResourceKind]

type KmsDocumentRef interface {
	GetNamespaceID() common.IdentifierWithKind[models.NamespaceKind]
	GetID() common.IdentifierWithKind[models.ResourceKind]
}

var _ KmsDocumentRef = (*docLocatorType)(nil)

type KmsDocument interface {
	KmsDocumentRef
	GetLocator() models.ResourceLocator
	stampUpdatedWithAuth(context.Context) TimeStorable
	setETag(azcore.ETag)
}

type BaseDoc struct {
	NamespaceID docNsIDType `json:"namespaceId"`
	ID          docIDType   `json:"id"`

	Updated       TimeStorable  `json:"updated"`
	Deleted       *TimeStorable `json:"deleted"`
	UpdatedBy     string        `json:"updatedBy"`
	SchemaVersion int           `json:"schemaVersion"`

	LinkTo   *docLocatorType  `json:"@link.to,omitempty"`
	LinkFrom []docLocatorType `json:"@link.from,omitempty"`

	ETag azcore.ETag         `json:"-"`    // populated during read
	Kind models.ResourceKind `json:"kind"` // populate during write for index
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

func (doc *BaseDoc) stampUpdatedWithAuth(c context.Context) TimeStorable {
	var callerPrincipalIdStr string
	var callerPrincipalName string
	if identity, ok := auth.GetAuthIdentity(c); ok {
		callerPrincipalIdStr = identity.ClientPrincipalID().String()
		callerPrincipalName = identity.ClientPrincipalName()
	}
	doc.Kind = doc.ID.Kind()
	doc.Updated = TimeStorable(time.Now().UTC())
	doc.UpdatedBy = fmt.Sprintf("%s:%s", callerPrincipalIdStr, callerPrincipalName)
	return doc.Updated
}

var _ KmsDocument = (*BaseDoc)(nil)

func Read[D KmsDocument](c common.ServiceContext, locator docLocatorType, target D) error {
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

func (d *BaseDoc) PopulateResourceRef(r *models.ResourceRef) bool {
	if d == nil || r == nil {
		return false
	}
	r.Id = d.ID.Identifier()
	r.Locator = models.NewResourceLocator(d.NamespaceID, d.ID)
	r.Updated = d.Updated.TimePtr()
	r.UpdatedBy = &d.UpdatedBy
	r.Deleted = d.Deleted.TimePtr()
	return true
}
