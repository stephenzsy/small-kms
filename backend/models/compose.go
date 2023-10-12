package models

import "github.com/stephenzsy/small-kms/backend/shared"

type Identifier = shared.Identifier

type ProfileRefComposed struct {
	shared.ResourceRef
	ProfileRefFields
}

type ProfileComposed = ProfileRefComposed

type CertificateTemplateRefComposed struct {
	shared.ResourceRef
	CertificateTemplateRefFields
}

type CertificateTemplateComposed struct {
	CertificateTemplateRefComposed
	CertificateTemplateFields
}

// Deprecated: Use shared.CertificateInfo instead
type CertificateInfoComposed = shared.CertificateInfo

type ServiceConfigComposed struct {
	shared.ResourceRef
	ServiceConfigFields
}
