package cert

import (
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type CertIssuerDocACME struct {
	DirectoryURL   string   `json:"directoryUrl"`
	AccountURI     string   `json:"accountUri"`
	AccountKeyID   string   `json:"accountKeyId"`
	AccountContact []string `json:"accountContact"`
	AccountStatus  string   `json:"accountStatus"`
}

type CertIssuerDoc struct {
	resdoc.ResourceDoc
	DisplayName string `json:"displayName"`

	ACME    *CertIssuerDocACME `json:"acme,omitempty"`
	Version []byte             `json:"version"`
}

func (d *CertIssuerDoc) ToModel() *certmodels.CertificateExternalIssuer {
	if d == nil {
		return nil
	}
	m := &certmodels.CertificateExternalIssuer{
		Ref: d.ResourceDoc.ToRef(),
		CertificateExternalIssuerFields: certmodels.CertificateExternalIssuerFields{
			Acme: &certmodels.CertificateExternalIssuerAcme{
				AccountKeyID: d.ACME.AccountKeyID,
				AccountURL:   d.ACME.AccountURI,
				//AzureDNSZoneResourceID: d.ACME.AzureDNSZoneResourceID,
				DirectoryURL: d.ACME.DirectoryURL,
				Contacts:     d.ACME.AccountContact,
			},
		},
	}
	return m
}
