package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getTemplateReservedDefault(nsID shared.NamespaceIdentifier, id shared.Identifier) *models.CertificateTemplateRefComposed {
	return &models.CertificateTemplateRefComposed{
		ResourceRef: shared.ResourceRef{
			Id:      id,
			Locator: shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindCertTemplate, id)),
		},
	}
}

// ListCertificateTemplates implements CertificateTemplateService.
func ListCertificateTemplates(c RequestContext) ([]*models.CertificateTemplateRefComposed, error) {

	nsc := ns.GetNamespaceContext(c)
	nsID := nsc.GetID()

	itemsPager := kmsdoc.QueryItemsPager[*CertificateTemplateDoc](c,
		nsID,
		shared.ResourceKindCertTemplate,
		kmsdoc.CosmosQueryBuilder{
			ExtraColumns: []string{"c.subjectCn", kmsdoc.QueryColumnNameOwner},
		})
	mappedPager := utils.NewMappedItemsPager(itemsPager, func(doc *CertificateTemplateDoc) *models.CertificateTemplateRefComposed {
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
