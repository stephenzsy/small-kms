package cert

import "github.com/stephenzsy/small-kms/backend/base"

func getNamespaceCertificateRuleDocFullIdentifier(
	nsKind base.NamespaceKind, nsID base.Identifier, ruleName base.NamespaceConfigName) base.DocFullIdentifier {
	return base.NewDocFullIdentifier(nsKind, nsID, base.ResourceKindNamespaceConfig, base.StringIdentifier(string(ruleName)))
}

type CertRulePolicyDoc struct {
	base.BaseDoc
	PolicyID base.Identifier `json:"policyId"`
}

func (d *CertRulePolicyDoc) init(
	nsKind base.NamespaceKind, nsIdentifier base.Identifier, ruleName base.NamespaceConfigName,
) {
	d.BaseDoc.Init(nsKind, nsIdentifier, base.ResourceKindNamespaceConfig, base.StringIdentifier(string(ruleName)))
}

type CertRuleIssuerLatestCertificateDoc struct {
	CertRulePolicyDoc
	CertificateID base.Identifier `json:"certificateId"`
}

type CertRuleIssuerLastNCertificateDoc struct {
	CertRulePolicyDoc
	CertificateIDs []base.Identifier `json:"certificateIds"`
}
