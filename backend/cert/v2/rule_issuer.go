package cert

import (
	"context"
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type CertRuleIssuerDoc struct {
	base.BaseDoc
	PolicyID      base.Identifier `json:"policyId"`
	CertificateID base.Identifier `json:"certificateId"`
}

func (d *CertRuleIssuerDoc) init(
	nsKind base.NamespaceKind, nsIdentifier base.Identifier,
) {
	d.BaseDoc.Init(nsKind, nsIdentifier, base.ResourceKindNamespaceConfig, base.StringIdentifier(CertRuleNameIssuer))
}

// PopulateModel implements base.ModelPopulater.
func (d *CertRuleIssuerDoc) PopulateModel(r *CertificateRuleIssuer) {
	if d == nil || r == nil {
		return
	}
	r.CertificateId = &d.CertificateID
	r.PolicyId = d.PolicyID
}

var _ base.ModelPopulater[CertificateRuleIssuer] = (*CertRuleIssuerDoc)(nil)

func readCertRuleIssuerDoc(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertRuleIssuerDoc, error) {
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleIssuerDoc)
	err := docSvc.Read(c, getNamespaceCertificateRuleDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.Identifier(), CertRuleNameIssuer), ruleDoc, nil)
	return ruleDoc, err
}

func getNamespaceIssuerCert(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertDoc, error) {

	ruleDoc, err := readCertRuleIssuerDoc(c, nsIdentifier)
	if err != nil {
		return nil, err
	}

	docSvc := base.GetAzCosmosCRUDService(c)
	certDoc := new(CertDoc)
	err = docSvc.Read(c, base.NewDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.Identifier(), base.ResourceKindCert, ruleDoc.CertificateID), certDoc, nil)
	return certDoc, err
}

func apiGetCertRuleIssuer(c ctx.RequestContext) error {
	nsCtx := ns.GetNSContext(c)
	ruleDoc, err := readCertRuleIssuerDoc(c, base.NewNamespaceIdentifier(nsCtx.Kind(), nsCtx.Identifier()))
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w, issuer configuration not found: %s", base.ErrResponseStatusNotFound, CertRuleNameIssuer)
		}
		return err
	}
	m := new(CertificateRuleIssuer)
	ruleDoc.PopulateModel(m)
	return c.JSON(200, m)
}

func apiPutCertRuleIssuer(c ctx.RequestContext, p *CertificateRuleIssuer) error {
	nsCtx := ns.GetNSContext(c)
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleIssuerDoc)
	ruleDoc.init(nsCtx.Kind(), nsCtx.Identifier())
	ruleDoc.PolicyID = p.PolicyId
	if p.CertificateId.IsNilOrEmpty() {
		certId, err := queryLatestCertificateIdIssuedByPolicy(c, base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, p.PolicyId))
		if err != nil {
			return err
		}
		if certId != nil {
			ruleDoc.CertificateID = *certId
		} else {
			return fmt.Errorf("%w, no certificate issued by policy: %s", base.ErrResponseStatusNotFound, p.PolicyId.String())
		}
	} else {
		ruleDoc.CertificateID = *p.CertificateId
	}
	docSvc.Upsert(c, ruleDoc, nil)
	ruleDoc, err := readCertRuleIssuerDoc(c, base.NewNamespaceIdentifier(nsCtx.Kind(), nsCtx.Identifier()))
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w, issuer configuration not found: %s", base.ErrResponseStatusNotFound, CertRuleNameIssuer)
		}
		return err
	}
	m := new(CertificateRuleIssuer)
	ruleDoc.PopulateModel(m)
	return c.JSON(200, m)
}
