package cert

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertQueryDoc struct {
	base.QueryBaseDoc
	ThumbprintSHA1 base.Base64RawURLEncodedBytes `json:"x5t"`
	NotAfter       base.NumericDate              `json:"exp"`
}

// PopulateModelRef implements base.ModelRefPopulater.
func (d *CertQueryDoc) PopulateModelRef(m *CertificateRef) {
	if d == nil || m == nil {
		return
	}
	d.QueryBaseDoc.PopulateModelRef(&m.ResourceReference)
	m.Thumbprint = d.ThumbprintSHA1.HexString()
	m.Attributes.Exp = &d.NotAfter
}

var _ base.ModelRefPopulater[CertificateRef] = (*CertQueryDoc)(nil)

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

	pager := base.NewQueryDocPager[*CertQueryDoc](docService, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertQueryDoc) *CertificateRef {
		r := &CertificateRef{}
		d.PopulateModelRef(r)
		return r
	})
	return utils.PagerToSlice(c, modelPager)
}

func queryLatestCertificateIdsIssuedByPolicy(c ctx.RequestContext, policyFullIdentifier base.DocFullIdentifier, limit uint) ([]base.Identifier, error) {
	qb := base.NewDefaultCosmoQueryBuilder().
		WithOrderBy(fmt.Sprintf("%s DESC", certDocQueryColumnCreated)).
		WithOffsetLimit(0, limit)
	qb.WhereClauses = append(qb.WhereClauses, "c.policy = @policy", "NOT IS_DEFINED(c.deleted)", "c.status = 'issued'")
	qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyFullIdentifier.String()})
	docService := base.GetAzCosmosCRUDService(c)
	pager := base.NewQueryDocPager[*CertQueryDoc](docService, qb, base.NewDocNamespacePartitionKey(policyFullIdentifier.NamespaceKind(), policyFullIdentifier.NamespaceIdentifier(), base.ResourceKindCert))

	return utils.PagerAllItems(utils.NewMappedItemsPager(pager, func(d *CertQueryDoc) Identifier {
		return d.ID
	}), c)
}
