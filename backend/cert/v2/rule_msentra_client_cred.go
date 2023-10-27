package cert

import (
	"context"
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type CertRuleMsEntraClientCredDoc = CertRuleIssuerLastNCertificateDoc

// PopulateModel implements base.ModelPopulater.
func (d *CertRuleMsEntraClientCredDoc) PopulateModel(r *CertificateRuleMsEntraClientCredential) {
	if d == nil || r == nil {
		return
	}
	r.CertificateIds = d.CertificateIDs
	r.PolicyId = d.PolicyID
}

// var _ base.ModelPopulater[CertificateRuleIssuer] = (*CertRuleIssuerDoc)(nil)

func readCertRuleMsEntraClientDoc(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertRuleIssuerLastNCertificateDoc, error) {
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleIssuerLastNCertificateDoc)
	err := docSvc.Read(c, getNamespaceCertificateRuleDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.Identifier(), CertRuleNameMsEntraClientCredential), ruleDoc, nil)
	return ruleDoc, err
}

// func getNamespaceIssuerCert(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertDoc, error) {

// 	ruleDoc, err := readCertRuleIssuerDoc(c, nsIdentifier)
// 	if err != nil {
// 		return nil, err
// 	}

// 	docSvc := base.GetAzCosmosCRUDService(c)
// 	certDoc := new(CertDoc)
// 	err = docSvc.Read(c, base.NewDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.Identifier(), base.ResourceKindCert, ruleDoc.CertificateID), certDoc, nil)
// 	return certDoc, err
// }

func apiGetCertRuleMsEntraClientCredential(c ctx.RequestContext) error {
	nsCtx := ns.GetNSContext(c)
	ruleDoc, err := readCertRuleMsEntraClientDoc(c, base.NewNamespaceIdentifier(nsCtx.Kind(), nsCtx.Identifier()))
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w, ms entra credential configuration not found: %s", base.ErrResponseStatusNotFound, CertRuleNameMsEntraClientCredential)
		}
		return err
	}
	m := new(CertificateRuleMsEntraClientCredential)
	ruleDoc.PopulateModel(m)
	return c.JSON(200, m)
}

func apiPutCertRuleMsEntraClientCredentrial(c ctx.RequestContext, p *CertificateRuleMsEntraClientCredential) error {
	nsCtx := ns.GetNSContext(c)
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleMsEntraClientCredDoc)
	ruleDoc.init(nsCtx.Kind(), nsCtx.Identifier(), CertRuleNameMsEntraClientCredential)
	ruleDoc.PolicyID = p.PolicyId
	if len(p.CertificateIds) == 0 {
		certIds, err := queryLatestCertificateIdsIssuedByPolicy(c, base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, p.PolicyId), 2)
		if err != nil {
			return err
		}
		ruleDoc.CertificateIDs = certIds
	}
	docSvc.Upsert(c, ruleDoc, nil)
	m := new(CertificateRuleMsEntraClientCredential)
	ruleDoc.PopulateModel(m)
	return c.JSON(200, m)
}
