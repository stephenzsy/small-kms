package cert

import (
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"golang.org/x/crypto/acme"
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

	acmeClient *acme.Client
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

func (d *CertIssuerDoc) ACMEClient(c ctx.RequestContext) (*acme.Client, error) {
	if d.acmeClient == nil {
		// load cloudKey
		keyDoc, err := key.GetKeyInternal(c, models.NamespaceProviderExternalCA, d.PartitionKey.NamespaceID, d.ACME.AccountKeyID)
		if err != nil {
			return nil, err
		}
		ck := cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, kv.GetAzKeyVaultService(c).AzKeysClient(), keyDoc.KeyID, cloudkey.SignatureAlgorithmES384, true, keyDoc.PublicKey())
		d.acmeClient = &acme.Client{
			DirectoryURL: d.ACME.DirectoryURL,
			Key:          ck,
		}
	}

	return d.acmeClient, nil
}
