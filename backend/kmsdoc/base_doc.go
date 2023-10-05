package kmsdoc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
)

type KmsDocType byte

const (
	DocTypeUnknown KmsDocType = 90 // Z

	DocTypeCert                  KmsDocType = 67 // C
	DocTypeMsGraphObject         KmsDocType = 71 // G
	DocTypeLatestCertForTemplate KmsDocType = 76 // L
	DocTypeDirectoryObject       KmsDocType = 79 // O, deprecated
	DocTypePendingCert           KmsDocType = 80 // P,
	DocTypeNamespaceRelation     KmsDocType = 82 // R
	DocTypePolicyState           KmsDocType = 83 // S, deprecated
	DocTypeCertTemplate          KmsDocType = 84 // T
)

type KmsDocTypeName string

const (
	DocTypeNameUnknown KmsDocTypeName = "unknown"

	DocTypeNameCert                  KmsDocTypeName = "cert"
	DocTypeNameMsGraphObject         KmsDocTypeName = "msgraph-object"
	DocTypeNameLatestCertForTemplate KmsDocTypeName = "cert-latest"
	DocTypeNameDirectoryObject       KmsDocTypeName = "directory-object"
	DocTypeNamePendingCert           KmsDocTypeName = "cert-pending"
	DocTypeNameNamespaceRelation     KmsDocTypeName = "namespace-relation"
	DocTypeNamePolicyState           KmsDocTypeName = "policy-state"
	DocTypeNameCertTemplate          KmsDocTypeName = "cert-template"
)

// KmsDocID is a unique identifier for a KmsDoc, is comparable
type KmsDocID struct {
	typeByte KmsDocType
	uuid     uuid.UUID
}

func NewKmsDocID(typ KmsDocType, id uuid.UUID) KmsDocID {
	return KmsDocID{
		typeByte: typ,
		uuid:     id,
	}
}

func (k KmsDocID) String() string {
	return fmt.Sprintf("%s%s", string(k.typeByte), k.uuid.String())
}

func (k *KmsDocID) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func (k *KmsDocID) GetUUID() uuid.UUID {
	return k.uuid
}

func (k *KmsDocID) UnmarshalJSON(b []byte) (err error) {
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return
	}
	k.typeByte = KmsDocType(s[0])
	k.uuid, err = uuid.Parse(s[1:])
	return err
}

func (k *KmsDocID) GetType() KmsDocType {
	return k.typeByte
}

type BaseDoc struct {
	ID            KmsDocID   `json:"id"`
	NamespaceID   uuid.UUID  `json:"namespaceId"`
	Updated       time.Time  `json:"updated"`
	UpdatedBy     string     `json:"updatedBy"`
	UpdatedByName string     `json:"updatedByName"`
	Deleted       *time.Time `json:"deleted,omitempty"`

	// used only for serialization and query
	TypeName KmsDocTypeName `json:"type"`

	// metadata
	ETag azcore.ETag `json:"-"`
}

func GetBaseDocQueryColumns(prefix string) string {
	return fmt.Sprintf("%s.id,%s.namespaceId,%s.updated,%s.updatedBy,%s.deleted", prefix, prefix, prefix, prefix, prefix)
}

type KmsDocument interface {
	GetNamespaceID() uuid.UUID
	StampUpdatedWithAuth(context.Context) time.Time
	StampDeletedWithAuth(context.Context) time.Time

	GetUUID() uuid.UUID
	GetDocID() KmsDocID
	GetUpdated() time.Time
	GetUpdatedBy() (string, string)
	GetDeleted() *time.Time
	SetETag(azcore.ETag)

	fillTypeName()
}

func (k *BaseDoc) GetUUID() uuid.UUID {
	return k.ID.GetUUID()
}

func (doc *BaseDoc) GetNamespaceID() uuid.UUID {
	return doc.NamespaceID
}

func (doc *BaseDoc) GetDocID() KmsDocID {
	return doc.ID
}

func (doc *BaseDoc) GetUpdated() time.Time {
	return doc.Updated
}

func (doc *BaseDoc) GetUpdatedBy() (string, string) {
	return doc.UpdatedBy, doc.UpdatedByName
}

func (doc *BaseDoc) GetDeleted() *time.Time {
	return doc.Deleted
}

func (doc *BaseDoc) SetETag(etag azcore.ETag) {
	doc.ETag = etag
}

func (doc *BaseDoc) StampUpdated(callerId string, callerName string) time.Time {
	doc.Updated = time.Now()
	doc.UpdatedBy = callerId
	doc.UpdatedByName = callerName
	return doc.Updated
}

func (doc *BaseDoc) StampUpdatedWithAuth(c context.Context) time.Time {
	var callerPrincipalIdStr string
	var callerPrincipalName string
	if identity, ok := auth.GetAuthIdentity(c); ok {
		callerPrincipalIdStr = identity.ClientPrincipalID().String()
		callerPrincipalName = identity.ClientPrincipalName()
	}
	return doc.StampUpdated(callerPrincipalIdStr, callerPrincipalName)
}

func (doc *BaseDoc) StampDeletedWithAuth(c context.Context) time.Time {
	time := doc.StampUpdatedWithAuth(c)
	doc.Deleted = &time
	return time
}

var docTypeNameMap = map[KmsDocType]KmsDocTypeName{
	DocTypeCert:                  DocTypeNameCert,
	DocTypeMsGraphObject:         DocTypeNameMsGraphObject,
	DocTypeLatestCertForTemplate: DocTypeNameLatestCertForTemplate,
	DocTypeDirectoryObject:       DocTypeNameDirectoryObject,
	DocTypePendingCert:           DocTypeNamePendingCert,
	DocTypeNamespaceRelation:     DocTypeNameNamespaceRelation,
	DocTypePolicyState:           DocTypeNamePolicyState,
	DocTypeCertTemplate:          DocTypeNameCertTemplate,
}

func (doc *BaseDoc) fillTypeName() {
	doc.TypeName = docTypeNameMap[doc.ID.typeByte]
}

func AzCosmosCreate[D KmsDocument](ctx context.Context, cc *azcosmos.ContainerClient, doc D) error {
	doc.StampUpdatedWithAuth(ctx)
	doc.fillTypeName()
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	resp, err := cc.CreateItem(ctx, azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String()), content, nil)
	doc.SetETag(resp.ETag)
	return err
}

func AzCosmosUpsert[D KmsDocument](ctx context.Context, cc *azcosmos.ContainerClient, doc D) error {
	doc.StampUpdatedWithAuth(ctx)
	doc.fillTypeName()
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	resp, err := cc.UpsertItem(ctx, azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String()), content, nil)

	doc.SetETag(resp.ETag)
	return err
}

func AzCosmosReadItem[D KmsDocument](ctx common.ServiceContext, nsID uuid.UUID, docID KmsDocID, target D) error {
	resp, err := common.GetAzCosmosContainerClient(ctx).ReadItem(ctx, azcosmos.NewPartitionKeyString(nsID.String()), docID.String(), nil)
	if err != nil {
		return common.WrapAzRsNotFoundErr(err, fmt.Sprintf("doc:%s/%s", nsID, docID))
	}
	err = json.Unmarshal(resp.Value, target)
	target.SetETag(resp.ETag)
	return err
}

// Deprecated: use AzCosmosReadItem instead
func AzCosmosRead[D KmsDocument](ctx context.Context,
	cc *azcosmos.ContainerClient,
	namespaceID uuid.UUID,
	docID KmsDocID, target D) error {
	resp, err := cc.ReadItem(ctx, azcosmos.NewPartitionKeyString(namespaceID.String()), docID.String(), nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp.Value, target)
	target.SetETag(resp.ETag)
	return err
}

func AzCosmosPatch[D KmsDocument](ctx context.Context, cc *azcosmos.ContainerClient, doc D, getPatchOps ...func(*azcosmos.PatchOperations, D)) error {
	ops := azcosmos.PatchOperations{}
	for _, getPatchOpsFunc := range getPatchOps {
		getPatchOpsFunc(&ops, doc)
	}
	doc.StampUpdatedWithAuth(ctx)
	ops.AppendSet("/updated", doc.GetUpdated())
	updatedBy, updatedByName := doc.GetUpdatedBy()
	ops.AppendSet("/updatedBy", updatedBy)
	ops.AppendSet("/updatedByName", updatedByName)
	_, err := cc.PatchItem(ctx, azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String()), doc.GetDocID().String(), ops, nil)
	return err
}

func AzCosmosDelete[D KmsDocument](ctx context.Context, cc *azcosmos.ContainerClient, doc D) (err error) {
	_, err = cc.DeleteItem(ctx, azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String()), doc.GetDocID().String(), nil)
	return
}

func AzCosmosSoftDelete[D KmsDocument](ctx context.Context, cc *azcosmos.ContainerClient, doc D) error {
	ops := azcosmos.PatchOperations{}
	doc.StampDeletedWithAuth(ctx)
	ops.AppendSet("/updated", doc.GetUpdated())
	updatedBy, updatedByName := doc.GetUpdatedBy()
	ops.AppendSet("/updatedBy", updatedBy)
	ops.AppendSet("/updatedByName", updatedByName)
	ops.AppendSet("/deleted", doc.GetDeleted())
	_, err := cc.PatchItem(ctx, azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String()), doc.GetDocID().String(), ops, nil)

	return err
}
