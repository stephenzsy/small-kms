package models

type NamespaceProvider string

const (
	NamespaceProviderProfile NamespaceProvider = "profile"
	NamespaceProviderAgent   NamespaceProvider = "agent"
)

type ResourceProvider string

const (
	ProfileResourceProviderSystem ResourceProvider = "sys"
	ProfileResourceProviderAgent  ResourceProvider = "agent"
)
