package cert

import "github.com/stephenzsy/small-kms/backend/base"

type CertRuleName string

const (
	CertRuleNameIssuer                  = "issuer"
	CertRuleNameMsEntraClientCredential = "msEntraClientCredential"
)

func getNamespaceCertificateRuleDocFullIdentifier(
	nsKind base.NamespaceKind, nsID base.Identifier, name CertRuleName) base.DocFullIdentifier {
	return base.NewDocFullIdentifier(nsKind, nsID, base.ResourceKindNamespaceConfig, base.StringIdentifier(CertRuleNameIssuer))
}
