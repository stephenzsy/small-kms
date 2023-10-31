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
	"github.com/stephenzsy/small-kms/backend/utils"
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
	Policy        base.DocFullIdentifier   `json:"policy"`
	PolicyVersion HexDigest                `json:"policyVersion"`
	Created       base.NumericDate         `json:"iat"`
	NotBefore     base.NumericDate         `json:"nbf"`
	NotAfter      base.NumericDate         `json:"exp"`
	Flags         []CertificateFlag        `json:"flags"`
	KeyVaultStore CertDocKeyVaultStore     `json:"keyVaultStore"`
	Issuer        base.DocFullIdentifier   `json:"issuer"`
}

const (
	certDocQueryColumnThumbprintSHA1 = "c.keySpec.x5t"
	certDocQueryColumnCreated        = "c.iat"
	certDocQueryColumnNotAfter       = "c.exp"
)

type CertDocSigningPatch struct {
	KeySpec       key.SigningKeySpec     `json:"keySpec"`
	KeyVaultStore CertDocKeyVaultStore   `json:"keyVaultStore"`
	Issuer        base.DocFullIdentifier `json:"issuer"`
}

func GetKeyStoreName(nsKind base.NamespaceKind, nsID Identifier, policyIdentifier Identifier) string {
	return fmt.Sprintf("%s-%s-%s", nsKind, nsID.String(), policyIdentifier.String())
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
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindCert, base.UUIDIdentifier(certID))

	d.Status = CertificateStatusPending
	d.KeySpec = pDoc.KeySpec
	d.KeyExportable = pDoc.KeyExportable
	d.Subject = pDoc.Subject
	d.SANs = pDoc.SANs
	d.Flags = pDoc.Flags
	d.Policy = pDoc.GetStorageFullIdentifier()
	d.PolicyVersion = pDoc.Version
	d.KeyVaultStore.Name = GetKeyStoreName(nsKind, nsID, pDoc.ID)

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

func (d *CertDoc) PopulateModel(m *Certificate) {
	if d == nil || m == nil {
		return
	}
	d.PopulateModelRef(&m.CertificateRef)
	m.Jwk = key.JsonWebKey{
		Kty:              d.KeySpec.Kty,
		Crv:              d.KeySpec.Crv,
		E:                d.KeySpec.E,
		N:                d.KeySpec.N,
		X:                d.KeySpec.X,
		Y:                d.KeySpec.Y,
		X5t:              d.KeySpec.X5t,
		X5tS256:          d.KeySpec.X5tS256,
		KeyOperations:    &d.KeySpec.KeyOperations,
		KeyID:            d.KeySpec.KeyID,
		CertificateChain: d.KeySpec.CertificateChain,
	}
	if d.KeySpec.Alg != nil {
		m.Alg = *d.KeySpec.Alg
	}
	m.Subject = d.Subject
	m.Flags = d.Flags
	m.Attributes.Nbf = &d.NotBefore
	m.Attributes.Iat = &d.Created
	m.SubjectAlternativeNames = d.SANs
}

var _ base.ModelPopulater[Certificate] = (*CertDoc)(nil)

func (d *CertDoc) getSigningParams() kv.SigningParams {
	params := kv.SigningParams{
		SigAlg: azkeys.SignatureAlgorithm(*d.KeySpec.Alg),
	}
	if d.KeySpec.KeyID != nil {
		params.CertID = azcertificates.ID(*d.KeySpec.KeyID)
	}
	return params
}

func (d *CertDoc) getCSRProviderParams() kv.CSRProviderParams {
	params := kv.CSRProviderParams{
		CertName:      d.KeyVaultStore.Name,
		KeyProperties: azcertificates.KeyProperties{},
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
	certID := d.ID.UUID()
	cert := &x509.Certificate{
		SerialNumber: new(big.Int).SetBytes(certID[:]),
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
	nextKeySpec.E = patch.KeySpec.E
	nextKeySpec.N = patch.KeySpec.N
	nextKeySpec.X = patch.KeySpec.X
	nextKeySpec.Y = patch.KeySpec.Y
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

func (d *CertDoc) getIssuedX509Certificate() (*x509.Certificate, [][]byte, error) {
	chain := utils.MapSlice(d.KeySpec.CertificateChain, func(certBytes base.Base64RawURLEncodedBytes) []byte {
		return certBytes
	})
	cert, err := x509.ParseCertificate(chain[0])
	return cert, chain, err
}

func (d *CertDoc) getX509SignatureAlgorithm() (sa x509.SignatureAlgorithm) {
	switch *d.KeySpec.Alg {
	case key.JsonWebKeySignatureAlgorithmRS256:
		sa = x509.SHA256WithRSA
	case key.JsonWebKeySignatureAlgorithmRS384:
		sa = x509.SHA384WithRSA
	case key.JsonWebKeySignatureAlgorithmRS512:
		sa = x509.SHA512WithRSA
	case key.JsonWebKeySignatureAlgorithmPS256:
		sa = x509.SHA256WithRSAPSS
	case key.JsonWebKeySignatureAlgorithmPS384:
		sa = x509.SHA384WithRSAPSS
	case key.JsonWebKeySignatureAlgorithmPS512:
		sa = x509.SHA512WithRSAPSS
	case key.JsonWebKeySignatureAlgorithmES256:
		sa = x509.ECDSAWithSHA256
	case key.JsonWebKeySignatureAlgorithmES384:
		sa = x509.ECDSAWithSHA384
	case key.JsonWebKeySignatureAlgorithmES512:
		sa = x509.ECDSAWithSHA512
	}
	return sa
}

// func (d *CertDoc) x5cPEMBlocks() []pem.Block {
// 	return utils.MapSlice(d.KeySpec.CertificateChain, func(certBytes base.Base64RawURLEncodedBytes) pem.Block {
// 		return pem.Block{
// 			Type:  "CERTIFICATE",
// 			Bytes: certBytes,
// 		}
// 	})
// }
