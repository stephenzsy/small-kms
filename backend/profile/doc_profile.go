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

var supportedMsGraphOdataTypeToDocNsID = map[MsGraphOdataType]models.ProfileType{
	MsGraphOdataTypeDevice:           models.ProfileTypeDevice,
	MsGraphOdataTypeUser:             models.ProfileTypeUser,
	MsGraphOdataTypeGroup:            models.ProfileTypeGroup,
	MsGraphOdataTypeApplication:      models.ProfileTypeApplication,
	MsGraphOdataTypeServicePrincipal: models.ProfileTypeServicePrincipal,
}

type ProfileDoc struct {
	kmsdoc.BaseDoc

	OdataType              MsGraphOdataType `json:"@odata.type"`
	DispalyName            *string          `json:"displayName,omitempty"`            // all
	AppID                  *string          `json:"appId,omitempty"`                  // application, service-principal
	DeviceID               *string          `json:"deviceId,omitempty"`               // device
	AccountEnabled         *bool            `json:"accountEnabled,omitempty"`         // device
	OperatingSystem        *string          `json:"operatingSystem,omitempty"`        // device
	OperatingSystemVersion *string          `json:"operatingSystemVersion,omitempty"` // device
	TrustType              *string          `json:"trustType,omitempty"`              // device
	MDMAppID               *string          `json:"mdmAppId,omitempty"`               // device
	IsCompliant            *bool            `json:"isCompliant,omitempty"`            // device
	UserPrincipalName      *string          `json:"userPrincipalName,omitempty"`      // user
}

func getProfileDocKey(objectIdentifier models.Identifier) (string, string) {
	return string(kmsdoc.DocNsTypeTenant), objectIdentifier.String()
}

func (d *ProfileDoc) init(dirObj gmodels.DirectoryObjectable) error {
	if dirObj == nil {
		return fmt.Errorf("nil directory object from graph api")
	}

	d.SchemaVersion = 1
	d.NsType = kmsdoc.DocNsTypeTenant
	d.DocType = kmsdoc.DocTypeProfile
	odataType := dirObj.GetOdataType()
	if odataType == nil {
		return fmt.Errorf("nil odata type from graph api")
	}
	if _, ok := supportedMsGraphOdataTypeToDocNsID[MsGraphOdataType(*odataType)]; ok {
		d.OdataType = MsGraphOdataType(*odataType)
	} else {
		return fmt.Errorf("%w:unsupported odata type from graph api: %s", common.ErrStatusBadRequest, *odataType)
	}
	id := dirObj.GetId()
	if id != nil {
		return fmt.Errorf("missing id from graph api")
	}
	d.DocID = models.IdentifierFromString(*id)

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

func (d *ProfileDoc) toModel() (p *models.Profile) {
	if d == nil {
		return nil
	}
	p = &models.Profile{
		Identifier: d.DocID,
		Metadata: models.ResourceMetadata{
			Updated:   utils.ToPtr(d.Updated),
			UpdatedBy: utils.ToPtr(d.UpdatedBy),
			Deleted:   d.Deleted,
		},
		Type: supportedMsGraphOdataTypeToDocNsID[d.OdataType],
	}
	if d.DispalyName != nil {
		p.DisplayName = *d.DispalyName
	}
	return p
}
