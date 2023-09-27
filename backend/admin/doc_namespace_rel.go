package admin

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type NsRelDoc struct {
	kmsdoc.BaseDoc
	AllowEnrollDeviceCertificate *bool  `json:"allowEnrollDeviceCertificate,omitempty"`
	DisplayName                  string `json:"displayName,omitempty"`
}

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

func (s *adminServer) getNsRel(ctx context.Context, namespaceID uuid.UUID, objectID uuid.UUID) (*NsRelDoc, error) {
	doc := new(NsRelDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.azCosmosContainerClientCerts, namespaceID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeNamespaceRelation, objectID), doc)
	return doc, err
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
