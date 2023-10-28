package base

type NamespaceConfigName string

const (
	CertRuleNameIssuer                  NamespaceConfigName = "issuer"
	CertRuleNameMsEntraClientCredential NamespaceConfigName = "ms-entra-client-credential"

	AgentConfigNameServer NamespaceConfigName = "agent-server"
)
