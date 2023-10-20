package profile

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getBuiltInRootCaProfiles() []ProfileDoc {
	return []ProfileDoc{
		rootCaProfileDocs[idCaRoot],
		rootCaProfileDocs[idCaRootTest],
	}
}

func getBuiltInIntermediateCaProfiles() []ProfileDoc {
	return []ProfileDoc{
		intCaProfileDocs[idIntCaServices],
		intCaProfileDocs[idIntCaIntranet],
		intCaProfileDocs[idIntCaMsEntraClientSecret],
		intCaProfileDocs[idIntCaTest],
	}
}

// ListProfiles implements ProfileService.
func ListProfiles(c RequestContext, profileType models.NamespaceKind) ([]*models.ProfileRefComposed, error) {
	switch profileType {
	case shared.NamespaceKindCaRoot:
		return utils.MapSlice(getBuiltInRootCaProfiles(), func(doc ProfileDoc) *models.ProfileRefComposed { return doc.toModelRef() }), nil
	case shared.NamespaceKindCaInt:
		return utils.MapSlice(getBuiltInIntermediateCaProfiles(), func(doc ProfileDoc) *models.ProfileRefComposed { return doc.toModel() }), nil
	}
	itemsPager := kmsdoc.QueryItemsPager[*ProfileDoc](c,
		docNsIDProfileTenant,
		shared.ResourceKindMsGraph,
		kmsdoc.CosmosQueryBuilder{
			ExtraColumns:      []string{"c.displayName", "c.isAppManaged"},
			ExtraWhereClauses: []string{"c.profileType = @profileType"},
			ExtraParameters: []azcosmos.QueryParameter{
				{Name: "@profileType", Value: profileType},
			},
		})
	allItems, err := utils.PagerAllItems[*models.ProfileRefComposed](
		utils.NewMappedItemsPager(itemsPager, func(doc *ProfileDoc) *models.ProfileRefComposed {
			return doc.toModel()
		}), c)
	if allItems == nil {
		allItems = make([]*models.ProfileRefComposed, 0)
	}
	return allItems, err
}
