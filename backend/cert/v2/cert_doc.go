package cert

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"math/big"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
	"github.com/stephenzsy/small-kms/backend/utils/caldur"
)

type CertificateStatus string

type CertDocKeyVaultStore struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	SID  string `json:"sid"`
}

type CertDoc struct {
	resdoc.ResourceDoc

	Status           certmodels.CertificateStatus        `json:"status"`
	JsonWebKey       cloudkey.JsonWebSignatureKey        `json:"jwk"`
	KeyExportable    bool                                `json:"keyExportable"`
	Subject          certmodels.CertificateSubject       `json:"subject"`
	SANs             *certmodels.SubjectAlternativeNames `json:"sans,omitempty"`
	PolicyIdentifier resdoc.DocIdentifier                `json:"policy"`
	PolicyVersion    []byte                              `json:"policyVersion"`
	IssuedAt         jwt.NumericDate                     `json:"iat"`
	NotBefore        jwt.NumericDate                     `json:"nbf"`
	NotAfter         jwt.NumericDate                     `json:"exp"`
	Flags            []certmodels.CertificateFlag        `json:"flags"`
	KeyVaultStore    *CertDocKeyVaultStore               `json:"keyVaultStore,omitempty"`
	Issuer           resdoc.DocIdentifier                `json:"issuer"`
	Checksum         []byte                              `json:"checksum"` // sha256 of the cloud certificate and critical fields
}

type certDocSelfSignedGeneratePending struct {
	CertDoc
	rsaKeySize       int
	serialNumber     *big.Int
	templateX509Cert *x509.Certificate
	issuerX509Cert   *x509.Certificate
	issuerCertChain  []cloudkey.Base64RawURLEncodableBytes
	publicKey        crypto.PublicKey
	signer           crypto.Signer
	createCertResp   *azcertificates.CreateCertificateResponse
}

const (
	certDocQueryColStatus         = "c.status"
	certDocQueryColThumbprintSHA1 = "c.jwk.x5t"
	certDocQueryColIssuedAt       = "c.iat"
	certDocQueryColNotAfter       = "c.exp"
	certDocQueryColPolicy         = "c.policy"
)

// func GetKeyVaultStoreName(nsProvider models.NamespaceProvider, nsID string, policyID string) string {

// 	return fmt.Sprintf("c-%s-%s-%s", nsProvider, nsID, policyID)
// }

// upon success, this function craetes a key in keyvault
func (d *certDocSelfSignedGeneratePending) init(
	c ctx.RequestContext,
	nsProvider models.NamespaceProvider, nsID string,
	pDoc *CertPolicyDoc) error {
	certID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	d.ID = certID.String()
	d.serialNumber = new(big.Int).SetBytes(certID[:])

	d.Status = certmodels.CertificateStatusPending
	d.JsonWebKey.KeyType = pDoc.KeySpec.Kty
	d.JsonWebKey.Curve = pDoc.KeySpec.Crv
	if pDoc.KeySpec.KeySize != nil {
		d.rsaKeySize = *pDoc.KeySpec.KeySize
	}
	d.JsonWebKey.Alg = cloudkey.JsonWebSignatureAlgorithm(pDoc.KeySpec.Alg)
	d.JsonWebKey.KeyOperations = pDoc.KeySpec.KeyOperations
	d.KeyExportable = pDoc.KeyExportable
	d.Subject = pDoc.Subject
	d.SANs = pDoc.SANs
	d.Flags = pDoc.Flags
	d.PolicyIdentifier = pDoc.Identifier()
	d.PolicyVersion = pDoc.Version

	now := time.Now().Truncate(time.Second)

	d.NotBefore.Time = now
	d.NotAfter.Time = caldur.Shift(now, pDoc.ExpiryTime)

	azKeysClient := kv.GetAzKeyVaultService(c).AzKeysClient()
	d.templateX509Cert = d.generateCertificateTemplate()
	if pDoc.IssuerPolicy.IsEmpty() {
		d.issuerX509Cert = d.templateX509Cert
		ckParams, err := d.getAzCreateKeyParams()
		if err != nil {
			return err
		}
		d.KeyVaultStore = &CertDocKeyVaultStore{
			Name: kv.GetMaterialName(kv.MaterialNameKindCertificateKey, nsProvider, nsID, pDoc.ID),
		}
		ckResp, ck, err := cloudkeyaz.CreateCloudSignatureKey(c,
			azKeysClient, d.KeyVaultStore.Name, ckParams, d.JsonWebKey.Alg, true)
		if err != nil {
			return err
		}
		d.JsonWebKey.N = ckResp.Key.N
		d.JsonWebKey.E = ckResp.Key.E
		d.JsonWebKey.X = ckResp.Key.X
		d.JsonWebKey.Y = ckResp.Key.Y
		d.JsonWebKey.KeyID = string(*ckResp.Key.KID)
		d.publicKey = ck.Public()
		d.signer = ck
	} else {
		issuerPolicy, err := getCertificatePolicyInternal(c, pDoc.IssuerPolicy.NamespaceProvider, pDoc.IssuerPolicy.NamespaceID, pDoc.IssuerPolicy.ID)
		if err != nil {
			return err
		}
		signerCert, err := issuerPolicy.getIssuerCert(c)
		if err != nil {
			return err
		} else if signerCert.Status != certmodels.CertificateStatusIssued {
			return fmt.Errorf("issuer certificate is not issued")
		} else if time.Until(signerCert.NotAfter.Time) < 24*time.Hour {
			return fmt.Errorf("issuer certificate is expiring soon or has expired")
		}
		d.issuerCertChain = signerCert.JsonWebKey.CertificateChain
		d.issuerX509Cert, err = x509.ParseCertificate(signerCert.JsonWebKey.CertificateChain[0])
		d.KeyVaultStore = &CertDocKeyVaultStore{
			Name: kv.GetMaterialName(kv.MaterialNameKindCertificate, nsProvider, nsID, pDoc.ID),
		}
		if err != nil {
			return err
		}
		d.signer = cloudkeyaz.NewAzCloudSignatureKeyWithKID(c, azKeysClient, signerCert.JsonWebKey.KeyID, signerCert.JsonWebKey.Alg)

		// now needs public key from keyvault
		azCertClient := kv.GetAzKeyVaultService(c).AzCertificatesClient()
		createCertParams, err := d.getAzCreateCertParams()
		if err != nil {
			return err
		}
		resp, err := azCertClient.CreateCertificate(c, d.KeyVaultStore.Name, createCertParams, nil)
		if err != nil {
			return err
		}
		d.createCertResp = &resp
		csrParsed, err := x509.ParseCertificateRequest(resp.CSR)
		if err != nil {
			return err
		}
		d.publicKey = csrParsed.PublicKey

	}

	return nil
}

func (d *certDocSelfSignedGeneratePending) getAzCreateKeyParams() (params azkeys.CreateKeyParameters, err error) {
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

func (d *certDocSelfSignedGeneratePending) getAzCreateCertParams() (params azcertificates.CreateCertificateParameters, err error) {
	params.CertificateAttributes = &azcertificates.CertificateAttributes{
		Enabled:   to.Ptr(true),
		NotBefore: &d.NotBefore.Time,
		Expires:   &d.NotAfter.Time,
	}
	params.CertificatePolicy = &azcertificates.CertificatePolicy{
		KeyProperties: &azcertificates.KeyProperties{
			Exportable: &d.KeyExportable,
		},
		SecretProperties: &azcertificates.SecretProperties{
			ContentType: to.Ptr("application/x-pem-file"),
		},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			Subject: to.Ptr(d.Subject.String()),
		},
	}
	kp := params.CertificatePolicy.KeyProperties
	switch d.JsonWebKey.KeyType {
	case cloudkey.KeyTypeEC:
		kp.KeyType = to.Ptr(azcertificates.KeyTypeEC)
		switch d.JsonWebKey.Curve {
		case cloudkey.CurveNameP256:
			kp.Curve = to.Ptr(azcertificates.CurveNameP256)
		case cloudkey.CurveNameP384:
			kp.Curve = to.Ptr(azcertificates.CurveNameP384)
		case cloudkey.CurveNameP521:
			kp.Curve = to.Ptr(azcertificates.CurveNameP521)
		default:
			return params, cloudkey.ErrInvalidCurve
		}
	case cloudkey.KeyTypeRSA:
		kp.KeyType = to.Ptr(azcertificates.KeyTypeRSA)
		switch d.rsaKeySize {
		case 2048, 3072, 4096:
			kp.KeySize = to.Ptr(int32(d.rsaKeySize))
		}
	default:
		return params, cloudkey.ErrInvalidKeyType
	}
	return params, nil
}

func (d *certDocSelfSignedGeneratePending) collectSignedCert(c context.Context, cert []byte) (err error) {

	d.JsonWebKey.CertificateChain = append([]cloudkey.Base64RawURLEncodableBytes(nil), cert)
	d.JsonWebKey.CertificateChain = append(d.JsonWebKey.CertificateChain, d.issuerCertChain...)
	if d.createCertResp != nil {
		certClient := kv.GetAzKeyVaultService(c).AzCertificatesClient()
		resp, err := certClient.MergeCertificate(c, d.createCertResp.ID.Name(), azcertificates.MergeCertificateParameters{
			X509Certificates: utils.MapSlice(d.JsonWebKey.CertificateChain, func(e base.Base64RawURLEncodedBytes) []byte {
				return e
			}),
		}, nil)
		if err != nil {
			return err
		}
		d.JsonWebKey.KeyID = string(*resp.KID)
		d.KeyVaultStore.ID = string(*resp.ID)
		if resp.SID != nil {
			d.KeyVaultStore.SID = string(*resp.SID)
		}
	}

	sha1d := sha1.New()
	sha1d.Write(cert)
	d.JsonWebKey.ThumbprintSHA1 = sha1d.Sum(nil)
	sha256d := sha256.New()
	sha256d.Write(cert)
	d.JsonWebKey.ThumbprintSHA256 = sha256d.Sum(nil)
	d.Status = certmodels.CertificateStatusIssued
	d.IssuedAt.Time = time.Now().Truncate(time.Second)
	d.Checksum = d.calculateChecksum()
	return nil
}

func (d *certDocSelfSignedGeneratePending) calculateChecksum() []byte {
	digest := sha512.New384()
	// serial number
	digest.Write(d.serialNumber.Bytes())
	// subject
	io.WriteString(digest, d.Subject.String())
	// key and cert
	d.JsonWebKey.Digest(digest)
	if d.KeyExportable {
		digest.Write([]byte{1})
	} else {
		digest.Write([]byte{0})
	}
	// subject alternative names
	d.SANs.Digest(digest)
	// validity period
	if m, _ := d.NotBefore.MarshalJSON(); m != nil {
		digest.Write(m)
	}
	if m, _ := d.NotAfter.MarshalJSON(); m != nil {
		digest.Write(m)
	}
	// flags
	for _, v := range d.Flags {
		digest.Write([]byte(v))
	}
	if d.KeyVaultStore != nil {
		io.WriteString(digest, d.KeyVaultStore.SID)
		io.WriteString(digest, d.KeyVaultStore.ID)
	}
	return digest.Sum(nil)
}

func (d *CertDoc) ToRef() (m certmodels.CertificateRef) {
	m.Ref = d.ResourceDoc.ToRef()
	m.Thumbprint = d.JsonWebKey.ThumbprintSHA1.HexString()
	m.Status = d.Status
	m.PolicyIdentifier = d.PolicyIdentifier.String()
	m.Iat = &d.IssuedAt
	m.Exp = d.NotAfter
	return m
}

func (d *CertDoc) ToModel(includeJwk bool) (m certmodels.Certificate) {
	m.CertificateRef = d.ToRef()
	if includeJwk {
		m.Jwk = &keymodels.JsonWebSignatureKey{
			JsonWebKeyBase: cloudkey.JsonWebKeyBase{
				KeyType:          d.JsonWebKey.KeyType,
				Curve:            d.JsonWebKey.Curve,
				E:                d.JsonWebKey.E,
				N:                d.JsonWebKey.N,
				X:                d.JsonWebKey.X,
				Y:                d.JsonWebKey.Y,
				ThumbprintSHA1:   d.JsonWebKey.ThumbprintSHA1,
				ThumbprintSHA256: d.JsonWebKey.ThumbprintSHA256,
				KeyOperations:    d.JsonWebKey.KeyOperations,
				CertificateChain: d.JsonWebKey.CertificateChain,
			},
			Alg: d.JsonWebKey.Alg,
		}
	}
	m.Subject = d.Subject.String()
	m.Flags = d.Flags
	m.Nbf = d.NotBefore
	m.SubjectAlternativeNames = d.SANs
	return m
}

func (d *CertDoc) cleanupKeyVault(c context.Context) error {
	if d.KeyVaultStore != nil && d.KeyVaultStore.ID != "" {
		certClient := kv.GetAzKeyVaultService(c).AzCertificatesClient()
		cid := azcertificates.ID(d.KeyVaultStore.ID)
		_, err := certClient.UpdateCertificate(c, cid.Name(), cid.Version(), azcertificates.UpdateCertificateParameters{
			CertificateAttributes: &azcertificates.CertificateAttributes{
				Enabled: to.Ptr(false),
			},
		}, nil)
		if err != nil {
			return err
		}
	} else if d.JsonWebKey.KeyID != "" {
		kid := azkeys.ID(d.JsonWebKey.KeyID)
		azKeysClient := kv.GetAzKeyVaultService(c).AzKeysClient()
		_, err := azKeysClient.UpdateKey(c, kid.Name(), kid.Version(), azkeys.UpdateKeyParameters{
			KeyAttributes: &azkeys.KeyAttributes{
				Enabled: to.Ptr(false),
			},
		}, nil)
		if err != nil {
			err = base.HandleAzKeyVaultError(err)
			if !errors.Is(err, base.ErrAzKeyVaultItemNotFound) {
				return err
			}
		}
	}
	return nil
}

func (d *certDocSelfSignedGeneratePending) cleanupKeyVault(c context.Context) error {
	if d.createCertResp != nil {
		azCertClient := kv.GetAzKeyVaultService(c).AzCertificatesClient()
		_, err := azCertClient.DeleteCertificateOperation(c, d.createCertResp.ID.Name(), nil)
		if err != nil {
			return err
		}
	}
	return d.CertDoc.cleanupKeyVault(c)
}

func (d *certDocSelfSignedGeneratePending) generateCertificateTemplate() *x509.Certificate {

	cert := &x509.Certificate{
		SerialNumber: d.serialNumber,
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

	cert.SignatureAlgorithm = d.JsonWebKey.Alg.X509SignatureAlgorithm()

	return cert
}
