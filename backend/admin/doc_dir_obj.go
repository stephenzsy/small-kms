package admin

import (
	"context"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/graph"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type DirectoryObjectDocDevice struct {
	DeviceID               uuid.UUID `json:"deviceId,omitempty"`
	OperatingSystem        *string   `json:"operatingSystem,omitempty"`
	OperatingSystemVersion *string   `json:"operatingSystemVersion,omitempty"`
	DeviceOwnership        *string   `json:"deviceOwnership,omitempty"`
	IsCompliant            *bool     `json:"isCompliant,omitempty"`
}

type DirectoryObjectDocServicePrincipal struct {
	ServicePrincipalType string `json:"servicePrincipalType"`
}
type DirectoryObjectDocApplication struct {
	AppID string `json:"appId,omitempty"`
}

type DirectoryObjectDoc struct {
	kmsdoc.BaseDoc
	OdataType         string                              `json:"odType"`
	DisplayName       string                              `json:"displayName"`
	UserPrincipalName *string                             `json:"userPrincipalName,omitempty"`
	Application       *DirectoryObjectDocApplication      `json:"application,omitempty"`
	Device            *DirectoryObjectDocDevice           `json:"device,omitempty"`
	ServicePrincipal  *DirectoryObjectDocServicePrincipal `json:"servicePrincipal,omitempty"`
}

func (s *adminServer) getDirectoryObjectDoc(ctx context.Context, objectID uuid.UUID) (*DirectoryObjectDoc, error) {
	doc := new(DirectoryObjectDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.AzCosmosContainerClient(), common.WellKnownID_TenantDirectory,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeDirectoryObject, objectID), doc)
	return doc, err
}

func toNsType(odataType graph.MsGraphOdataType) NamespaceTypeShortName {
	switch odataType {
	case graph.MsGraphOdataTypeDevice:
		return NSTypeDevice
	case graph.MsGraphOdataTypeGroup:
		return NSTypeGroup
	case graph.MsGraphOdataTypeUser:
		return NSTypeUser
	case graph.MsGraphOdataTypeApplication:
		return NSTypeApplication
	case graph.MsGraphOdataTypeServicePrincipal:
		return NSTypeServicePrincipal
	}
	return NSTypeUnknown
}

func newNamespaceInfoFromProfileDoc(doc graph.GraphProfileDocument) *NamespaceInfo {
	if doc == nil {
		return nil
	}
	p := new(NamespaceInfo)
	profileDocPopulateRefWithMetadata(doc, &p.Ref, toNsType(doc.GetOdataType()))
	return p
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
		if item.ServicePrincipal != nil {
			ref.ServicePrincipalType = ToPtr(item.ServicePrincipal.ServicePrincipalType)
		}
	case "#microsoft.graph.device":
		if item.Device != nil {
			ref.DeviceID = ToPtr(item.Device.DeviceID.String())
			ref.DeviceOwnership = item.Device.DeviceOwnership
			ref.IsCompliant = item.Device.IsCompliant
			ref.OperatingSystem = item.Device.OperatingSystem
			ref.OperatingSystemVersion = item.Device.OperatingSystemVersion
		}
	case string(NamespaceTypeMsGraphApplication):
		if item.Application != nil {
			ref.AppID = ToPtr(item.Application.AppID)
		}
	}
}

func newDirectoryObjectDocFromApplicationable() *DirectoryObjectDoc {
	return &DirectoryObjectDoc{
		OdataType: "#microsoft.graph.application",
	}
}
