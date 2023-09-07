package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type DirectoryObjectDoc struct {
	kmsdoc.BaseDoc
	OdataType            string  `json:"odType"`
	DisplayName          string  `json:"displayName"`
	UserPrincipalName    *string `json:"userPrincipalName,omitempty"`
	ServicePrincipalType *string `json:"servicePrincipalType,omitempty"`
}

func (s *adminServer) GetDirectoryObjectDoc(ctx context.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	doc := new(DirectoryObjectDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.azCosmosContainerClientCerts, directoryID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeDirectoryObject, objectID), doc)
	return doc, err
}

func (s *adminServer) ListDirectoryObjectByType(ctx context.Context, nsType NamespaceType) (results []DirectoryObjectDoc, err error) {
	switch nsType {
	case NamespaceTypeMsGraphUser:
	case NamespaceTypeMsGraphServicePrincipal:
	default:
		return nil, fmt.Errorf("namespace type not supported")
	}
	partitionKey := azcosmos.NewPartitionKeyString(directoryID.String())
	pager := s.azCosmosContainerClientCerts.NewQueryItemsPager(`SELECT * FROM c
WHERE c.namespaceId = @namespaceId
  AND c.odType = @odType`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: directoryID.String()},
				{Name: "@odType", Value: nsType},
			},
		})

	for pager.More() {
		t, scanErr := pager.NextPage(ctx)
		if scanErr != nil {
			err = fmt.Errorf("faild to get list of certificates: %w", scanErr)
			return
		}
		for _, itemBytes := range t.Items {
			item := DirectoryObjectDoc{}
			if err = json.Unmarshal(itemBytes, &item); err != nil {
				err = fmt.Errorf("faild to serialize db entry: %w", err)
				return
			}
			results = append(results, item)
		}
	}
	return
}

func (item *DirectoryObjectDoc) PopulateNamespaceRef(ref *NamespaceRef) {
	ref.NamespaceID = directoryID
	ref.ID = item.ID.GetUUID()
	ref.DisplayName = item.DisplayName
	ref.ObjectType = NamespaceType(item.OdataType)
	ref.UserPrincipalName = item.UserPrincipalName
	ref.ServicePrincipalType = item.ServicePrincipalType
	ref.Updated = item.Updated
	ref.UpdatedBy = item.UpdatedBy
}
