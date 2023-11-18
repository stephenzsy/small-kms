package cert

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/admin"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertQueryDoc struct {
	resdoc.ResourceQueryDoc
	Status         certmodels.CertificateStatus        `json:"status"`
	ThumbprintSHA1 cloudkey.Base64RawURLEncodableBytes `json:"x5t"`
	NotAfter       jwt.NumericDate                     `json:"exp"`
	IssuedAt       jwt.NumericDate                     `json:"iat"`
	Policy         resdoc.DocIdentifier                `json:"policy"`
}

func (d *CertQueryDoc) ToRef() (m certmodels.CertificateRef) {
	m.Ref = d.ResourceQueryDoc.ToRef()
	m.Thumbprint = d.ThumbprintSHA1.HexString()
	m.Exp = d.NotAfter
	if !d.IssuedAt.Time.IsZero() {
		m.Iat = &d.IssuedAt
	}
	m.PolicyIdentifier = d.Policy.String()
	m.Status = d.Status
	return m
}

// ListCertificates implements ServerInterface.
func (*CertServer) ListCertificates(ec echo.Context, namespaceProvider models.NamespaceProvider, namespaceId string, params admin.ListCertificatesParams) error {
	c := ec.(ctx.RequestContext)
	namespaceId = ns.ResolveMeNamespace(c, namespaceId)
	if _, authOk := authz.Authorize(c, authz.AllowAdmin, authz.AllowSelf(namespaceId)); !authOk {
		return base.ErrResponseStatusForbidden
	}

	qb := resdoc.NewDefaultCosmoQueryBuilder().
		WithExtraColumns(certDocQueryColStatus, certDocQueryColIssuedAt, certDocQueryColNotAfter, certDocQueryColThumbprintSHA1).
		WithOrderBy("c.iat DESC")
	if params.PolicyId != nil && *params.PolicyId != "" {
		policyIdentifer := resdoc.NewDocIdentifier(
			namespaceProvider, namespaceId,
			models.ResourceProviderCertPolicy,
			*params.PolicyId)
		qb.WithWhereClauses("c.policy = @policy")
		qb.Parameters = append(qb.Parameters, azcosmos.QueryParameter{Name: "@policy", Value: policyIdentifer.String()})
	} else {
		qb.WithExtraColumns(certDocQueryColPolicy)
	}
	pager := resdoc.NewQueryDocPager[*CertQueryDoc](c, qb, resdoc.PartitionKey{
		NamespaceProvider: namespaceProvider,
		NamespaceID:       namespaceId,
		ResourceProvider:  models.ResourceProviderCert,
	})

	modelPager := utils.NewMappedItemsPager(pager, func(doc *CertQueryDoc) certmodels.CertificateRef {
		return doc.ToRef()
	})
	return api.RespondPagerList(c, utils.NewSerializableItemsPager(modelPager))
}
