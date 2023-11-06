package cert

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/api"
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
func apiListCertificates(c ctx.RequestContext, params ListCertificatesParams) error {
	qb := base.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(certDocQueryColumnThumbprintSHA1, certDocQueryColumnNotAfter).
		WithOrderBy(fmt.Sprintf("%s DESC", certDocQueryColumnCreated))
	nsCtx := ns.GetNSContext(c)
	storageNsID := base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCert)

	if params.PolicyId != nil {
		policyIdentifier := base.ParseID(*params.PolicyId)

		policyLocator := base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy, policyIdentifier)

		qb.WhereClauses = append(qb.WhereClauses, "c.policy = @policy")
		qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyLocator.String()})
	}

	pager := base.NewQueryDocPager[*CertQueryDoc](c, qb, storageNsID)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertQueryDoc) *CertificateRef {
		r := &CertificateRef{}
		d.PopulateModelRef(r)
		return r
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}

func QueryLatestCertificateIdsIssuedByPolicy(c ctx.RequestContext, policyFullIdentifier base.DocFullIdentifier, limit uint) ([]ID, error) {
	qb := base.NewDefaultCosmoQueryBuilder().
		WithOrderBy(fmt.Sprintf("%s DESC", certDocQueryColumnCreated)).
		WithOffsetLimit(0, limit)
	qb.WhereClauses = append(qb.WhereClauses, "c.policy = @policy", "NOT IS_DEFINED(c.deleted)", "c.status = 'issued'")
	qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyFullIdentifier.String()})
	pager := base.NewQueryDocPager[*CertQueryDoc](c,
		qb,
		base.NewDocNamespacePartitionKey(policyFullIdentifier.NamespaceKind(), policyFullIdentifier.NamespaceID(), base.ResourceKindCert))

	return utils.PagerToSlice(utils.NewMappedItemsPager(pager, func(d *CertQueryDoc) ID {
		return d.ID
	}))
}
