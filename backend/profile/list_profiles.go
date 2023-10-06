package profile

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
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
func (*profileService) ListProfiles(c common.ServiceContext, profileType models.ProfileType) ([]*models.ProfileRef, error) {
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	switch profileType {
	case models.ProfileTypeRootCA:
		return utils.MapSlices(getBuiltInRootCaProfiles(), func(doc ProfileDoc) *models.ProfileRef { return doc.toModel() }), nil
	case models.ProfileTypeIntermediateCA:
		return utils.MapSlices(getBuiltInIntermediateCaProfiles(), func(doc ProfileDoc) *models.ProfileRef { return doc.toModel() }), nil
	}
	itemsPager := kmsdoc.QueryItemsPager[*ProfileDoc](c, docNsIDProfileTenant, func(items []string) []string {
		return append(items, "displayName")
	}, func(tbl string) string {
		sb := strings.Builder{}
		sb.WriteString(tbl)
		sb.WriteString(".profileType = @profileType")
		return sb.String()
	}, []azcosmos.QueryParameter{
		{Name: "@profileType", Value: profileType},
	})
	allItems, err := utils.PagerAllItems[*models.ProfileRef](utils.NewMappedPager(itemsPager, func(doc *ProfileDoc) *models.ProfileRef {
		return doc.toModel()
	}), c)
	if allItems == nil {
		allItems = make([]*models.ProfileRef, 0)
	}
	return allItems, err
}
