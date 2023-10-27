package cert

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/base"
)

type CertRuleIssuerDoc struct {
	base.BaseDoc
	PolicyID      base.Identifier `json:"policyId"`
	CertificateID base.Identifier `json:"certificateId"`
}

func getNamespaceIssuerCert(c context.Context, nsIdentifier base.NamespaceIdentifier) (*CertDoc, error) {
	docSvc := base.GetAzCosmosCRUDService(c)

	ruleDoc := new(CertRuleIssuerDoc)
	err := docSvc.Read(c, getNamespaceCertificateRuleDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.Identifier(), CertRuleNameIssuer), ruleDoc, nil)
	if err != nil {
		return nil, err
	}

	certDoc := new(CertDoc)
	err = docSvc.Read(c, base.NewDocFullIdentifier(nsIdentifier.Kind(), nsIdentifier.Identifier(), base.ResourceKindCert, ruleDoc.CertificateID), certDoc, nil)

	return certDoc, err
}
