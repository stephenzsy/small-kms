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
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/base"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	cloudkeyaz "github.com/stephenzsy/small-kms/backend/cloud/key/az"
	kv "github.com/stephenzsy/small-kms/backend/internal/keyvault"
	"github.com/stephenzsy/small-kms/backend/models"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
	keymodels "github.com/stephenzsy/small-kms/backend/models/key"
	"github.com/stephenzsy/small-kms/backend/resdoc"
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
	publicKey        crypto.PublicKey
	signer           crypto.Signer
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
	c context.Context,
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

	d.templateX509Cert = d.generateCertificateTemplate()
	if pDoc.IssuerPolicy.IsEmpty() {
		d.issuerX509Cert = d.templateX509Cert
		if pDoc.AllowGenerate {
			ckParams, err := d.getAzCreateKeyParams()
			if err != nil {
				return err
			}
			d.KeyVaultStore = &CertDocKeyVaultStore{
				Name: kv.GetMaterialName(kv.MaterialNameKindfCertificateKey, nsProvider, nsID, pDoc.ID),
			}
			azKeysClient := kv.GetAzKeyVaultService(c).AzKeysClient()
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
		}
	} else {
		// TODO: load issuer
		return fmt.Errorf("unimplemented")
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

func (d *certDocSelfSignedGeneratePending) collectSignedCert(cert []byte) (err error) {
	d.JsonWebKey.CertificateChain = []cloudkey.Base64RawURLEncodableBytes{cert}
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
	if d.JsonWebKey.KeyID != "" {
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

// var _ base.ModelPopulater[Certificate] = (*CertDoc)(nil)

// func (d *CertDoc) getSigningParams() kv.SigningParams {
// 	params := kv.SigningParams{
// 		SigAlg: azkeys.SignatureAlgorithm(*d.KeySpec.Alg),
// 	}
// 	if d.KeySpec.KeyID != nil {
// 		params.CertID = azcertificates.ID(*d.KeySpec.KeyID)
// 	}
// 	return params
// }

// func (d *CertDoc) getCSRProviderParams() kv.CSRProviderParams {
// 	params := kv.CSRProviderParams{
// 		CertName:      d.KeyVaultStore.Name,
// 		KeyProperties: azcertificates.KeyProperties{},
// 	}
// 	switch d.KeySpec.Kty {
// 	case cloudkey.KeyTypeRSA:
// 		params.KeyProperties.KeyType = to.Ptr(azcertificates.KeyTypeRSA)
// 		params.KeyProperties.KeySize = d.KeySpec.KeySize
// 	case cloudkey.KeyTypeEC:
// 		params.KeyProperties.KeyType = to.Ptr(azcertificates.KeyTypeEC)
// 		switch d.KeySpec.Crv {
// 		case cloudkey.CurveNameP256:
// 			params.KeyProperties.Curve = to.Ptr(azcertificates.CurveNameP256)
// 		case cloudkey.CurveNameP384:
// 			params.KeyProperties.Curve = to.Ptr(azcertificates.CurveNameP384)
// 		case cloudkey.CurveNameP521:
// 			params.KeyProperties.Curve = to.Ptr(azcertificates.CurveNameP521)
// 		}
// 	}
// 	params.KeyProperties.Exportable = &d.KeyExportable
// 	return params
// }

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

// func (d *CertDoc) applyPatch(c context.Context,
// 	docService base.AzCosmosCRUDDocService,
// 	patch *CertDocSigningPatch) error {
// 	nextKeySpec := d.KeySpec
// 	nextKeySpec.E = patch.KeySpec.E
// 	nextKeySpec.N = patch.KeySpec.N
// 	nextKeySpec.X = patch.KeySpec.X
// 	nextKeySpec.Y = patch.KeySpec.Y
// 	nextKeySpec.CertificateChain = patch.KeySpec.CertificateChain
// 	nextKeySpec.KeyID = patch.KeySpec.KeyID
// 	nextKeySpec.X5t = patch.KeySpec.X5t
// 	nextKeySpec.X5tS256 = patch.KeySpec.X5tS256

// 	patchOps := azcosmos.PatchOperations{}
// 	patchOps.AppendSet("/keySpec", nextKeySpec)
// 	patchOps.AppendSet("/keyVaultStore", patch.KeyVaultStore)
// 	patchOps.AppendSet("/issuer", patch.Issuer)
// 	patchOps.AppendSet("/status", CertificateStatusIssued)
// 	err := docService.Patch(c, d, patchOps, &azcosmos.ItemOptions{
// 		IfMatchEtag: d.ETag,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	d.KeySpec = nextKeySpec
// 	d.KeyVaultStore = patch.KeyVaultStore
// 	d.Issuer = patch.Issuer
// 	return nil
// }

// func (d *CertDoc) getIssuedX509Certificate() (*x509.Certificate, [][]byte, error) {
// 	chain := utils.MapSlice(d.KeySpec.CertificateChain, func(certBytes base.Base64RawURLEncodedBytes) []byte {
// 		return certBytes
// 	})
// 	cert, err := x509.ParseCertificate(chain[0])
// 	return cert, chain, err
// }

// func (d *CertDoc) getX509SignatureAlgorithm() (sa x509.SignatureAlgorithm) {
// 	switch *d.KeySpec.Alg {
// 	case cloudkey.SignatureAlgorithmRS256:
// 		sa = x509.SHA256WithRSA
// 	case cloudkey.SignatureAlgorithmRS384:
// 		sa = x509.SHA384WithRSA
// 	case cloudkey.SignatureAlgorithmRS512:
// 		sa = x509.SHA512WithRSA
// 	case cloudkey.SignatureAlgorithmPS256:
// 		sa = x509.SHA256WithRSAPSS
// 	case cloudkey.SignatureAlgorithmPS384:
// 		sa = x509.SHA384WithRSAPSS
// 	case cloudkey.SignatureAlgorithmPS512:
// 		sa = x509.SHA512WithRSAPSS
// 	case cloudkey.SignatureAlgorithmES256:
// 		sa = x509.ECDSAWithSHA256
// 	case cloudkey.SignatureAlgorithmES384:
// 		sa = x509.ECDSAWithSHA384
// 	case cloudkey.SignatureAlgorithmES512:
// 		sa = x509.ECDSAWithSHA512
// 	}
// 	return sa
// }

// func (d *CertDoc) x5cPEMBlocks() []pem.Block {
// 	return utils.MapSlice(d.KeySpec.CertificateChain, func(certBytes base.Base64RawURLEncodedBytes) pem.Block {
// 		return pem.Block{
// 			Type:  "CERTIFICATE",
// 			Bytes: certBytes,
// 		}
// 	})
// }
