package cert

import (
	"context"
	"crypto/x509"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/key"
)

type CertificateStatus string

const (
	CertificateStatusPending CertificateStatus = "pending"
	CertificateStatusIssued  CertificateStatus = "issued"
	CertificateStatusError   CertificateStatus = "error"
)

type CertDocKeyVaultStore struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	SID  string `json:"sid"`
}

type CertDoc struct {
	base.BaseDoc

	Status        CertificateStatus        `json:"status"`
	KeySpec       key.SigningKeySpec       `json:"keySpec"`
	KeyExportable bool                     `json:"keyExportable"`
	Subject       CertificateSubject       `json:"subject"`
	SANs          *SubjectAlternativeNames `json:"sans,omitempty"`
	Policy        base.SLocator            `json:"policyLocator"`
	PolicyVersion HexDigest                `json:"policyVersion"`
	Created       base.NumericDate         `json:"iat"`
	NotBefore     base.NumericDate         `json:"nbf"`
	NotAfter      base.NumericDate         `json:"exp"`
	Flags         []CertificateFlag        `json:"flags"`
	KeyVaultStore CertDocKeyVaultStore     `json:"keyVaultStore"`
	Issuer        base.SLocator            `json:"issuer"`
}

type CertListQueryDoc struct {
	base.BaseDoc
	ThumbprintSHA1  base.Base64RawURLEncodedBytes `json:"x5t"`
	NotAfter        base.NumericDate              `json:"exp"`
	IssuerForPolicy *base.SLocator                `json:"issuerCertPolicyId,omitempty"`
}

const (
	certDocQueryColumnThumbprintSHA1 = "c.keySpec.x5t"
	certDocQueryColumnCreated        = "c.iat"
	certDocQueryColumnNotAfter       = "c.exp"
)

type CertDocSigningPatch struct {
	KeySpec       key.SigningKeySpec   `json:"keySpec"`
	KeyVaultStore CertDocKeyVaultStore `json:"keyVaultStore"`
	Issuer        base.SLocator        `json:"issuer"`
}

// GetStorageID implements base.CRUDDocHasCustomStorageID.
func (d *CertDoc) GetStorageID(context.Context) uuid.UUID {
	return d.ResourceIdentifier.UUID()
}

func (d *CertDoc) Init(
	nsKind base.NamespaceKind,
	nsID Identifier,
	pDoc *CertPolicyDoc) error {
	if d == nil {
		return nil
	}
	certID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	d.NamespaceKind = nsKind
	d.NamespaceIdentifier = nsID
	d.ResourceKind = base.ResourceKindCert
	d.ResourceIdentifier = base.UUIDIdentifier(certID)

	d.Status = CertificateStatusPending
	d.KeySpec = pDoc.KeySpec
	d.KeyExportable = pDoc.KeyExportable
	d.Subject = pDoc.Subject
	d.SANs = pDoc.SANs
	d.Flags = pDoc.Flags
	d.Policy = pDoc.GetPersistedSLocator()
	d.PolicyVersion = pDoc.Version
	d.KeyVaultStore.Name =
		fmt.Sprintf("%s-%s-%s", d.NamespaceKind, d.NamespaceIdentifier.String(), pDoc.ResourceIdentifier.String())

	now := time.Now()
	d.Created = *jwt.NewNumericDate(now)
	d.NotBefore = d.Created
	d.NotAfter = *jwt.NewNumericDate(base.AddPeriod(now, pDoc.ExpiryTime))

	return nil
}

func (d *CertDoc) PopulateModelRef(m *CertificateRef) {
	if d == nil || m == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&m.ResourceReference)
	if d.KeySpec.X5t != nil {
		m.Thumbprint = d.KeySpec.X5t.HexString()
	}
	m.Attributes.Exp = &d.NotAfter
}

func (d *CertListQueryDoc) PopulateModelRef(m *CertificateRef) {
	if d == nil || m == nil {
		return
	}
	d.BaseDoc.PopulateModelRef(&m.ResourceReference)
	m.Thumbprint = d.ThumbprintSHA1.HexString()
	m.Attributes.Exp = &d.NotAfter
	m.IssuerForPolicy = d.IssuerForPolicy
}

func (d *CertDoc) PopulateModel(m *Certificate) {
	if d == nil || m == nil {
		return
	}
	d.PopulateModelRef(&m.CertificateRef)
	if d.KeySpec.Alg != nil {
		m.Alg = *d.KeySpec.Alg
	}
	if d.KeySpec.X5t != nil {
		m.X5t = *d.KeySpec.X5t
	}
	if d.KeySpec.X5tS256 != nil {
		m.X5tS256 = *d.KeySpec.X5tS256
	}
	m.Subject = d.Subject
	m.Flags = d.Flags
	m.Attributes.Nbf = &d.NotBefore
	m.Attributes.Iat = &d.Created
	m.SubjectAlternativeNames = d.SANs
}

var _ base.CRUDDocHasCustomStorageID = (*CertDoc)(nil)

var _ base.ModelRefPopulater[certificateRefComposed] = (*CertListQueryDoc)(nil)
var _ base.ModelPopulater[Certificate] = (*CertDoc)(nil)

func (d *CertDoc) getSigningParams() kv.SigningParams {
	certID := d.KeyVaultStore.Name
	params := kv.SigningParams{
		CertName:      certID,
		KeyProperties: azcertificates.KeyProperties{},
		SigAlg:        azkeys.SignatureAlgorithm(*d.KeySpec.Alg),
	}
	switch d.KeySpec.Kty {
	case key.JsonWebKeyTypeRSA:
		params.KeyProperties.KeyType = to.Ptr(azcertificates.KeyTypeRSA)
		params.KeyProperties.KeySize = d.KeySpec.KeySize
	case key.JsonWebKeyTypeEC:
		params.KeyProperties.KeyType = to.Ptr(azcertificates.KeyTypeEC)
		switch *d.KeySpec.Crv {
		case key.JsonWebKeyCurveNameP256:
			params.KeyProperties.Curve = to.Ptr(azcertificates.CurveNameP256)
		case key.JsonWebKeyCurveNameP384:
			params.KeyProperties.Curve = to.Ptr(azcertificates.CurveNameP384)
		case key.JsonWebKeyCurveNameP521:
			params.KeyProperties.Curve = to.Ptr(azcertificates.CurveNameP521)
		}
	}
	params.KeyProperties.Exportable = &d.KeyExportable
	return params
}

func (d *CertDoc) getX509CertTemplate() *x509.Certificate {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(0).SetBytes(d.ResourceIdentifier.Bytes()),
		Subject:      d.Subject.ToPkixName(),
		NotBefore:    d.NotBefore.Time,
		NotAfter:     d.NotAfter.Time,
	}

	if slices.Contains(d.Flags, CertificateFlagCA) {
		cert.KeyUsage |= x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature
		cert.BasicConstraintsValid = true
		cert.IsCA = true
		if slices.Contains(d.Flags, CertificateFlagRootCA) {
			cert.MaxPathLen = 1
			cert.MaxPathLenZero = false
		} else {
			cert.MaxPathLenZero = true
		}
	} else {
		cert.KeyUsage |= x509.KeyUsageDigitalSignature
		if slices.Contains(d.Flags, CertificateFlagServerAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
		}
		if slices.Contains(d.Flags, CertificateFlagClientAuth) {
			cert.ExtKeyUsage = append(cert.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		}
		if d.KeySpec.Kty == key.JsonWebKeyTypeRSA {
			cert.KeyUsage |= x509.KeyUsageKeyEncipherment
		}
	}

	switch *d.KeySpec.Alg {
	case key.JsonWebKeySignatureAlgorithmRS256:
		cert.SignatureAlgorithm = x509.SHA256WithRSA
	case key.JsonWebKeySignatureAlgorithmRS384:
		cert.SignatureAlgorithm = x509.SHA384WithRSA
	case key.JsonWebKeySignatureAlgorithmRS512:
		cert.SignatureAlgorithm = x509.SHA512WithRSA
	case key.JsonWebKeySignatureAlgorithmPS256:
		cert.SignatureAlgorithm = x509.SHA256WithRSAPSS
	case key.JsonWebKeySignatureAlgorithmPS384:
		cert.SignatureAlgorithm = x509.SHA384WithRSAPSS
	case key.JsonWebKeySignatureAlgorithmPS512:
		cert.SignatureAlgorithm = x509.SHA512WithRSAPSS
	case key.JsonWebKeySignatureAlgorithmES256:
		cert.SignatureAlgorithm = x509.ECDSAWithSHA256
	case key.JsonWebKeySignatureAlgorithmES384:
		cert.SignatureAlgorithm = x509.ECDSAWithSHA384
	case key.JsonWebKeySignatureAlgorithmES512:
		cert.SignatureAlgorithm = x509.ECDSAWithSHA512
	}

	if d.SANs != nil {
		cert.DNSNames = d.SANs.DNSNames
		cert.EmailAddresses = d.SANs.Emails
		cert.IPAddresses = d.SANs.IPAddresses
	}

	return cert
}

func (d *CertDoc) applyPatch(c context.Context,
	docService base.AzCosmosCRUDDocService,
	patch *CertDocSigningPatch) error {
	nextKeySpec := d.KeySpec
	nextKeySpec.CertificateChain = patch.KeySpec.CertificateChain
	nextKeySpec.KeyID = patch.KeySpec.KeyID
	nextKeySpec.X5t = patch.KeySpec.X5t
	nextKeySpec.X5tS256 = patch.KeySpec.X5tS256

	patchOps := azcosmos.PatchOperations{}
	patchOps.AppendSet("/keySpec", nextKeySpec)
	patchOps.AppendSet("/keyVaultStore", patch.KeyVaultStore)
	patchOps.AppendSet("/issuer", patch.Issuer)
	patchOps.AppendSet("/status", CertificateStatusIssued)
	err := docService.Patch(c, d, patchOps, &azcosmos.ItemOptions{
		IfMatchEtag: d.ETag,
	})
	if err != nil {
		return err
	}
	d.KeySpec = nextKeySpec
	d.KeyVaultStore = patch.KeyVaultStore
	d.Issuer = patch.Issuer
	return nil
}
