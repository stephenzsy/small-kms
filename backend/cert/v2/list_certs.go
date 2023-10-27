package cert

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListCertificates implements ServerInterface.
func listCertificates(c ctx.RequestContext, params ListCertificatesParams) ([]*CertificateRef, error) {
	docService := base.GetAzCosmosCRUDService(c)
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(certDocQueryColumnThumbprintSHA1, certDocQueryColumnNotAfter).
		WithOrderBy(fmt.Sprintf("%s DESC", certDocQueryColumnCreated))
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCert)

	if params.PolicyId != nil {
		policyIdentifier := base.ParseIdentifier(*params.PolicyId)

		policyLocator := base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, policyIdentifier)

		qb.WhereClauses = append(qb.WhereClauses, "c.policy = @policy")
		qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyLocator.String()})
	}

	pager := base.NewQueryDocPager[*CertListQueryDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertListQueryDoc) *CertificateRef {
		r := &CertificateRef{}
		d.PopulateModelRef(r)
		return r
	})
	return utils.PagerToSlice(c, modelPager)
}
