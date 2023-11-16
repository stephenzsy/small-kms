package models

type ResourceProvider string

const (
	ProfileResourceProviderSystem ResourceProvider = "sys"
	ProfileResourceProviderAgent  ResourceProvider = "agent"
	ResourceProviderKeyPolicy     ResourceProvider = "key-policy"
)
