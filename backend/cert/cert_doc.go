package cert

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
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
	Policy        base.DocLocator          `json:"policy"`
	PolicyVersion HexDigest                `json:"policyVersion"`
	Created       base.NumericDate         `json:"iat"`
	NotBefore     base.NumericDate         `json:"nbf"`
	NotAfter      base.NumericDate         `json:"exp"`
	Flags         []CertificateFlag        `json:"flags"`
	KeyVaultStore CertDocKeyVaultStore     `json:"keyVaultStore"`
	Issuer        base.DocLocator          `json:"issuer"`
}

const (
	certDocQueryColumnThumbprintSHA1 = "c.keySpec.x5t"
	certDocQueryColumnCreated        = "c.iat"
	certDocQueryColumnNotAfter       = "c.exp"
)

type CertDocSigningPatch struct {
	KeySpec       key.SigningKeySpec   `json:"keySpec"`
	KeyVaultStore CertDocKeyVaultStore `json:"keyVaultStore"`
	Issuer        base.DocLocator      `json:"issuer"`
}

func GetKeyStoreName(nsKind base.NamespaceKind, nsID ID, policyIdentifier ID) string {
	return fmt.Sprintf("%s-%s-%s", nsKind, nsID, policyIdentifier)
}

func (d *CertDoc) init(
	c context.Context,
	nsKind base.NamespaceKind,
	nsID ID,
	pDoc *CertPolicyDoc) error {
	if d == nil {
		return nil
	}
	certID, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	d.BaseDoc.Init(nsKind, nsID, base.ResourceKindCert, base.IDFromUUID(certID))

	d.Status = CertificateStatusPending
	d.KeySpec = pDoc.KeySpec
	d.KeyExportable = pDoc.KeyExportable
	if d.Subject, err = pDoc.Subject.processTemplate(c); err != nil {
		return err
	}
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
		KeyType:          d.KeySpec.Kty,
		Curve:            d.KeySpec.Crv,
		E:                d.KeySpec.E,
		N:                d.KeySpec.N,
		X:                d.KeySpec.X,
		Y:                d.KeySpec.Y,
		ThumbprintSHA1:   d.KeySpec.X5t,
		ThumbprintSHA256: d.KeySpec.X5tS256,
		KeyOperations:    d.KeySpec.KeyOperations,
		KeyID:            *d.KeySpec.KeyID,
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

// func (d *CertDoc) x5cPEMBlocks() []pem.Block {
// 	return utils.MapSlice(d.KeySpec.CertificateChain, func(certBytes base.Base64RawURLEncodedBytes) pem.Block {
// 		return pem.Block{
// 			Type:  "CERTIFICATE",
// 			Bytes: certBytes,
// 		}
// 	})
// }
