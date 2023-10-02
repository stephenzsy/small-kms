package graph

import (
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// fields for application graph profile
type ApplicationDoc struct {
	GraphDoc

	// client ID of the application
	AppID uuid.UUID `json:"appId"`
}

func (doc *ApplicationDoc) init(
	tenantID uuid.UUID,
	graphObj GraphProfileable,
	_ MsGraphOdataType,
) {
	if graphObj == nil {
		return
	}

	doc.GraphDoc.init(tenantID, graphObj, MsGraphOdataTypeApplication)
	if graphObj, ok := graphObj.(msgraphmodels.ServicePrincipalable); ok {
		doc.AppID, _ = uuid.Parse(utils.NilToDefault(graphObj.GetAppId()))
	}
}
