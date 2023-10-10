package models

type ProfileRefComposed struct {
	ResourceRef
	ProfileRefFields
}

type ProfileComposed = ProfileRefComposed

type CertificateTemplateRefComposed struct {
	ResourceRef
	CertificateTemplateRefFields
}

type CertificateTemplateComposed struct {
	CertificateTemplateRefComposed
	CertificateTemplateFields
}

type CertificateRefComposed struct {
	ResourceRef
	CertificateRefFields
}

type CertificateInfoComposed struct {
	CertificateRefComposed
	CertificateInfoFields
}

type ServiceConfigComposed struct {
	ResourceRef
	ServiceConfigFields
}
