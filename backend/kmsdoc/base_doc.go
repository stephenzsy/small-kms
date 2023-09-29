package kmsdoc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/auth"
)

type KmsDocType byte

const (
	DocTypeUnknown KmsDocType = 90 // Z

	DocTypeCert                KmsDocType = 67 // C
	DocTypeDeviceLink          KmsDocType = 68 // D
	DocTypeLatestCertForPolicy KmsDocType = 76 // L
	DocTypeDirectoryObject     KmsDocType = 79 // O
	DocTypePolicy              KmsDocType = 80 // P, deprecated
	DocTypeNamespaceRelation   KmsDocType = 82 // R, deprecated
	DocTypePolicyState         KmsDocType = 83 // S, deprecated
	DocTypeCertTemplate        KmsDocType = 84 // T
)

type KmsDocTypeName string

const (
	DocTypeNameUnknown KmsDocTypeName = "unknown"

	DocTypeNameCert                KmsDocTypeName = "cert"
	DocTypeNameDeviceLink          KmsDocTypeName = "device-link"
	DocTypeNameLatestCertForPolicy KmsDocTypeName = "cert-latest"
	DocTypeNameDirectoryObject     KmsDocTypeName = "directory-object"
	DocTypeNamePolicy              KmsDocTypeName = "policy"
	DocTypeNameNamespaceRelation   KmsDocTypeName = "namespace-relation"
	DocTypeNamePolicyState         KmsDocTypeName = "policy-state"
	DocTypeNameCertTemplate        KmsDocTypeName = "cert-template"
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

func (k *KmsDocID) String() string {
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
	switch k.typeByte {
	case DocTypeCert,
		DocTypePolicy,
		DocTypePolicyState,
		DocTypeLatestCertForPolicy:
		// accept
	default:
		k.typeByte = DocTypeUnknown
	}
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
}

func GetBaseDocQueryColumns(prefix string) string {
	return fmt.Sprintf("%s.id,%s.namespaceId,%s.updated,%s.updatedBy,%s.deleted", prefix, prefix, prefix, prefix, prefix)
}

type KmsDocument interface {
	GetNamespaceID() uuid.UUID
	StampUpdatedWithAuth(c *gin.Context)
	GetUUID() uuid.UUID

	fillTypeName()
}

func (k *BaseDoc) GetUUID() uuid.UUID {
	return k.ID.GetUUID()
}

func (doc *BaseDoc) GetNamespaceID() uuid.UUID {
	return doc.NamespaceID
}

func (doc *BaseDoc) StampUpdated(callerId string, callerName string) {
	doc.Updated = time.Now()
	doc.UpdatedBy = callerId
	doc.UpdatedByName = callerName
}

func (doc *BaseDoc) StampUpdatedWithAuth(c *gin.Context) {
	doc.StampUpdated(auth.CallerPrincipalId(c).String(), auth.CallerPrincipalName(c))
}

var docTypeNameMap = map[KmsDocType]KmsDocTypeName{
	DocTypeCert:                DocTypeNameCert,
	DocTypeDeviceLink:          DocTypeNameDeviceLink,
	DocTypeLatestCertForPolicy: DocTypeNameLatestCertForPolicy,
	DocTypeDirectoryObject:     DocTypeNameDirectoryObject,
	DocTypePolicy:              DocTypeNamePolicy,
	DocTypeNamespaceRelation:   DocTypeNameNamespaceRelation,
	DocTypePolicyState:         DocTypeNamePolicyState,
	DocTypeCertTemplate:        DocTypeNameCertTemplate,
}

func (doc *BaseDoc) fillTypeName() {
	doc.TypeName = docTypeNameMap[doc.ID.typeByte]
}

func AzCosmosUpsert[D KmsDocument](ctx *gin.Context, cc *azcosmos.ContainerClient, doc D) error {
	doc.StampUpdatedWithAuth(ctx)
	doc.fillTypeName()
	content, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	_, err = cc.UpsertItem(ctx, azcosmos.NewPartitionKeyString(doc.GetNamespaceID().String()), content, nil)
	return err
}

func AzCosmosRead[D KmsDocument](ctx context.Context, cc *azcosmos.ContainerClient, namespaceID uuid.UUID, docID KmsDocID, target D) error {
	resp, err := cc.ReadItem(ctx, azcosmos.NewPartitionKeyString(namespaceID.String()), docID.String(), nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp.Value, target)
}

func AzCosmosDelete(ctx *gin.Context, cc *azcosmos.ContainerClient, namespaceID uuid.UUID, docID KmsDocID, purge bool) (err error) {
	if purge {
		_, err = cc.DeleteItem(ctx, azcosmos.NewPartitionKeyString(namespaceID.String()), docID.String(), nil)
	} else {
		ops := azcosmos.PatchOperations{}
		tsStr := time.Now().UTC().Format(time.RFC3339)
		ops.AppendSet("/deleted", tsStr)
		ops.AppendSet("/updated", tsStr)
		ops.AppendSet("/updatedBy", auth.CallerPrincipalId(ctx).String())
		ops.AppendSet("/updatedByName", auth.CallerPrincipalName(ctx))
		_, err = cc.PatchItem(ctx, azcosmos.NewPartitionKeyString(namespaceID.String()), docID.String(), ops, nil)
	}
	return
}
