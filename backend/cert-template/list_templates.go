package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/auth"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/profile"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getTemplateReservedDefault(id common.Identifier) *models.CertificateTemplateRef {
	return &models.CertificateTemplateRef{
		Id: id,
	}
}

// ListCertificateTemplates implements CertificateTemplateService.
func (*certTmplService) ListCertificateTemplates(c common.ServiceContext) ([]*models.CertificateTemplateRef, error) {
	if err := auth.AuthorizeAdminOnly(c); err != nil {
		return nil, err
	}

	pcs := profile.GetProfileContextService(c)
	nsID := pcs.GetResourceDocNsID()

	nsCap, err := ns.GetNamespaceCapabilities(nsID)
	if err != nil {
		return nil, err
	}

	itemsPager := kmsdoc.QueryItemsPager[*CertificateTemplateDoc](c,
		nsID,
		kmsdoc.DocKindCertificateTemplate,
		func(items []string) []string {
			return append(items, "subjectCn")
		},
		kmsdoc.DefaultQueryGetWhereClause,
		nil)
	mappedPager := utils.NewMappedPager(itemsPager, func(doc *CertificateTemplateDoc) *models.CertificateTemplateRef {
		return doc.toModelRef()
	})
	allItems, err := utils.PagerAllItems[*models.CertificateTemplateRef](mappedPager, c)
	if err != nil {
		return nil, err
	}
	reservedMapping := nsCap.GetReservedCertificateTemplateNames(pcs.GetRequestProfileType())
	if reservedMapping != nil {
		reservedDefaults := make([]*models.CertificateTemplateRef, len(reservedMapping))
		for i, v := range reservedMapping {
			reservedDefaults[v] = getTemplateReservedDefault(i)
		}
		return utils.ReservedFirst(allItems, reservedDefaults, func(item *models.CertificateTemplateRef) int {
			if ind, inMap := reservedMapping[item.Id]; inMap {
				return ind
			}
			return -1
		}), nil
	} else if allItems == nil {
		return make([]*models.CertificateTemplateRef, 0), nil
	} else {
		return allItems, nil
	}
}
