package admin

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type NsRelStatus string

const (
	NsRelStatusUnknown  NsRelStatus = "unknown"
	NsRelStatusVerified NsRelStatus = "verified"
	NsRelStatusError    NsRelStatus = "error"
	NsRelStatusLink     NsRelStatus = "link"
)

type NsRelDocNamespaces struct {
	Device           *uuid.UUID `json:"device,omitempty"`           // device object id
	Application      *uuid.UUID `json:"application,omitempty"`      // app object id
	ServicePrincipal *uuid.UUID `json:"servicePrincipal,omitempty"` // service principal object id
}

type NsRelDocAttributes struct {
	DeviceID *uuid.UUID `json:"deviceId,omitempty"` // device id associated with the device, not object id of the device directory object
	AppID    *uuid.UUID `json:"appId,omitempty"`    // app id associated with the application, not object id of the app directory object
}

type NsRelDoc struct {
	kmsdoc.BaseDoc
	SourceNamespaceID uuid.UUID          `json:"sourceNamespaceId"`
	LinkedNamespaces  NsRelDocNamespaces `json:"linked"`
	Attributes        NsRelDocAttributes `json:"attributes"`
	Status            NsRelStatus        `json:"status"`
	StatusMessage     string             `json:"statusMessage"`
}

func patchNsRelDocSourceNamespaceID(ops *azcosmos.PatchOperations, doc *NsRelDoc) {
	ops.AppendSet("/sourceNamespaceId", doc.SourceNamespaceID)
}

func patchNsRelDocLinkedNamespacesDevice(ops *azcosmos.PatchOperations, doc *NsRelDoc) {
	ops.AppendSet("/linked/device", doc.LinkedNamespaces.Device)
	ops.AppendSet("/attributes/deviceId", doc.Attributes.DeviceID)
}

func patchNsRelDocLinkedNamespacesApplication(ops *azcosmos.PatchOperations, doc *NsRelDoc) {
	ops.AppendSet("/linked/application", doc.LinkedNamespaces.Application)
	ops.AppendSet("/attributes/appId", doc.Attributes.AppID)
}

func patchNsRelDocLinkedNamespacesServicePrincipal(ops *azcosmos.PatchOperations, doc *NsRelDoc) {
	ops.AppendSet("/linked/servicePrincipal", doc.LinkedNamespaces.ServicePrincipal)
}

/*
func (s *adminServer) queryNsRelHasPermission(ctx context.Context, namespaceID uuid.UUID, permissionKey NamespacePermissionKey) ([]*NsRelDoc, error) {
	partitionKey := azcosmos.NewPartitionKeyString(namespaceID.String())
	pager := s.AzCosmosContainerClient().NewQueryItemsPager(`SELECT c.id,c.displayName FROM c
WHERE c.namespaceId = @namespaceId
  AND c.type = @type
  AND c.`+string(permissionKey),
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: namespaceID.String()},
				{Name: "@type", Value: kmsdoc.DocTypeNameNamespaceRelation},
			},
		})
	return PagerToList[NsRelDoc](ctx, pager)
}
*/

func (s *adminServer) readNsRelWithFollow(ctx context.Context, nsID uuid.UUID, relID uuid.UUID, follow bool) (*NsRelDoc, error) {
	doc := new(NsRelDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.AzCosmosContainerClient(), nsID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, relID), doc)
	if !follow || err == nil || doc.NamespaceID == nsID {
		return doc, err
	} else {
		return s.readNsRelWithFollow(ctx, doc.NamespaceID, relID, false)
	}
}

func (s *adminServer) readNsRel(ctx context.Context, nsID uuid.UUID, relID uuid.UUID) (*NsRelDoc, error) {
	return s.readNsRelWithFollow(ctx, nsID, relID, true)
}

/*
func (s *adminServer) patchNsRelStatus(c *gin.Context, doc *NsRelDoc, status NsRelStatus, statusMessage string) error {
	_, err := kmsdoc.AzCosmosPatch(c, s.AzCosmosContainerClient(), doc.NamespaceID,
		doc.ID, func(t time.Time) *azcosmos.PatchOperations {
			ops := azcosmos.PatchOperations{}
			ops.AppendSet("/status", status)
			ops.AppendSet("/statusMessage", statusMessage)
			return &ops
		})
	if err != nil {
		return err
	}
	doc.Status = status
	doc.StatusMessage = statusMessage
	return nil
}

/*
func (s *adminServer) putNsRelShadow(c *gin.Context, doc *NsRelDoc, targetNsID uuid.UUID) error {
	return kmsdoc.AzCosmosUpsert(c, s.AzCosmosContainerClient(), &NsRelDoc{
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: targetNsID,
			ID:          doc.ID,
		},
		Status:            NsRelStatusLink,
		SourceNamespaceID: doc.NamespaceID,
	})

}

func (s *adminServer) patchNsRelLinkedNamespaces(c *gin.Context, doc *NsRelDoc, keys ...string) error {
	_, err := kmsdoc.AzCosmosPatch(c, s.AzCosmosContainerClient(), doc.NamespaceID,
		doc.ID, func(t time.Time) *azcosmos.PatchOperations {
			ops := azcosmos.PatchOperations{}
			for _, key := range keys {
				ops.AppendSet(fmt.Sprintf("/linked/%s", key), doc.LinkedNamespaces[key])
			}
			return &ops
		})
	if err != nil {
		return err
	}
	return nil
}


func (s *adminServer) hasAllowEnrollDeviceCertificatePermission(ctx context.Context, namespaceID uuid.UUID, objectID uuid.UUID) (bool, error) {
	doc, err := s.getNsRel(ctx, namespaceID, objectID)
	if err != nil {
		if common.IsAzNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return doc.AllowEnrollDeviceCertificate != nil && *doc.AllowEnrollDeviceCertificate, nil
}
*/
