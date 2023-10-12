package cert

import (
	"errors"
	"slices"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getLatestCertificateByTemplateDoc(c RequestContext, templateLocator shared.ResourceLocator) (doc *CertDoc, err error) {
	doc = &CertDoc{}
	err = kmsdoc.Read[*CertDoc](c,
		shared.NewResourceLocator(templateLocator.GetNamespaceID(), shared.NewResourceIdentifier(shared.ResourceKindLatestCertForTemplate, templateLocator.GetID().Identifier())), doc)
	return
}

func ListCertificatesByTemplate(c RequestContext) ([]*shared.CertificateRef, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	ctc := ct.GetCertificateTemplateContext(c)
	tmplLocator := ctc.GetCertificateTemplateLocator(c)

	itemsPager := kmsdoc.QueryItemsPager[*CertDoc](c,
		nsID,
		shared.ResourceKindCert,
		kmsdoc.CosmosQueryBuilder{
			ExtraColumns: []string{"c.thumbprint", queryColumnStatus},
			ExtraWhereClauses: []string{
				queryColumnTemplate + " = @templateId",
			},
			OrderBy: "c.notBefore DESC",
			ExtraParameters: []azcosmos.QueryParameter{
				{Name: "@templateId", Value: tmplLocator.String()},
			},
		})
	mappedPager := utils.NewMappedItemsPager(itemsPager, func(doc *CertDoc) *shared.CertificateRef {
		return doc.toModelRef()
	})
	allItems, err := utils.PagerAllItems[*shared.CertificateRef](mappedPager, c)
	if err != nil {
		return nil, err
	}
	if latestDoc, err := getLatestCertificateByTemplateDoc(c, tmplLocator); err != nil {
		if !errors.Is(err, common.ErrStatusNotFound) {
			return nil, err
		}
	} else {
		cmpId := latestDoc.AliasTo.GetID().Identifier()
		matchedInd := slices.IndexFunc(allItems, func(item *shared.CertificateRef) bool {
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
		return make([]*shared.CertificateRef, 0), nil
	}
	return allItems, nil
}
