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
)

type KmsDocument interface {
	stampUpdatedWithAuth(context.Context) time.Time
	setETag(azcore.ETag)
	GetNamespaceID() DocNsID
	GetID() DocID
}

type BaseDoc struct {
	NamespaceID DocNsID `json:"namespaceId"`
	ID          DocID   `json:"id"`

	Updated       time.Time  `json:"updated"`
	Deleted       *time.Time `json:"deleted"`
	UpdatedBy     string     `json:"updatedBy"`
	SchemaVersion int        `json:"schemaVersion"`

	ETag azcore.ETag `json:"-"`    // populated during read
	Kind DocKind     `json:"kind"` // populate during write for index
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
func (doc *BaseDoc) GetID() DocID {
	return doc.ID
}

// GetNamespaceID implements KmsDocument.
func (doc *BaseDoc) GetNamespaceID() DocNsID {
	return doc.NamespaceID
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
	doc.Kind = doc.ID.kind
	doc.Updated = time.Now()
	doc.UpdatedBy = fmt.Sprintf("%s:%s", callerPrincipalIdStr, callerPrincipalName)
	return doc.Updated
}

var _ KmsDocument = (*BaseDoc)(nil)

func Read[D KmsDocument](c common.ServiceContext, nsID DocNsID, id DocID, target D) error {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(nsID.String())
	resp, err := cc.ReadItem(c, partitionKey, id.String(), nil)
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
	return DeleteByKey(c, doc.GetNamespaceID(), doc.GetID())
}

func DeleteByKey(c common.ServiceContext, nsID DocNsID, id DocID) (err error) {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey := azcosmos.NewPartitionKeyString(nsID.String())
	_, err = cc.DeleteItem(c, partitionKey, id.String(), nil)
	return err
}
