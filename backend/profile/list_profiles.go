package profile

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
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
	case models.NamespaceKindCaRoot:
		return utils.MapSlices(getBuiltInRootCaProfiles(), func(doc ProfileDoc) *models.ProfileRefComposed { return doc.toModelRef() }), nil
	case models.NamespaceKindCaInt:
		return utils.MapSlices(getBuiltInIntermediateCaProfiles(), func(doc ProfileDoc) *models.ProfileRefComposed { return doc.toModel() }), nil
	}
	itemsPager := kmsdoc.QueryItemsPager[*ProfileDoc](c,
		docNsIDProfileTenant,
		models.ResourceKindMsGraph,
		func(tbl string) kmsdoc.CosmosQueryBuilder {
			return kmsdoc.CosmosQueryBuilder{
				ExtraColumns:      []string{"displayName"},
				ExtraWhereClauses: []string{tbl + ".profileType = @profileType"},
				ExtraParameters: []azcosmos.QueryParameter{
					{Name: "@profileType", Value: profileType},
				}}
		})
	allItems, err := utils.PagerAllItems[*models.ProfileRefComposed](utils.NewMappedPager(itemsPager, func(doc *ProfileDoc) *models.ProfileRefComposed {
		return doc.toModel()
	}), c)
	if allItems == nil {
		allItems = make([]*models.ProfileRefComposed, 0)
	}
	return allItems, err
}
