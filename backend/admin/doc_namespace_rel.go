package admin

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
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
