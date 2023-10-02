package graph

import (
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// fields for user graph profile
type UserDoc struct {
	GraphDoc

	UserPrincipalName string `json:"userPrincipalName"`
}

func (doc *UserDoc) init(
	tenantID uuid.UUID,
	graphObj GraphProfileable,
	_ MsGraphOdataType,
) {
	if graphObj == nil {
		return
	}

	doc.GraphDoc.init(tenantID, graphObj, MsGraphOdataTypeUser)
	if graphObj, ok := graphObj.(msgraphmodels.Userable); ok {
		doc.UserPrincipalName = utils.NilToDefault(graphObj.GetUserPrincipalName())
	}
}
