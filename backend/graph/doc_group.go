package graph

import (
	"github.com/google/uuid"
)

// fields for group graph profile
type GroupDoc struct {
	GraphDoc
}

func (doc *GroupDoc) init(
	tenantID uuid.UUID,
	graphObj GraphProfileable,
	_ MsGraphOdataType,
) {
	if graphObj == nil {
		return
	}

	doc.GraphDoc.init(tenantID, graphObj, MsGraphOdataTypeApplication)
}
