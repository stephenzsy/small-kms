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
	qb := base.NewDefaultCosmoQueryBuilder(base.ResourceKindCert).
		WithExtraColumns(certDocQueryColumnThumbprintSHA1, certDocQueryColumnNotAfter, "c[\"@rels\"].namedFrom[\"issuer-cert\"] AS issuerCertPolicyId").
		WithOrderBy(fmt.Sprintf("%s DESC", certDocQueryColumnCreated))
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.GetDefaultStorageNamespaceID(c, nsCtx.Kind(), nsCtx.Identifier())

	if params.PolicyId != nil {
		policyIdentifier := base.IdentifierFromString(*params.PolicyId)

		policyLocator := base.GetDefaultStorageLocator(c, nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, policyIdentifier)

		qb.ExtraWhereClauses = append(qb.ExtraWhereClauses, "c.policyLocator = @policyLocator")
		qb.ExtraParameters = append(qb.ExtraParameters, azcosmos.QueryParameter{Name: "@policyLocator", Value: policyLocator.String()})
	}

	pager := base.NewQueryDocPager[*CertListQueryDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertListQueryDoc) *CertificateRef {
		r := &CertificateRef{}
		d.PopulateModelRef(r)
		r.Id.NID = storageNsID
		r.NamespaceKind = nsCtx.Kind()
		r.NamespaceIdentifier = nsCtx.Identifier()
		r.ResourceKind = base.ResourceKindCert
		return r
	})
	return utils.PagerToSlice(c, modelPager)
}
