package models

type ResourceProvider string

const (
	ProfileResourceProviderSystem ResourceProvider = "sys"
	ProfileResourceProviderAgent  ResourceProvider = "agent"
	ResourceProviderAgentConfig   ResourceProvider = "agent-config"
	ResourceProviderKeyPolicy     ResourceProvider = "key-policy"
	ResourceProviderCertPolicy    ResourceProvider = "cert-policy"
)
