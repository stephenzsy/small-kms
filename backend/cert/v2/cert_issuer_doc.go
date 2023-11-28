package cert

import (
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type CertIssuerDocACME struct {
}

type CertIssuerDoc struct {
	resdoc.ResourceDoc
	DisplayName string `json:"displayName"`

	ACME    *CertIssuerDocACME `json:"acme,omitempty"`
	Version []byte             `json:"version"`
}
