package graph

import (
	"context"

	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
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

func (s *graphService) NewDeviceDocFromGraph(
	graphObj msgraphmodels.Deviceable,
) *DeviceDoc {
	if graphObj == nil {
		return nil
	}
	doc := &DeviceDoc{}
	s.init(&doc.GraphDoc, graphObj, kmsdoc.DocTypeExtNameDevice)
	doc.DeviceID, _ = uuid.Parse(utils.NilToDefault(graphObj.GetId()))
	doc.AccountEnabled = utils.NilToDefault(graphObj.GetAccountEnabled())
	doc.OperatingSystem = graphObj.GetOperatingSystem()
	doc.OperatingSystemVersion = graphObj.GetOperatingSystemVersion()
	doc.TrustType = graphObj.GetTrustType()
	doc.MDMAppID = graphObj.GetMdmAppId()
	doc.IsCompliant = graphObj.GetIsCompliant()
	return doc
}

func (doc *DeviceDoc) Persist(ctx context.Context) error {
	return kmsdoc.AzCosmosUpsert(ctx, doc.service.AzCosmosContainerClient(), doc)
}
