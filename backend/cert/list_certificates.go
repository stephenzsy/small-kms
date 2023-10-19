package cert

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func ListActiveCertDocsByTemplateID(c context.Context, templateId shared.Identifier) ([]*CertDoc, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	templateLocator := shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindCertTemplate, templateId))
	itemsPager := kmsdoc.QueryItemsPager[*CertDoc](c,
		nsID,
		shared.ResourceKindCert,
		kmsdoc.CosmosQueryBuilder{
			ExtraColumns: []string{"c.thumbprint", queryColumnStatus},
			ExtraWhereClauses: []string{
				queryColumnTemplate + " = @templateId",
				queryColumnStatus + " = @status",
				"IS_NULL(c.deleted)",
				queryColumnNotAfter + " > GetCurrentDateTime()",
			},
			OrderBy: "c.notBefore DESC",
			ExtraParameters: []azcosmos.QueryParameter{
				{Name: "@templateId", Value: templateLocator.String()},
				{Name: "@status", Value: CertStatusIssued},
			},
		})
	return utils.PagerAllItems[*CertDoc](itemsPager, c)
}

func ApiListCertificatesByTemplate(c RequestContext) error {
	result, err := listCertificatesByTemplate(c)
	if err != nil {
		return err
	}
	if result == nil {
		result = make([]*shared.CertificateRef, 0)
	}
	return c.JSON(200, result)
}

func listCertificatesByTemplate(c RequestContext) ([]*shared.CertificateRef, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	ctc := ct.GetCertificateTemplateContext(c)
	tmplLocator := ctc.GetCertificateTemplateLocator(c)

	itemsPager := kmsdoc.QueryItemsPager[*CertDoc](c,
		nsID,
		shared.ResourceKindCert,
		kmsdoc.CosmosQueryBuilder{
			ExtraColumns: []string{queryColumnThumbprint, queryColumnNotAfter, queryColumnStatus},
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
	return utils.PagerAllItems[*shared.CertificateRef](mappedPager, c)
}
