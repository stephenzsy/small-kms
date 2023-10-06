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

// document name space type
type DocNsType string

const (
	DocNsTypeCaRoot           DocNsType = "ca-root"
	DocNsTypeCaInt            DocNsType = "ca-int"
	DocNsTypeDevice           DocNsType = "device"
	DocNsTypeApplication      DocNsType = "application"
	DocNsTypeServicePrincipal DocNsType = "service-principal"
	DocNsTypeUser             DocNsType = "user"
	DocNsTypeGroup            DocNsType = "group"
	DocNsTypeTenant           DocNsType = "tenant"
)

type DocType string

const (
	DocTypeProfile DocType = "profile"
)

type BaseDoc struct {
	NsType        DocNsType         `json:"nsType"`
	NsID          models.Identifier `json:"nsId"`
	DocType       DocType           `json:"docType"`
	DocID         models.Identifier `json:"docId"`
	Updated       time.Time         `json:"updated"`
	Deleted       *time.Time        `json:"deleted"`
	UpdatedBy     string            `json:"updatedBy"`
	SchemaVersion int               `json:"schemaVersion"`

	// document db used for storage in cosmos, to be populated prior to save
	StorageID   string      `json:"id"`
	NamespaceID string      `json:"namespaceId"`
	ETag        azcore.ETag `json:"-"` // populated during read

}

type KmsDocument interface {
	stampUpdatedWithAuth(context.Context) time.Time
	getDBKeys() (azcosmos.PartitionKey, string)
	setETag(azcore.ETag)
}

func (doc *BaseDoc) getDBKeys() (azcosmos.PartitionKey, string) {
	if doc.NsID.IsNilOrEmpty() {
		doc.NamespaceID = string(doc.NsType)
	} else {
		doc.NamespaceID = string(doc.NsType) + "/" + doc.NsID.String()
	}
	if doc.DocType == "" {
		doc.StorageID = string(doc.DocID.String())
	} else {
		doc.StorageID = string(doc.DocType) + "/" + doc.DocID.String()
	}
	return azcosmos.NewPartitionKeyString(doc.NamespaceID), doc.StorageID
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
	doc.Updated = time.Now()
	doc.UpdatedBy = fmt.Sprintf("%s:%s", callerPrincipalIdStr, callerPrincipalName)
	return doc.Updated
}

func ReadByKeyFunc[D KmsDocument](c common.ServiceContext, getKeys func() (string, string), doc D) error {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey, id := getKeys()
	resp, err := cc.ReadItem(c, azcosmos.NewPartitionKeyString(partitionKey), id, nil)
	if err != nil {
		return common.WrapAzRsNotFoundErr(err, fmt.Sprintf("doc:%s/%s", partitionKey, id))
	}
	err = json.Unmarshal(resp.Value, doc)
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
	partitionKey, _ := doc.getDBKeys()
	resp, err := cc.UpsertItem(c, partitionKey, content, nil)
	if err != nil {
		return err
	}
	doc.setETag(resp.ETag)
	return err
}

func Delete[D KmsDocument](c common.ServiceContext, doc D) (err error) {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey, id := doc.getDBKeys()
	_, err = cc.DeleteItem(c, partitionKey, id, nil)
	return err
}

func DeleteByKeyFunc(c common.ServiceContext, getKeys func() (string, string)) (err error) {
	cc := common.GetClientProvider(c).AzCosmosContainerClient()
	partitionKey, id := getKeys()
	_, err = cc.DeleteItem(c, azcosmos.NewPartitionKeyString(partitionKey), id, nil)
	return err
}
