package cert

import (
	"context"
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type CertRuleIssuerDoc = CertRuleIssuerLatestCertificateDoc

// PopulateModel implements base.ModelPopulater.
func (d *CertRuleIssuerDoc) PopulateModel(r *CertificateRuleIssuer) {
	if d == nil || r == nil {
		return
	}
	r.CertificateId = d.CertificateID
	r.PolicyId = d.PolicyID
}

var _ base.ModelPopulater[CertificateRuleIssuer] = (*CertRuleIssuerDoc)(nil)

func readCertRuleIssuerDoc(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertRuleIssuerDoc, error) {
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleIssuerDoc)
	err := docSvc.Read(c, getNamespaceCertificateRuleDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.ID(), base.CertRuleNameIssuer), ruleDoc, nil)
	return ruleDoc, err
}

func apiGetCertRuleIssuer(c ctx.RequestContext) error {
	nsCtx := ns.GetNSContext(c)
	ruleDoc, err := readCertRuleIssuerDoc(c, base.NewNamespaceIdentifier(nsCtx.Kind(), nsCtx.ID()))
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w, issuer configuration not found: %s", base.ErrResponseStatusNotFound, base.CertRuleNameIssuer)
		}
		return err
	}
	m := new(CertificateRuleIssuer)
	ruleDoc.PopulateModel(m)
	return c.JSON(200, m)
}
