package graph

import (
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// fields for device service principal profile
type ServicePrincipalDoc struct {
	GraphDoc

	// client ID of the service principal
	AppID uuid.UUID `json:"appId"`
}

func (doc *ServicePrincipalDoc) init(
	tenantID uuid.UUID,
	graphObj GraphProfileable,
	_ MsGraphOdataType,
) {
	if graphObj == nil {
		return
	}

	doc.GraphDoc.init(tenantID, graphObj, MsGraphOdataTypeServicePrincipal)
	if graphObj, ok := graphObj.(msgraphmodels.ServicePrincipalable); ok {
		doc.AppID, _ = uuid.Parse(utils.NilToDefault(graphObj.GetAppId()))
	}
}
