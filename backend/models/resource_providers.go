package models

type NamespaceProvider string

const (
	NamespaceProviderProfile NamespaceProvider = "profile"
	NamespaceProviderAgent   NamespaceProvider = "agent"
)

type ResourceProvider string

const ()

type ProfileResourceProvider ResourceProvider

const (
	ProfileResourceProviderAgent ProfileResourceProvider = "agent"
)
