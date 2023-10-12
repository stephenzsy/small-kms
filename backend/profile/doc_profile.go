package profile

import (
	"fmt"

	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type MsGraphOdataType string

const (
	MsGraphOdataTypeAny MsGraphOdataType = ""

	MsGraphOdataTypeDevice           MsGraphOdataType = "#microsoft.graph.device"
	MsGraphOdataTypeUser             MsGraphOdataType = "#microsoft.graph.user"
	MsGraphOdataTypeGroup            MsGraphOdataType = "#microsoft.graph.group"
	MsGraphOdataTypeApplication      MsGraphOdataType = "#microsoft.graph.application"
	MsGraphOdataTypeServicePrincipal MsGraphOdataType = "#microsoft.graph.servicePrincipal"
)

var supportedMsGraphOdataTypeToDocNsID = map[MsGraphOdataType]models.NamespaceKind{
	MsGraphOdataTypeDevice:           shared.NamespaceKindDevice,
	MsGraphOdataTypeUser:             shared.NamespaceKindUser,
	MsGraphOdataTypeGroup:            shared.NamespaceKindGroup,
	MsGraphOdataTypeApplication:      shared.NamespaceKindApplication,
	MsGraphOdataTypeServicePrincipal: shared.NamespaceKindServicePrincipal,
}

type ProfileDocGraphData struct {
	DispalyName            *string `json:"displayName,omitempty"`            // all
	AppID                  *string `json:"appId,omitempty"`                  // application, service-principal
	DeviceID               *string `json:"deviceId,omitempty"`               // device
	AccountEnabled         *bool   `json:"accountEnabled,omitempty"`         // device
	OperatingSystem        *string `json:"operatingSystem,omitempty"`        // device
	OperatingSystemVersion *string `json:"operatingSystemVersion,omitempty"` // device
	TrustType              *string `json:"trustType,omitempty"`              // device
	MDMAppID               *string `json:"mdmAppId,omitempty"`               // device
	IsCompliant            *bool   `json:"isCompliant,omitempty"`            // device
	UserPrincipalName      *string `json:"userPrincipalName,omitempty"`      // user
}

type ProfileDoc struct {
	kmsdoc.BaseDoc

	ProfileType shared.NamespaceKind `json:"profileType"`

	GraphSyncCode string               `json:"graphSyncCode"` // field in the doc to indicate object is managed my this app
	Graph         *ProfileDocGraphData `json:"graph,omitempty"`
	IsAppManaged  *bool                `json:"isAppManaged"` // field in the doc to indicate object is managed my this app
	DispalyName   string               `json:"displayName"`
	OdataType     MsGraphOdataType     `json:"@odata.type"`
	IsBuiltIn     bool                 `json:"-"` // builtin does not persist in DB
}

func (d *ProfileDoc) init(dirObj gmodels.DirectoryObjectable) error {
	if dirObj == nil {
		return fmt.Errorf("nil directory object from graph api")
	}

	d.SchemaVersion = 1
	d.NamespaceID = docNsIDProfileTenant
	odataType := dirObj.GetOdataType()
	if odataType == nil {
		return fmt.Errorf("nil odata type from graph api")
	}
	if profileType, ok := supportedMsGraphOdataTypeToDocNsID[MsGraphOdataType(*odataType)]; ok {
		d.ProfileType = profileType
		d.OdataType = MsGraphOdataType(*odataType)
	} else {
		return fmt.Errorf("%w:unsupported odata type from graph api: %s", common.ErrStatusBadRequest, *odataType)
	}

	id := shared.UUIDIdentifierFromStringPtr(dirObj.GetId())
	if dirObjUuid, isUuid := id.TryGetUUID(); !isUuid || dirObjUuid.Version() != 4 {
		return fmt.Errorf("invalid graph object id from api: %s", id.String())
	}
	d.ID = shared.NewResourceIdentifier(shared.ResourceKindMsGraph, id)

	g := ProfileDocGraphData{}
	d.Graph = &g
	switch dirObj := dirObj.(type) {
	case gmodels.Deviceable:
		g.DispalyName = dirObj.GetDisplayName()
		d.DispalyName = utils.NilToDefault(g.DispalyName)
		g.DeviceID = dirObj.GetDeviceId()
		g.AccountEnabled = dirObj.GetAccountEnabled()
		g.OperatingSystem = dirObj.GetOperatingSystem()
		g.OperatingSystemVersion = dirObj.GetOperatingSystemVersion()
		g.TrustType = dirObj.GetTrustType()
		g.MDMAppID = dirObj.GetMdmAppId()
		g.IsCompliant = dirObj.GetIsCompliant()
	case gmodels.Userable:
		g.DispalyName = dirObj.GetDisplayName()
		d.DispalyName = utils.NilToDefault(g.DispalyName)
		g.UserPrincipalName = dirObj.GetUserPrincipalName()
	case gmodels.Groupable:
		g.DispalyName = dirObj.GetDisplayName()
		d.DispalyName = utils.NilToDefault(g.DispalyName)
	case gmodels.Applicationable:
		g.DispalyName = dirObj.GetDisplayName()
		d.DispalyName = utils.NilToDefault(g.DispalyName)
		g.AppID = dirObj.GetAppId()
	case gmodels.ServicePrincipalable:
		g.DispalyName = dirObj.GetDisplayName()
		d.DispalyName = utils.NilToDefault(g.DispalyName)
		g.AppID = dirObj.GetAppId()
	}
	return nil
}

func (d *ProfileDoc) populateRef(r *models.ProfileRefComposed) {
	if d == nil || r == nil {
		return
	}
	d.BaseDoc.PopulateResourceRef(&r.ResourceRef)
	r.Type = d.ProfileType
	r.DisplayName = d.DispalyName
	r.IsAppManaged = d.IsAppManaged
}

func (d *ProfileDoc) toModelRef() *models.ProfileRefComposed {
	if d == nil {
		return nil
	}
	p := models.ProfileComposed{}
	d.populateRef(&p)

	return &p
}

func (d *ProfileDoc) toModel() *models.ProfileComposed {
	return d.toModelRef()
}
