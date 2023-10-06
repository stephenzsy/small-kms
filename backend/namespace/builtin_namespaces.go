package ns

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

type BuiltiInCertTemplateName string

const (
	CertTemplateNameDefault                   BuiltiInCertTemplateName = "default"
	CertTemplateNameDefaultIntranetAccess     BuiltiInCertTemplateName = "default-intranet-access"
	CertTemplateNameDefaultMsEntraClientCreds BuiltiInCertTemplateName = "default-ms-entra-client-creds"
)
