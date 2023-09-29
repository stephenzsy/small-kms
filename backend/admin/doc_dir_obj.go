package admin

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type DirectoryObjectDocDeviceSection struct {
	DeviceID               uuid.UUID `json:"deviceId,omitempty"`
	OperatingSystem        *string   `json:"operatingSystem,omitempty"`
	OperatingSystemVersion *string   `json:"operatingSystemVersion,omitempty"`
	DeviceOwnership        *string   `json:"deviceOwnership,omitempty"`
	IsCompliant            *bool     `json:"isCompliant,omitempty"`

	// custom link info, not from ms graph
	LinkedApplicationObjectID      *uuid.UUID `json:"linkedApplicationObjectId,omitempty"`
	LinkedApplicationClientID      *uuid.UUID `json:"linkedApplicationClientId,omitempty"`
	LinkedServicePrincipalObjectID *uuid.UUID `json:"linkedServicePrincipalObjectID,omitempty"`
}

type DirectoryObjectDoc struct {
	kmsdoc.BaseDoc
	OdataType            string                           `json:"odType"`
	DisplayName          string                           `json:"displayName"`
	UserPrincipalName    *string                          `json:"userPrincipalName,omitempty"`
	ServicePrincipalType *string                          `json:"servicePrincipalType,omitempty"`
	AppID                *string                          `json:"appId,omitempty"`
	Device               *DirectoryObjectDocDeviceSection `json:"device,omitempty"`
}

func (s *adminServer) getDirectoryObjectDoc(ctx context.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	doc := new(DirectoryObjectDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.azCosmosContainerClientCerts, common.WellKnownID_TenantDirectory,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeDirectoryObject, objectID), doc)
	return doc, err
}

func (s *adminServer) listDirectoryObjectByType(ctx context.Context, odType string) ([]*DirectoryObjectDoc, error) {
	partitionKey := azcosmos.NewPartitionKeyString(common.WellKnownID_TenantDirectory.String())
	pager := s.azCosmosContainerClientCerts.NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.odType,c.displayName FROM c
WHERE c.namespaceId = @namespaceId
  AND c.odType = @odType`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: common.WellKnownID_TenantDirectory.String()},
				{Name: "@odType", Value: odType},
			},
		})

	return PagerToList[DirectoryObjectDoc](ctx, pager)
}

func toNsType(odataType string) NamespaceTypeShortName {
	switch odataType {
	case string(NamespaceTypeMsGraphServicePrincipal):
		return NSTypeServicePrincipal
	case string(NamespaceTypeMsGraphGroup):
		return NSTypeGroup
	case string(NamespaceTypeMsGraphDevice):
		return NSTypeDevice
	case string(NamespaceTypeMsGraphUser):
		return NSTypeUser
	case string(NamespaceTypeMsGraphApplication):
		return NSTypeApplication
	}
	return NSTypeUnknown
}

func (item *DirectoryObjectDoc) toNamespaceInfo() *NamespaceInfo {
	if item == nil {
		return nil
	}
	r := new(NamespaceInfo)
	baseDocPopulateRefWithMetadata(&item.BaseDoc, &r.Ref, toNsType(item.OdataType))

	r.Ref.Metadata[RefPropertyKeyDisplayName] = item.DisplayName
	switch r.Ref.NamespaceType {
	}
	return r
}

func (item *DirectoryObjectDoc) PopulateNamespaceRef(ref *NamespaceRef) {
	ref.NamespaceID = common.WellKnownID_TenantDirectory
	ref.ID = item.ID.GetUUID()
	ref.DisplayName = item.DisplayName
	ref.ObjectType = NamespaceType(item.OdataType)
	ref.Updated = item.Updated
	ref.UpdatedBy = item.UpdatedBy
}

func (item *DirectoryObjectDoc) PopulateNamespaceProfile(ref *NamespaceProfile) {
	ref.NamespaceID = common.WellKnownID_TenantDirectory
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
		if item.Device != nil {
			ref.DeviceID = ToPtr(item.Device.DeviceID.String())
			ref.DeviceOwnership = item.Device.DeviceOwnership
			ref.IsCompliant = item.Device.IsCompliant
			ref.OperatingSystem = item.Device.OperatingSystem
			ref.OperatingSystemVersion = item.Device.OperatingSystemVersion
		}
	case string(NamespaceTypeMsGraphApplication):
		ref.AppID = item.AppID
	}
}
