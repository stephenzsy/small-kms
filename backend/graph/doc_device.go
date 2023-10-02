package graph

import (
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// fields for device graph profile
type DeviceDoc struct {
	GraphDoc

	DeviceID       uuid.UUID `json:"deviceId"`
	AccountEnabled bool      `json:"accountEnabled"`

	OperatingSystem        *string `json:"operatingSystem,omitempty"`
	OperatingSystemVersion *string `json:"operatingSystemVersion,omitempty"`
	TrustType              *string `json:"trustType,omitempty"`
	MDMAppID               *string `json:"mdmAppId,omitempty"`
	IsCompliant            *bool   `json:"isCompliant,omitempty"`
}

func GetProfileGraphSelectDeviceDoc() (r []string) {
	r = append(r, GetProfileGraphSelectGraphDoc()...)
	r = append(r, "deviceId", "accountEnabled", "operatingSystem", "operatingSystemVersion", "trustType", "mdmAppId", "isCompliant")
	return r
}

func (doc *DeviceDoc) init(
	tenantID uuid.UUID,
	graphObj GraphProfileable,
	_ MsGraphOdataType,
) {
	if graphObj == nil {
		return
	}

	doc.GraphDoc.init(tenantID, graphObj, MsGraphOdataTypeDevice)
	if graphObj, ok := graphObj.(msgraphmodels.Deviceable); ok {
		doc.DeviceID, _ = uuid.Parse(utils.NilToDefault(graphObj.GetDeviceId()))
		doc.AccountEnabled = utils.NilToDefault(graphObj.GetAccountEnabled())
		doc.OperatingSystem = graphObj.GetOperatingSystem()
		doc.OperatingSystemVersion = graphObj.GetOperatingSystemVersion()
		doc.TrustType = graphObj.GetTrustType()
		doc.MDMAppID = graphObj.GetMdmAppId()
		doc.IsCompliant = graphObj.GetIsCompliant()
	}
}
