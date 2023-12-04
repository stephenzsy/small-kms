package certda

import (
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type CertDocKeyVaultStore struct {
	Name          string `json:"name"`
	ID            string `json:"id"`
	SID           string `json:"sid"`
	KeyExportable bool   `json:"keyExportable"`
}

type certDoc struct {
	resdoc.ResourceDoc

	Status certmodels.CertificateStatus `json:"status"`

	JsonWebKey cloudkey.JsonWebKey `json:"jwk"`

	Subject          certmodels.CertificateSubject       `json:"subject"`
	SANs             *certmodels.SubjectAlternativeNames `json:"sans,omitempty"`
	PolicyIdentifier resdoc.DocIdentifier                `json:"policy"`
	PolicyVersion    []byte                              `json:"policyVersion"`

	Flags  []certmodels.CertificateFlag `json:"flags"`
	Issuer resdoc.DocIdentifier         `json:"issuer"`

	NotBefore resdoc.NumericDate `json:"nbf"`
	NotAfter  resdoc.NumericDate `json:"exp"`

	// enrolled certificate will have this as empty
	KeyVaultStore *CertDocKeyVaultStore `json:"keyVaultStore,omitempty"`

	SerialNumber []byte             `json:"serialNumber"`
	IssuedAt     resdoc.NumericDate `json:"iat"`

	Checksum []byte `json:"checksum"` // sha256 of the cloud certificate and critical fields
}
