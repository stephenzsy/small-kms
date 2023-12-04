package cert

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type certDocInternal struct {
	certDocPending
}

func (doc *certDocInternal) init(c ctx.RequestContext,
	nsProvider models.NamespaceProvider, nsID string,
	pDoc *CertPolicyDoc, publicKey *cloudkey.JsonWebKey) (err error) {
	if err = doc.certDocPending.init(c, nsProvider, nsID, pDoc, publicKey); err != nil {
		return err
	}
	if doc.PartitionKey.NamespaceProvider == models.NamespaceProviderRootCA {
		doc.Issuer = doc.Identifier()
	} else {
		issuerPolicy, err := GetCertificatePolicyInternal(c, pDoc.IssuerPolicy.NamespaceProvider, pDoc.IssuerPolicy.NamespaceID, pDoc.IssuerPolicy.ID)
		if err != nil {
			return err
		}
		doc.Issuer, err = issuerPolicy.getIssuerCertIdentifier(c)
		if err != nil {
			return err
		}
	}
	return err
}

// CreateCertificate implements CertDocument.
func (doc *certDocInternal) CreateCertificate(c ctx.RequestContext, csr CertCSR) ([][]byte, error) {
	template := doc.getCertificateTemplate()
	var issuerCert *x509.Certificate
	var signer crypto.Signer
	azKeysClient := kv.GetAzKeyVaultService(c).AzKeysClient()
	var signerChain [][]byte
	var publicKey crypto.PublicKey
	if doc.PartitionKey.NamespaceProvider == models.NamespaceProviderRootCA {
		issuerCert = template

		ckParams, err := doc.getAzCreateKeyParams()
		if err != nil {
			return nil, err
		}
		sigAlg := cloudkey.JsonWebSignatureAlgorithm(doc.JsonWebKey.Alg)
		ckResp, ck, err := cloudkeyaz.CreateCloudSignatureKey(c,
			azKeysClient, doc.KeyVaultStore.Name, ckParams, sigAlg, true)
		if err != nil {
			return nil, err
		}
		doc.JsonWebKey.KeyID = string(*ckResp.Key.KID)
		publicKey = ck.Public()
		signer = ck
		doc.Issuer = doc.Identifier()
		template.SignatureAlgorithm = sigAlg.X509SignatureAlgorithm()
	} else {
		issuerCertDoc, err := GetCertificateInternal(c, doc.Issuer.NamespaceProvider, doc.Issuer.NamespaceID, doc.Issuer.ID)
		if err != nil {
			return nil, err
		} else if issuerCertDoc.GetStatus() != certmodels.CertificateStatusIssued {
			return nil, fmt.Errorf("issuer certificate is not issued")
		} else if time.Until(issuerCertDoc.GetNotAfter()) < 24*time.Hour {
			return nil, fmt.Errorf("issuer certificate is expiring soon or has expired")
		}
		issuerCert, err = issuerCertDoc.X509Certificate()
		if err != nil {
			return nil, err
		}
		issuerJwk := issuerCertDoc.GetJsonWebKey()
		sigAlg := cloudkey.JsonWebSignatureAlgorithm(issuerJwk.Alg)

		signer = cloudkeyaz.NewAzCloudSignatureKeyWithKID(
			c, azKeysClient, issuerJwk.KeyID,
			sigAlg,
			true,
			issuerJwk.PublicKey())
		signerChain = utils.MapSlice(issuerCertDoc.GetJsonWebKey().CertificateChain, func(b cloudkey.Base64RawURLEncodableBytes) []byte { return b })
		publicKey, err = csr.PublicKey()
		if err != nil {
			return nil, err
		}
		template.SignatureAlgorithm = sigAlg.X509SignatureAlgorithm()

	}
	signed, err := x509.CreateCertificate(rand.Reader,
		template,
		issuerCert,
		publicKey,
		signer)
	if err != nil {
		return nil, err
	}
	der := make([][]byte, 1, len(signerChain)+1)
	der[0] = signed
	der = append(der, signerChain...)
	return der, nil
}

func (d *certDocInternal) getCertificateTemplate() *x509.Certificate {

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(0).SetBytes(d.certUUID[:]),
		Subject:      d.Subject.ToPkixName(),
		NotBefore:    d.NotBefore.Time,
		NotAfter:     d.NotAfter.Time,
	}

	if d.PartitionKey.NamespaceProvider == models.NamespaceProviderRootCA ||
		d.PartitionKey.NamespaceProvider == models.NamespaceProviderIntermediateCA {
		cert.KeyUsage |= x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature
		cert.BasicConstraintsValid = true
		cert.IsCA = true
		if d.PartitionKey.NamespaceProvider == models.NamespaceProviderRootCA {
			cert.MaxPathLen = 1
			cert.MaxPathLenZero = false
		} else {
			cert.MaxPathLenZero = true
		}
	} else {
		cert.KeyUsage |= x509.KeyUsageDigitalSignature
		if slices.Contains(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationWrapKey) &&
			slices.Contains(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationUnwrapKey) {
			cert.KeyUsage |= x509.KeyUsageKeyEncipherment
		}
		if slices.Contains(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationEncrypt) &&
			slices.Contains(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationDecrypt) {
			cert.KeyUsage |= x509.KeyUsageDataEncipherment
		}
		if slices.Contains(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationDeriveKey) &&
			slices.Contains(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationDeriveBits) {
			cert.KeyUsage |= x509.KeyUsageKeyAgreement
		}
		if slices.Contains(d.Flags, certmodels.CertificateFlagServerAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
		}
		if slices.Contains(d.Flags, certmodels.CertificateFlagClientAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		}
	}

	if d.SANs != nil {
		cert.DNSNames = d.SANs.DNSNames
		cert.EmailAddresses = d.SANs.Emails
		cert.IPAddresses = d.SANs.IPAddresses
	}

	return cert
}

func (d *certDocInternal) getAzCreateKeyParams() (params azkeys.CreateKeyParameters, err error) {
	switch d.JsonWebKey.KeyType {
	case cloudkey.KeyTypeEC:
		params.Kty = to.Ptr(azkeys.KeyTypeEC)
		switch d.JsonWebKey.Curve {
		case cloudkey.CurveNameP256:
			params.Curve = to.Ptr(azkeys.CurveNameP256)
		case cloudkey.CurveNameP384:
			params.Curve = to.Ptr(azkeys.CurveNameP384)
		case cloudkey.CurveNameP521:
			params.Curve = to.Ptr(azkeys.CurveNameP521)
		default:
			return params, cloudkey.ErrInvalidCurve
		}
	case cloudkey.KeyTypeRSA:
		params.Kty = to.Ptr(azkeys.KeyTypeRSA)
		switch d.rsaKeySize {
		case 2048, 3072, 4096:
			params.KeySize = to.Ptr(int32(d.rsaKeySize))
		}
	default:
		return params, cloudkey.ErrInvalidKeyType
	}
	// keyops
	params.KeyOps = make([]*azkeys.KeyOperation, len(d.JsonWebKey.KeyOperations))
	for i, keyOp := range d.JsonWebKey.KeyOperations {
		params.KeyOps[i] = to.Ptr(azkeys.KeyOperation(keyOp))
	}
	// exportable
	params.KeyAttributes = &azkeys.KeyAttributes{
		Exportable: &d.KeyExportable,
		NotBefore:  &d.NotBefore.Time,
		Expires:    &d.NotAfter.Time,
		Enabled:    to.Ptr(true),
	}
	return params, nil
}

var _ CertDocumentPending = (*certDocInternal)(nil)
