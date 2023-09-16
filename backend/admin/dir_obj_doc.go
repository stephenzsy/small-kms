package admin

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type DirectoryObjectDoc struct {
	kmsdoc.BaseDoc
	OdataType              string  `json:"odType"`
	DisplayName            string  `json:"displayName"`
	UserPrincipalName      *string `json:"userPrincipalName,omitempty"`
	ServicePrincipalType   *string `json:"servicePrincipalType,omitempty"`
	DeviceID               *string `json:"deviceId,omitempty"`
	OperatingSystem        *string `json:"operatingSystem,omitempty"`
	OperatingSystemVersion *string `json:"operatingSystemVersion,omitempty"`
	DeviceOwnership        *string `json:"deviceOwnership,omitempty"`
	IsCompliant            *bool   `json:"isCompliant,omitempty"`
}

func (s *adminServer) GetDirectoryObjectDoc(ctx context.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	doc := new(DirectoryObjectDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.azCosmosContainerClientCerts, directoryID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeDirectoryObject, objectID), doc)
	return doc, err
}

func (s *adminServer) ListDirectoryObjectByType(ctx context.Context, nsType NamespaceType) ([]*DirectoryObjectDoc, error) {
	switch nsType {
	case NamespaceTypeMsGraphUser,
		NamespaceTypeMsGraphGroup,
		NamespaceTypeMsGraphDevice,
		NamespaceTypeMsGraphServicePrincipal:
	default:
		return nil, fmt.Errorf("namespace type not supported")
	}
	partitionKey := azcosmos.NewPartitionKeyString(directoryID.String())
	pager := s.azCosmosContainerClientCerts.NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.odType,c.displayName FROM c
WHERE c.namespaceId = @namespaceId
  AND c.odType = @odType`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: directoryID.String()},
				{Name: "@odType", Value: nsType},
			},
		})

	return PagerToList[DirectoryObjectDoc](ctx, pager)
}

func (item *DirectoryObjectDoc) PopulateNamespaceRef(ref *NamespaceRef) {
	ref.NamespaceID = directoryID
	ref.ID = item.ID.GetUUID()
	ref.DisplayName = item.DisplayName
	ref.ObjectType = NamespaceType(item.OdataType)
	ref.Updated = item.Updated
	ref.UpdatedBy = item.UpdatedBy
}

func (item *DirectoryObjectDoc) PopulateNamespaceProfile(ref *NamespaceProfile) {
	ref.NamespaceID = directoryID
	ref.ID = item.ID.GetUUID()
	ref.DisplayName = item.DisplayName
	ref.ObjectType = NamespaceType(item.OdataType)
	ref.Updated = item.Updated
	ref.UpdatedBy = item.UpdatedBy

	switch item.OdataType {
	case "#microsoft.graph.user":
		ref.UserPrincipalName = item.UserPrincipalName
	case "#microsoft.graph.servicePrincipal":
		ref.ServicePrincipalType = item.ServicePrincipalType
	case "#microsoft.graph.device":
		ref.DeviceID = item.DeviceID
		ref.DeviceOwnership = item.DeviceOwnership
		ref.IsCompliant = item.IsCompliant
		ref.OperatingSystem = item.OperatingSystem
		ref.OperatingSystemVersion = item.OperatingSystemVersion
	}
}
