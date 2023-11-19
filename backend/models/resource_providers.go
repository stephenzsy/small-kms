package models

type ResourceProvider string

const (
	ProfileResourceProviderSystem           ResourceProvider = "sys"
	ProfileResourceProviderAgent            ResourceProvider = "agent"
	ProfileResourceProviderServicePrincipal ResourceProvider = "service-principal"
	ProfileResourceProviderUser             ResourceProvider = "user"
	ProfileResourceProviderGroup            ResourceProvider = "group"
	ResourceProviderAgentConfig             ResourceProvider = "agent-config"
	ResourceProviderKeyPolicy               ResourceProvider = "key-policy"
	ResourceProviderCert                    ResourceProvider = "cert"
	ResourceProviderCertPolicy              ResourceProvider = "cert-policy"
	ResourceProviderLink                    ResourceProvider = "link"
)
