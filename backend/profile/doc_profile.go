package profile

import (
	"fmt"

	gmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
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
	MsGraphOdataTypeDevice:           models.NamespaceKindDevice,
	MsGraphOdataTypeUser:             models.NamespaceKindUser,
	MsGraphOdataTypeGroup:            models.NamespaceKindGroup,
	MsGraphOdataTypeApplication:      models.NamespaceKindApplication,
	MsGraphOdataTypeServicePrincipal: models.NamespaceKindServicePrincipal,
}

type ProfileDoc struct {
	kmsdoc.BaseDoc

	ProfileType            models.NamespaceKind `json:"profileType"`
	OdataType              MsGraphOdataType     `json:"@odata.type"`
	DispalyName            *string              `json:"displayName,omitempty"`            // all
	AppID                  *string              `json:"appId,omitempty"`                  // application, service-principal
	DeviceID               *string              `json:"deviceId,omitempty"`               // device
	AccountEnabled         *bool                `json:"accountEnabled,omitempty"`         // device
	OperatingSystem        *string              `json:"operatingSystem,omitempty"`        // device
	OperatingSystemVersion *string              `json:"operatingSystemVersion,omitempty"` // device
	TrustType              *string              `json:"trustType,omitempty"`              // device
	MDMAppID               *string              `json:"mdmAppId,omitempty"`               // device
	IsCompliant            *bool                `json:"isCompliant,omitempty"`            // device
	UserPrincipalName      *string              `json:"userPrincipalName,omitempty"`      // user
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

	id := common.UUIDIdentifierFromStringPtr(dirObj.GetId())
	if dirObjUuid, isUuid := id.TryGetUUID(); !isUuid || dirObjUuid.Version() != 4 {
		return fmt.Errorf("invalid graph object id from api: %s", id.String())
	}
	d.ID = common.NewIdentifierWithKind(models.ResourceKindMsGraph, id)

	switch dirObj := dirObj.(type) {
	case gmodels.Deviceable:
		d.DispalyName = dirObj.GetDisplayName()
		d.DeviceID = dirObj.GetDeviceId()
		d.AccountEnabled = dirObj.GetAccountEnabled()
		d.OperatingSystem = dirObj.GetOperatingSystem()
		d.OperatingSystemVersion = dirObj.GetOperatingSystemVersion()
		d.TrustType = dirObj.GetTrustType()
		d.MDMAppID = dirObj.GetMdmAppId()
		d.IsCompliant = dirObj.GetIsCompliant()
	case gmodels.Userable:
		d.DispalyName = dirObj.GetDisplayName()
		d.UserPrincipalName = dirObj.GetUserPrincipalName()
	case gmodels.Groupable:
		d.DispalyName = dirObj.GetDisplayName()
	case gmodels.Applicationable:
		d.DispalyName = dirObj.GetDisplayName()
		d.AppID = dirObj.GetAppId()
	case gmodels.ServicePrincipalable:
		d.DispalyName = dirObj.GetDisplayName()
		d.AppID = dirObj.GetAppId()
	}
	return nil
}

func (d *ProfileDoc) populateRef(dst *models.ProfileRefComposed) bool {
	if ok := d.BaseDoc.PopulateResourceRef(&dst.ResourceRef); !ok {
		return ok
	}
	dst.Type = d.ProfileType
	dst.DisplayName = utils.NilToDefault(d.DispalyName)
	return true
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
