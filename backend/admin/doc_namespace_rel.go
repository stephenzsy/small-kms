package admin

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/gin-gonic/gin"
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

type NsRelDoc struct {
	kmsdoc.BaseDoc
	SourceNamespaceID uuid.UUID            `json:"sourceNamespaceId"`
	LinkedNamespaces  map[string]uuid.UUID `json:"linked"`
	Status            NsRelStatus          `json:"status"`
	StatusMessage     string               `json:"statusMessage"`
}

/*
func (s *adminServer) queryNsRelHasPermission(ctx context.Context, namespaceID uuid.UUID, permissionKey NamespacePermissionKey) ([]*NsRelDoc, error) {
	partitionKey := azcosmos.NewPartitionKeyString(namespaceID.String())
	pager := s.azCosmosContainerClientCerts.NewQueryItemsPager(`SELECT c.id,c.displayName FROM c
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
	err := kmsdoc.AzCosmosRead(ctx, s.azCosmosContainerClientCerts, nsID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, relID), doc)
	if !follow {
		return doc, err
	}
	if err != nil {
		return doc, err
	}
	if doc.NamespaceID == nsID {
		return doc, err
	} else {
		return s.readNsRelWithFollow(ctx, doc.NamespaceID, relID, false)
	}
}

func (s *adminServer) readNsRel(ctx context.Context, nsID uuid.UUID, relID uuid.UUID) (*NsRelDoc, error) {
	return s.readNsRelWithFollow(ctx, nsID, relID, true)
}

func (s *adminServer) patchNsRelStatus(c *gin.Context, doc *NsRelDoc, status NsRelStatus, statusMessage string) error {
	_, err := kmsdoc.AzCosmosPatch(c, s.azCosmosContainerClientCerts, doc.NamespaceID,
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

func (s *adminServer) putNsRelShadow(c *gin.Context, doc *NsRelDoc, targetNsID uuid.UUID) error {
	return kmsdoc.AzCosmosUpsert(c, s.azCosmosContainerClientCerts, &NsRelDoc{
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: targetNsID,
			ID:          doc.ID,
		},
		Status:            NsRelStatusLink,
		SourceNamespaceID: doc.NamespaceID,
	})

}

func (s *adminServer) patchNsRelLinkedNamespaces(c *gin.Context, doc *NsRelDoc, keys ...string) error {
	_, err := kmsdoc.AzCosmosPatch(c, s.azCosmosContainerClientCerts, doc.NamespaceID,
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

/*
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
