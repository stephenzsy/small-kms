package cert

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListCertificates implements ServerInterface.
func (s *server) ListCertificates(ec echo.Context,
	namespaceKind base.NamespaceKind,
	namespaceId base.ID, params ListCertificatesParams) error {
	c := ec.(ctx.RequestContext)

	c, nsCtx := ns.WithResovingMeNSContext(c, namespaceKind, namespaceId)

	c, authOk := authz.Authorize(c, authz.AllowAdmin, nsCtx.AllowSelf())
	if !authOk {
		return base.ErrResponseStatusForbidden
	}

	p := api.QueryPolicyItemsParams{
		ExtraColumns: []string{certDocQueryColumnThumbprintSHA1, certDocQueryColumnNotAfter},
	}
	if params.PolicyId != nil {
		p.PolicyLocator = utils.ToPtr(base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy, base.ParseID(*params.PolicyId)))
	}
	pager := api.QueryPolicyItems[*CertQueryDoc](
		c,
		base.NewDocNamespacePartitionKey(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCert),
		p)

	modelPager := utils.NewMappedItemsPager(pager, func(d *CertQueryDoc) *CertificateRef {
		r := &CertificateRef{}
		d.PopulateModelRef(r)
		return r
	})

	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}

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

func QueryLatestCertificateIdsIssuedByPolicy(c ctx.RequestContext, policyFullIdentifier base.DocLocator, limit uint) ([]ID, error) {
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
