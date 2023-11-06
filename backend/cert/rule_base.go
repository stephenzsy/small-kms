package cert

import "github.com/stephenzsy/small-kms/backend/base"

func getNamespaceCertificateRuleDocFullIdentifier(
	nsKind base.NamespaceKind, nsID base.ID, ruleName base.NamespaceConfigName) base.DocFullIdentifier {
	return base.NewDocFullIdentifier(nsKind, nsID, base.ResourceKindNamespaceConfig, base.ID(ruleName))
}

type CertRulePolicyDoc struct {
	base.BaseDoc
	PolicyID base.ID `json:"policyId"`
}

func (d *CertRulePolicyDoc) init(
	nsKind base.NamespaceKind, nsIdentifier base.ID, ruleName base.NamespaceConfigName,
) {
	d.BaseDoc.Init(nsKind, nsIdentifier, base.ResourceKindNamespaceConfig, base.ID(ruleName))
}

type CertRuleIssuerLatestCertificateDoc struct {
	CertRulePolicyDoc
	CertificateID base.ID `json:"certificateId"`
}

type CertRuleIssuerLastNCertificateDoc struct {
	CertRulePolicyDoc
	CertificateIDs []base.ID `json:"certificateIds"`
}
