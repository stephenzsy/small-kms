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
	MsGraphOdataTypeDevice           MsGraphOdataType = "#microsoft.graph.device"
	MsGraphOdataTypeUser             MsGraphOdataType = "#microsoft.graph.user"
	MsGraphOdataTypeGroup            MsGraphOdataType = "#microsoft.graph.group"
	MsGraphOdataTypeApplication      MsGraphOdataType = "#microsoft.graph.application"
	MsGraphOdataTypeServicePrincipal MsGraphOdataType = "#microsoft.graph.servicePrincipal"
)

// this kind of docs represent a graph object browsable in the app, it should not be used to persist data, always query live graph service
type GraphDoc struct {
	kmsdoc.BaseDoc

	DisplayName string `json:"displayName"`

	service *graphService
}

type GraphProfileable interface {
	msgraphmodels.DirectoryObjectable
	GetDisplayName() *string
}

func (s *graphService) init(doc *GraphDoc, graphObj GraphProfileable, extType kmsdoc.KmsDocTypeExtName) {
	oid, _ := uuid.Parse(utils.NilToDefault(graphObj.GetId()))
	doc.NamespaceID = s.TenantID()
	doc.BaseDoc.ID = kmsdoc.NewKmsDocIDExt(kmsdoc.DocTypeMsGraphObject, oid, extType)
	doc.DisplayName = utils.NilToDefault(graphObj.GetDisplayName())
}

func (s *graphService) queryProfilesByType(c ctx.Context, docExtension kmsdoc.KmsDocTypeExtName) *azruntime.Pager[azcosmos.QueryItemsResponse] {
	partitionKey := azcosmos.NewPartitionKeyString(s.TenantID().String())
	pager := s.AzCosmosContainerClient().NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.displayName FROM c
WHERE c.namespaceId = @namespaceId
  AND c.extType = @extType`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: s.TenantID().String()},
				{Name: "@extType", Value: docExtension},
			},
		})

	return pager
}
