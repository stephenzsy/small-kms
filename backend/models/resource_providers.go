package models

type ResourceProvider string

const (
	ProfileResourceProviderSystem           ResourceProvider = "sys"
	ProfileResourceProviderRootCA           ResourceProvider = ResourceProvider(NamespaceProviderRootCA)
	ProfileResourceProviderIntermediateCA   ResourceProvider = ResourceProvider(NamespaceProviderIntermediateCA)
	ProfileResourceProviderAgent            ResourceProvider = "agent"
	ProfileResourceProviderServicePrincipal ResourceProvider = "service-principal"
	ProfileResourceProviderUser             ResourceProvider = "user"
	ProfileResourceProviderGroup            ResourceProvider = "group"
	ResourceProviderAgentConfig             ResourceProvider = "agent-config"
	ResourceProviderAgentInstance           ResourceProvider = "agent-instance"
	ResourceProviderKey                     ResourceProvider = "key"
	ResourceProviderKeyPolicy               ResourceProvider = "key-policy"
	ResourceProviderCert                    ResourceProvider = "cert"
	ResourceProviderCertPolicy              ResourceProvider = "cert-policy"
	ResourceProviderLink                    ResourceProvider = "link"
)

type LinkProvider string

const (
	LinkProviderCAPolicyIssuerCertificate LinkProvider = "issuer-cert"
	LinkProviderGraphMemberOf             LinkProvider = "graph-member-of"
	LinkProviderGraphMember               LinkProvider = "graph-member"
)
