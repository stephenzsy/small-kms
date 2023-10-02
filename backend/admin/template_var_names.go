package admin

type TemplateVarName string

const (
	TemplateVarNameDeviceURI              TemplateVarName = "device.uri"
	TemplateVarNameDeviceAltURI           TemplateVarName = "device.altUri"
	TemplateVarNameApplicationURI         TemplateVarName = "application.uri"
	TemplateVarNameApplicationAltURI      TemplateVarName = "application.altUri"
	TemplateVarNameServicePrincipalURI    TemplateVarName = "servicePrincipal.uri"
	TemplateVarNameServicePrincipalAltURI TemplateVarName = "servicePrincipal.altUri"
	TemplateVarNameGroupURI               TemplateVarName = "group.uri"
)
