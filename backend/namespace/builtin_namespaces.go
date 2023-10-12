package ns

type SystemServiceName string

const (
	SystemServiceNameAgentPush SystemServiceName = "agent-push"
)

type RootCAName string

const (
	RootCANameDefault RootCAName = "default"
	RootCANameTest    RootCAName = "test"
)

type IntCAName string

const (
	IntCaNameServices            IntCAName = "services"
	IntCaNameIntranet            IntCAName = "intranet"
	IntCaNameMsEntraClientSecret IntCAName = "ms-entra-client-secret"
	IntCaNameTest                IntCAName = "test"
)

type ProfileNamespaceIDName string

const (
	ProfileNamespaceIDNameBuiltin ProfileNamespaceIDName = "builtin"
	ProfileNamespaceIDNameTenant  ProfileNamespaceIDName = "tenant"
)
