package graph

import (
	ctx "context"

	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type MsGraphOdataType string

const (
	MsGraphOdataTypeNone MsGraphOdataType = ""

	MsGraphOdataTypeDevice           MsGraphOdataType = "#microsoft.graph.device"
	MsGraphOdataTypeUser             MsGraphOdataType = "#microsoft.graph.user"
	MsGraphOdataTypeGroup            MsGraphOdataType = "#microsoft.graph.group"
	MsGraphOdataTypeApplication      MsGraphOdataType = "#microsoft.graph.application"
	MsGraphOdataTypeServicePrincipal MsGraphOdataType = "#microsoft.graph.servicePrincipal"
)

// this kind of docs represent a graph object browsable in the app, it should not be used to persist data, always query live graph service
type GraphDoc struct {
	kmsdoc.BaseDoc

	OdataType   MsGraphOdataType `json:"@odata.type"`
	DisplayName string           `json:"displayName"`
}

func GetProfileGraphSelectGraphDoc() []string {
	return []string{"id", "displayName"}
}

func (doc *GraphDoc) GetDisplayName() string {
	return doc.DisplayName
}

type GraphProfileable interface {
	msgraphmodels.DirectoryObjectable
	GetDisplayName() *string
}

func (doc *GraphDoc) init(tenantID uuid.UUID, obj GraphProfileable, odataType MsGraphOdataType) {
	if obj == nil {
		return
	}
	oid, _ := uuid.Parse(utils.NilToDefault(obj.GetId()))
	doc.NamespaceID = tenantID
	doc.BaseDoc.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeMsGraphObject, oid)
	if obj.GetOdataType() != nil {
		doc.OdataType = MsGraphOdataType(*obj.GetOdataType())
	} else {
		doc.OdataType = odataType
	}
	doc.DisplayName = utils.NilToDefault(obj.GetDisplayName())
}

func (s *graphService) queryProfilesByType(c ctx.Context, odataType MsGraphOdataType) *azruntime.Pager[azcosmos.QueryItemsResponse] {
	partitionKey := azcosmos.NewPartitionKeyString(s.TenantID().String())
	pager := s.AzCosmosContainerClient().NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.displayName FROM c
WHERE c.namespaceId = @namespaceId
  AND c["@odata.type"] = @odataType`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: s.TenantID().String()},
				{Name: "@odataType", Value: odataType},
			},
		})

	return pager
}

func (doc *GraphDoc) IsValid() bool {
	switch doc.OdataType {
	case MsGraphOdataTypeDevice,
		MsGraphOdataTypeUser,
		MsGraphOdataTypeGroup,
		MsGraphOdataTypeApplication,
		MsGraphOdataTypeServicePrincipal:
		return true
	}
	return false
}

func (doc *GraphDoc) GetOdataType() MsGraphOdataType {
	return doc.OdataType
}
