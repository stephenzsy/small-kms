package cert

import (
	"errors"
	"slices"

	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getLatestCertificateByTemplateDoc(c RequestContext, templateLocator models.ResourceLocator) (doc *CertDoc, err error) {
	doc = &CertDoc{}
	err = kmsdoc.Read[*CertDoc](c,
		common.NewLocator(templateLocator.GetNamespaceID(), common.NewIdentifierWithKind(models.ResourceKindLatestCertForTemplate, templateLocator.GetID().Identifier())), doc)
	return
}

func ListCertificatesByTemplate(c RequestContext) ([]*models.CertificateRefComposed, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	ctc := ct.GetCertificateTemplateContext(c)
	tmplLocator := ctc.GetCertificateTemplateLocator(c)

	itemsPager := kmsdoc.QueryItemsPager[*CertDoc](c,
		nsID,
		models.ResourceKindCert,
		func(tbl string) kmsdoc.CosmosQueryBuilder {
			return kmsdoc.CosmosQueryBuilder{
				ExtraColumns: []string{"thumbprint"},
				OrderBy:      tbl + ".notBefore DESC",
			}
		})
	mappedPager := utils.NewMappedItemsPager(itemsPager, func(doc *CertDoc) *models.CertificateRefComposed {
		return doc.toModelRef()
	})
	allItems, err := utils.PagerAllItems[*models.CertificateRefComposed](mappedPager, c)
	if err != nil {
		return nil, err
	}
	if latestDoc, err := getLatestCertificateByTemplateDoc(c, tmplLocator); err != nil {
		if !errors.Is(err, common.ErrStatusNotFound) {
			return nil, err
		}
	} else {
		cmpId := latestDoc.AliasTo.GetID().Identifier()
		matchedInd := slices.IndexFunc(allItems, func(item *models.CertificateRefComposed) bool {
			return item.Id == cmpId
		})
		if matchedInd >= 0 {
			// shift to first
			matched := allItems[matchedInd]
			for i := matchedInd; i > 0; i-- {
				allItems[i] = allItems[i-1]
			}
			allItems[0] = matched
			matched.Metadata = map[string]any{"latest": true}
		}
	}
	if allItems == nil {
		return make([]*models.CertificateRefComposed, 0), nil
	}
	return allItems, nil
}
