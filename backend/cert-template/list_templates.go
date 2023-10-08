package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getTemplateReservedDefault(nsID models.NamespaceID, id common.Identifier) *models.CertificateTemplateRefComposed {
	return &models.CertificateTemplateRefComposed{
		ResourceRef: models.ResourceRef{
			Id:      id,
			Locator: common.NewLocator(nsID, common.NewIdentifierWithKind(models.ResourceKindCertTemplate, id)),
		},
	}
}

// ListCertificateTemplates implements CertificateTemplateService.
func ListCertificateTemplates(c common.ServiceContext) ([]*models.CertificateTemplateRefComposed, error) {

	nsc := ns.GetNamespaceContext(c)
	nsID := nsc.GetID()

	itemsPager := kmsdoc.QueryItemsPager[*CertificateTemplateDoc](c,
		nsID,
		models.ResourceKindCertTemplate,
		func(tbl string) kmsdoc.CosmosQueryBuilder {
			return kmsdoc.CosmosQueryBuilder{
				ExtraColumns: []string{"subjectCn"},
			}
		})
	mappedPager := utils.NewMappedPager(itemsPager, func(doc *CertificateTemplateDoc) *models.CertificateTemplateRefComposed {
		return doc.toModelRef()
	})
	allItems, err := utils.PagerAllItems[*models.CertificateTemplateRefComposed](mappedPager, c)
	if err != nil {
		return nil, err
	}
	reservedMapping := ns.GetReservedCertificateTemplateNames(nsID)
	if reservedMapping != nil {
		reservedDefaults := make([]*models.CertificateTemplateRefComposed, len(reservedMapping))
		for i, v := range reservedMapping {
			reservedDefaults[v] = getTemplateReservedDefault(nsID, i)
		}
		return utils.ReservedFirst(allItems, reservedDefaults, func(item *models.CertificateTemplateRefComposed) int {
			if ind, inMap := reservedMapping[item.Id]; inMap {
				return ind
			}
			return -1
		}), nil
	} else if allItems == nil {
		return make([]*models.CertificateTemplateRefComposed, 0), nil
	} else {
		return allItems, nil
	}
}
