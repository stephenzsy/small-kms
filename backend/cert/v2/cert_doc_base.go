package cert

import (
	"context"
	"crypto"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	"github.com/stephenzsy/small-kms/backend/resdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
	"github.com/stephenzsy/small-kms/backend/utils/caldur"
)

type CertCSR interface {
	PublicKey() (crypto.PublicKey, error)
	X509CSRBytes() []byte
}

type CertDocument interface {
	resdoc.ResourceDocument
	ToModel(includeKey bool) *certmodels.Certificate
	X509Certificate() (*x509.Certificate, error)
	GetCertificateBytes() []byte
	GetJsonWebKey() *cloudkey.JsonWebKey
	GetStatus() certmodels.CertificateStatus
	IsExpired() bool
	GetNotBefore() time.Time
	GetNotAfter() time.Time
	KeyVaultSecretID() string
}

type CertDocumentPending interface {
	CertDocument
	Authorize(c context.Context) (bool, error)
	CreateCertificate(c ctx.RequestContext, csr CertCSR) ([][]byte, error)
	GetCertificateRequest(c ctx.RequestContext, skipCheckExisting bool) (CertCSR, error)
	CollectSignedCertificate(c ctx.RequestContext, der [][]byte) error
}

type CertDocKeyVaultStore struct {
	Name string `json:"name"`
	ID   string `json:"id,omitempty"`
	SID  string `json:"sid,omitempty"`
}

type certDocBase struct {
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

// KeyVaultSecretID implements CertDocument.
func (d *certDocBase) KeyVaultSecretID() string {
	if d.KeyVaultStore == nil {
		return ""
	}
	return d.KeyVaultStore.SID
}

type certDocPending struct {
	certDocBase
	certUUID   uuid.UUID
	rsaKeySize int32
}

// GetCertificateBytes implements CertDocument.
func (doc *certDocBase) GetCertificateBytes() []byte {
	return doc.JsonWebKey.CertificateChain[0]
}

// GetNotAfter implements CertDocument.
func (doc *certDocBase) GetNotAfter() time.Time {
	return doc.NotAfter.Time
}

// GetNotBefore implements CertDocument.
func (doc *certDocBase) GetNotBefore() time.Time {
	return doc.NotBefore.Time
}

// IsExpired implements CertDocument.
func (doc *certDocBase) IsExpired() bool {
	return doc.NotAfter.Before(time.Now())
}

// GetStatus implements CertDocument.
func (doc *certDocBase) GetStatus() certmodels.CertificateStatus {
	return doc.Status
}

// PublicJWK implements CertDocument.
func (doc *certDocBase) GetJsonWebKey() *cloudkey.JsonWebKey {
	return &doc.JsonWebKey
}

// X509Certificate implements CertDocument.
func (d *certDocBase) X509Certificate() (*x509.Certificate, error) {
	if len(d.JsonWebKey.CertificateChain) == 0 {
		return nil, fmt.Errorf("%w: certificate not issued", base.ErrResponseStatusBadRequest)
	}
	return x509.ParseCertificate(d.JsonWebKey.CertificateChain[0])
}

func (d *certDocPending) init(
	c ctx.RequestContext,
	nsProvider models.NamespaceProvider, nsID string,
	pDoc *CertPolicyDoc,
	publicJwk *cloudkey.JsonWebKey) (err error) {

	d.certUUID, err = uuid.NewRandom()
	if err != nil {
		return err
	}
	d.PartitionKey.NamespaceProvider = nsProvider
	d.PartitionKey.NamespaceID = nsID
	d.PartitionKey.ResourceProvider = models.ResourceProviderCert
	d.ID = d.certUUID.String()
	d.Status = certmodels.CertificateStatusPending

	d.JsonWebKey.KeyType = pDoc.KeySpec.Kty
	d.JsonWebKey.Curve = pDoc.KeySpec.Crv
	d.JsonWebKey.Alg = pDoc.KeySpec.Alg
	d.JsonWebKey.KeyOperations = pDoc.KeySpec.KeyOperations
	d.JsonWebKey.Extractable = pDoc.KeySpec.Extractable
	if publicJwk == nil {
		// cloud mastered key

		if pDoc.KeySpec.KeySize != nil {
			d.rsaKeySize = int32(*pDoc.KeySpec.KeySize)
		}

		var materialName string
		if nsProvider == models.NamespaceProviderRootCA {
			materialName = kv.GetMaterialName(kv.MaterialNameKindCertificateKey, nsProvider, nsID, pDoc.ID)
		} else {
			materialName = kv.GetMaterialName(kv.MaterialNameKindCertificate, nsProvider, nsID, pDoc.ID)
		}
		d.KeyVaultStore = &CertDocKeyVaultStore{
			Name: materialName,
		}
	} else {
		if d.JsonWebKey.KeyType != publicJwk.KeyType {
			return fmt.Errorf("%w: public key type does not match", base.ErrResponseStatusBadRequest)
		}
		switch d.JsonWebKey.KeyType {
		case cloudkey.KeyTypeRSA:
			if publicJwk.N.BitLen() != *pDoc.KeySpec.KeySize {
				return fmt.Errorf("%w: public key size does not match", base.ErrResponseStatusBadRequest)
			}
			d.JsonWebKey.N = publicJwk.N
			d.JsonWebKey.E = publicJwk.E
		case cloudkey.KeyTypeEC:
			if d.JsonWebKey.Curve != publicJwk.Curve {
				return fmt.Errorf("%w: public key curve does not match", base.ErrResponseStatusBadRequest)
			}
			d.JsonWebKey.X = publicJwk.X
			d.JsonWebKey.Y = publicJwk.Y
		default:
			return fmt.Errorf("%w: invalid key type", base.ErrResponseStatusBadRequest)
		}
	}

	d.Subject, err = processSubjectTemplate(c, pDoc.Subject)
	if err != nil {
		return err
	}
	d.SANs = pDoc.SANs
	d.Flags = pDoc.Flags
	d.PolicyIdentifier = pDoc.Identifier()
	d.PolicyVersion = pDoc.Version

	now := time.Now().Truncate(time.Second)

	d.NotBefore.Time = now
	d.NotAfter.Time = caldur.Shift(now, pDoc.ExpiryTime)

	return nil
}

// Authorize implements CertDocument.
func (*certDocBase) Authorize(c context.Context) (bool, error) {
	return true, nil
}

type kvCreateCertCSR struct {
	*azcertificates.CertificateOperation
}

// X509CSRBytes implements CertCSR.
func (csr *kvCreateCertCSR) X509CSRBytes() []byte {
	return csr.CSR
}

// PublicKey implements CertCSR.
func (csr *kvCreateCertCSR) PublicKey() (crypto.PublicKey, error) {

	parsed, err := x509.ParseCertificateRequest(csr.CSR)
	if err != nil {
		return nil, err
	}
	return parsed.PublicKey, nil
}

var _ CertCSR = (*kvCreateCertCSR)(nil)

type enrollPublicKeyCSR struct {
	pubKey crypto.PublicKey
}

// PublicKey implements CertCSR.
func (csr *enrollPublicKeyCSR) PublicKey() (crypto.PublicKey, error) {
	return csr.pubKey, nil
}

// X509CSRBytes implements CertCSR.
func (*enrollPublicKeyCSR) X509CSRBytes() []byte {
	return nil
}

var _ CertCSR = (*enrollPublicKeyCSR)(nil)

// GenerateCertificateRequest implements CertDocument.
func (doc *certDocPending) GetCertificateRequest(c ctx.RequestContext, skipCheckExisting bool) (CertCSR, error) {
	if doc.KeyVaultStore == nil {
		return &enrollPublicKeyCSR{doc.JsonWebKey.PublicKey()}, nil
	}
	client := kv.GetAzKeyVaultService(c).AzCertificatesClient()
	// if !skipCheckExisting {
	// 	if resp, err := client.GetCertificateOperation(c, doc.KeyVaultStore.Name, nil); err != nil {
	// 		err = kv.HandleAzKeyVaultError(err)
	// 		if !errors.Is(err, kv.ErrAzKeyVaultItemNotFound) {
	// 			return nil, err
	// 		}
	// 	} else if resp.CertificateOperation.Status == azcertificates. {
	// 		return &kvCreateCertCSR{&resp.CertificateOperation}, nil
	// 	}
	// }
	params, err := doc.getAzCreateCertParams()
	if err != nil {
		return nil, err
	}
	resp, err := client.CreateCertificate(c, doc.KeyVaultStore.Name, params, nil)
	if err != nil {
		return nil, err
	}
	return &kvCreateCertCSR{&resp.CertificateOperation}, nil
}

func (d *certDocPending) getAzCreateCertParams() (params azcertificates.CreateCertificateParameters, err error) {
	params.CertificateAttributes = &azcertificates.CertificateAttributes{
		Enabled:   to.Ptr(true),
		NotBefore: &d.NotBefore.Time,
		Expires:   &d.NotAfter.Time,
	}
	params.CertificatePolicy = &azcertificates.CertificatePolicy{
		KeyProperties: &azcertificates.KeyProperties{
			Exportable: d.JsonWebKey.Extractable,
		},
		SecretProperties: &azcertificates.SecretProperties{
			ContentType: to.Ptr("application/x-pem-file"),
		},
		X509CertificateProperties: &azcertificates.X509CertificateProperties{
			Subject: to.Ptr(d.Subject.String()),
		},
	}
	if d.SANs != nil {
		params.CertificatePolicy.X509CertificateProperties.SubjectAlternativeNames = &azcertificates.SubjectAlternativeNames{
			DNSNames: to.SliceOfPtrs(d.SANs.DNSNames...),
		}
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
			kp.KeySize = &d.rsaKeySize
		default:
			return params, cloudkey.ErrInvalidKeySize
		}
	default:
		return params, cloudkey.ErrInvalidKeyType
	}
	return params, nil
}

func (d *certDocBase) CollectSignedCertificate(c ctx.RequestContext, der [][]byte) error {

	parsed, err := x509.ParseCertificate(der[0])
	if err != nil {
		return err
	}
	if err := d.JsonWebKey.SetPublicKey(parsed.PublicKey); err != nil {
		return err
	}
	d.JsonWebKey.CertificateChain = utils.MapSlice(der, func(b []byte) cloudkey.Base64RawURLEncodableBytes { return b })
	switch {
	case parsed.KeyUsage&x509.KeyUsageDigitalSignature != 0:
		d.JsonWebKey.KeyOperations = append(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationSign, cloudkey.JsonWebKeyOperationVerify)
	case parsed.KeyUsage&x509.KeyUsageKeyEncipherment != 0:
		d.JsonWebKey.KeyOperations = append(d.JsonWebKey.KeyOperations, cloudkey.JsonWebKeyOperationWrapKey, cloudkey.JsonWebKeyOperationDeriveKey)
	}
	sha1d := sha1.New()
	sha1d.Write(der[0])
	d.JsonWebKey.ThumbprintSHA1 = sha1d.Sum(nil)
	sha256d := sha256.New()
	sha256d.Write(der[0])
	d.JsonWebKey.ThumbprintSHA256 = sha256d.Sum(nil)

	d.IssuedAt.Time = time.Now().Truncate(time.Second)
	d.Status = certmodels.CertificateStatusIssued
	d.SerialNumber = parsed.SerialNumber.Bytes()
	d.NotBefore.Time = parsed.NotBefore
	d.NotAfter.Time = parsed.NotAfter
	d.SANs = &certmodels.SubjectAlternativeNames{
		DNSNames:    parsed.DNSNames,
		Emails:      parsed.EmailAddresses,
		IPAddresses: parsed.IPAddresses,
	}
	d.Subject = certmodels.CertificateSubject{
		CommonName: parsed.Subject.CommonName,
	}
	d.Flags = make([]certmodels.CertificateFlag, 0, len(d.Flags))
	for _, ext := range parsed.ExtKeyUsage {
		switch ext {
		case x509.ExtKeyUsageServerAuth:
			d.Flags = append(d.Flags, certmodels.CertificateFlagServerAuth)
		case x509.ExtKeyUsageClientAuth:
			d.Flags = append(d.Flags, certmodels.CertificateFlagClientAuth)
		}
	}

	if d.KeyVaultStore != nil && d.PartitionKey.NamespaceProvider != models.NamespaceProviderRootCA {
		certClient := kv.GetAzKeyVaultService(c).AzCertificatesClient()
		resp, err := certClient.MergeCertificate(c, d.KeyVaultStore.Name, azcertificates.MergeCertificateParameters{
			X509Certificates: der,
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
	return nil
}

func (d *certDocBase) ToModel(includeKey bool) *certmodels.Certificate {
	if d == nil {
		return nil
	}
	m := &certmodels.Certificate{
		CertificateRef: certmodels.CertificateRef{
			Ref: d.ResourceDoc.ToRef(),
			CertificateRefFields: certmodels.CertificateRefFields{
				Exp:              d.NotAfter,
				Status:           d.Status,
				Thumbprint:       d.JsonWebKey.ThumbprintSHA1.HexString(),
				PolicyIdentifier: d.PolicyIdentifier.String(),
			},
		},
		CertificateFields: certmodels.CertificateFields{
			Identifier:              d.Identifier().String(),
			IssuerIdentifier:        d.Issuer.String(),
			SerialNumber:            hex.EncodeToString(d.SerialNumber),
			Nbf:                     d.NotBefore,
			Subject:                 d.Subject.String(),
			SubjectAlternativeNames: d.SANs,
			Flags:                   d.Flags,
		},
	}
	if !d.IssuedAt.Time.IsZero() {
		m.Iat = &d.IssuedAt
	}
	if includeKey {
		m.Jwk = &d.JsonWebKey
		if d.KeyVaultStore != nil {
			m.KeyVaultCertificateID = d.KeyVaultStore.ID
			m.KeyVaultSecretID = d.KeyVaultStore.SID
		}
	}
	return m
}

var _ CertDocument = (*certDocBase)(nil)
