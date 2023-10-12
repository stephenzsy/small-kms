package certtemplate

import (
	"github.com/stephenzsy/small-kms/backend/models"
)

type CertificateTemplateFlag string

const (
	CertTmplFlagRestrictKtyRsa CertificateTemplateFlag = "kty-rsa"
	CertTmplFlagDelegate       CertificateTemplateFlag = "delegate"
	CertTmplFlagTest           CertificateTemplateFlag = "test"
	CertTmplFlagHasKeyStore    CertificateTemplateFlag = "use-key-store"
	CertTmplFlagKeyExportable  CertificateTemplateFlag = "key-exportable"
)

type CertificateTemplateDocKeyProperties struct {
	// signature algorithm
	Kty      models.JwtKty  `json:"kty"`
	KeySize  *int           `json:"key_size,omitempty"`
	Crv      *models.JwtCrv `json:"crv,omitempty"`
	ReuseKey *bool          `json:"reuse_key,omitempty"`
}

type CertificateTemplateDocSubject struct {
	CN string  `json:"cn"`
	OU *string `json:"ou,omitempty"`
	O  *string `json:"o,omitempty"`
	C  *string `json:"c,omitempty"`
}

type CertificateTemplateDocSANs struct {
	EmailAddresses []string `json:"emails,omitempty"`
	URIs           []string `json:"uris,omitempty"`
}

type CertificateTemplateDocLifeTimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}
