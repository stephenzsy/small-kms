package admin

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyCertRequestDocSection struct {
	IssuerNamespaceID       uuid.UUID                           `json:"issuerNamespaceId"`
	KeyProperties           KeyProperties                       `json:"keyProperties"`
	KeyStorePath            string                              `json:"keyStorePath"`
	Subject                 CertificateSubject                  `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
	Usage                   CertificateUsage                    `json:"usage"`
	ValidityInMonths        int32                               `json:"validity_months"`
	LifetimeTrigger         *CertificateLifetimeTrigger         `json:"lifetimeTrigger,omitempty"`
}

const maxValidityInMonths = 120
const defaultValidityInMonths = 12

func getDefaultKeyProperties(namespaceID uuid.UUID) (kp KeyProperties) {
	kp.KeyType = KtyRSA
	kp.KeySize = ToPtr(KeySize2048)
	if IsCANamespace(namespaceID) {
		kp.KeySize = ToPtr(KeySize4096)
	}
	if IsTestCA(namespaceID) {
		kp.KeyType = KtyEC
		kp.KeySize = nil
		kp.CurveName = ToPtr(EcCurveP384)
	}
	return
}

func (t *PolicyCertRequestDocSection) validateAndFillWithParameters(p *CertificateRequestPolicyParameters, namespaceID uuid.UUID) error {
	if p == nil {
		return errors.New("missing CertRequest property")
	}
	// validate issuer namespace
	if !IsCANamespace(p.IssuerNamespaceID) {
		return fmt.Errorf("requested issuer namespace is not a CA namespace: %s", p.IssuerNamespaceID)
	}

	// validate usage
	switch p.Usage {
	case UsageRootCA:
		if !IsRootCANamespace(namespaceID) {
			return fmt.Errorf("usage %s is not valid for namespace: %s", p.Usage, namespaceID)
		}
	case UsageIntCA:
		if !IsIntCANamespace(namespaceID) {
			return fmt.Errorf("usage %s is not valid for namespace: %s", p.Usage, namespaceID)
		}
	case UsageServerAndClient,
		UsageClientOnly,
		UsageServerOnly:
		if IsCANamespace(namespaceID) {
			return fmt.Errorf("usage %s is not valid for namespace: %s", p.Usage, namespaceID)
		}
	default:
		return fmt.Errorf("invalid usage: %s", p.Usage)
	}
	t.Usage = p.Usage

	// key store path to store certificate in keyvault is required except for client only certificates
	if p.Usage == UsageClientOnly {
		t.KeyStorePath = ""
	} else if len(t.KeyStorePath) == 0 {
		t.KeyStorePath = strings.TrimSpace(p.KeyStorePath)
		return fmt.Errorf("missing KeyStorePath for usage: %s", p.Usage)
	}

	t.Subject = p.Subject
	if t.Subject.CN = strings.TrimSpace(t.Subject.CN); len(t.Subject.CN) == 0 {
		return fmt.Errorf("missing CN in subject")
	}
	t.SubjectAlternativeNames = p.SubjectAlternativeNames

	// validate life time trigger, invalid configurations are dropped
	if p.LifetimeTrigger != nil {
		if p.LifetimeTrigger.DaysBeforeExpiry != nil && *p.LifetimeTrigger.DaysBeforeExpiry > 0 {
			t.LifetimeTrigger = &CertificateLifetimeTrigger{
				DaysBeforeExpiry: p.LifetimeTrigger.DaysBeforeExpiry,
			}
		} else if p.LifetimeTrigger.LifetimePercentage != nil &&
			*p.LifetimeTrigger.LifetimePercentage > 0 &&
			*p.LifetimeTrigger.LifetimePercentage < 100 {
			t.LifetimeTrigger = &CertificateLifetimeTrigger{
				LifetimePercentage: p.LifetimeTrigger.LifetimePercentage,
			}
		}
	}

	// validate validity in months, max 10 years, default 12 months, minimum 1
	if p.ValidityInMonths == nil || *p.ValidityInMonths == 0 {
		t.ValidityInMonths = defaultValidityInMonths
	} else if *p.ValidityInMonths > maxValidityInMonths {
		t.ValidityInMonths = maxValidityInMonths
	} else if *p.ValidityInMonths < 1 {
		t.ValidityInMonths = 1
	} else {
		t.ValidityInMonths = *p.ValidityInMonths
	}

	t.KeyProperties = getDefaultKeyProperties(namespaceID)
	// keyspec
	if p.KeyProperties != nil {
		switch p.KeyProperties.KeyType {
		case KtyRSA:
			if t.KeyProperties.KeyType != KtyRSA {
				t.KeyProperties.KeyType = KtyRSA
				t.KeyProperties.CurveName = nil
				t.KeyProperties.KeySize = ToPtr(KeySize2048)
			}
			if t.KeyProperties.KeySize != nil {
				switch *p.KeyProperties.KeySize {
				case KeySize2048,
					KeySize3072,
					KeySize4096:
					t.KeyProperties.KeySize = p.KeyProperties.KeySize
				}
			}
		case KtyEC:
			if t.KeyProperties.KeyType != KtyEC {
				t.KeyProperties.KeyType = KtyEC
				t.KeyProperties.CurveName = ToPtr(EcCurveP256)
				t.KeyProperties.KeySize = nil
			}
			if t.KeyProperties.CurveName != nil {
				switch *p.KeyProperties.CurveName {
				case EcCurveP256,
					EcCurveP384:
					t.KeyProperties.CurveName = p.KeyProperties.CurveName
				}
			}
		}
	}

	return nil
}

type PolicyCertRequestAction string

const (
	PolicyCertRequestActionIssue PolicyCertRequestAction = "issue-certificate"
)

type PolicyStateCertRequestDocSection struct {
	LastCertCUID    kmsdoc.KmsDocID         `json:"lastCertId"`
	LastCertIssued  time.Time               `json:"lastCertIssued"`
	LastCertExpires time.Time               `json:"lastCertExpired"`
	LastAction      PolicyCertRequestAction `json:"lastAction"`
}

func (p *PolicyCertRequestDocSection) evaluateForAction(ctx context.Context, s *adminServer, namespaceID uuid.UUID, policyDoc *PolicyDoc, forceFlag *bool) (bool, string, error) {
	if forceFlag != nil && *forceFlag {
		return true, "forced", nil
	}
	// read policy state
	ps, err := s.GetPolicyStateDoc(ctx, namespaceID, policyDoc.GetUUID())
	if err != nil {
		if kmsdoc.IsNotFound(err) {
			return true, "initial policy run", nil
		}
		return true, "error reading state", err
	}
	if p.LifetimeTrigger == nil {
		return false, "no renewal configured", nil
	}
	if p.LifetimeTrigger.DaysBeforeExpiry != nil {
		testExpireAfter := time.Now().AddDate(0, 0, int(*p.LifetimeTrigger.DaysBeforeExpiry))
		if ps.CertRequest.LastCertExpires.Before(testExpireAfter) {
			return true, fmt.Sprintf("renew before %d days till expiry", *p.LifetimeTrigger.DaysBeforeExpiry), nil
		}
	} else if p.LifetimeTrigger.LifetimePercentage != nil {
		testCutoff := ps.CertRequest.LastCertIssued.Add(ps.CertRequest.LastCertExpires.Sub(ps.CertRequest.LastCertIssued) *
			time.Duration(*p.LifetimeTrigger.LifetimePercentage) / 100)
		if testCutoff.Before(time.Now()) {
			return true, fmt.Sprintf("renew after lifetime percentage %d%%", *p.LifetimeTrigger.LifetimePercentage), nil
		}
	}
	return false, "no renewal needed", nil
}

func (p *KeyProperties) ToAzCertificatesKeyProperties() (r azcertificates.KeyProperties) {
	r.KeyType = ToPtr(azcertificates.KeyTypeRSA)
	r.KeySize = ToPtr(int32(2048))
	r.ReuseKey = p.ReuseKey
	switch p.KeyType {
	case KtyRSA:
		if p.KeySize != nil {
			switch *p.KeySize {
			case KeySize3072:
				r.KeySize = ToPtr(int32(3072))
			case KeySize4096:
				r.KeySize = ToPtr(int32(4096))
			}
		}
	case KtyEC:
		r.KeyType = ToPtr(azcertificates.KeyTypeEC)
		r.Curve = ToPtr(azcertificates.CurveNameP256)
		if p.CurveName != nil {
			switch *p.CurveName {
			case EcCurveP384:
				r.Curve = ToPtr(azcertificates.CurveNameP384)
			}
		}
	}
	return

}

func (san *CertificateSubjectAlternativeNames) ToAzCertificatesSubjectAlternativeNames() (r *azcertificates.SubjectAlternativeNames) {
	if san == nil {
		return nil
	}
	if san.DNSNames != nil && len(*san.DNSNames) > 0 {
		if r != nil {
			r = new(azcertificates.SubjectAlternativeNames)
		}
		r.DNSNames = make([]*string, len(*san.DNSNames))
		for i, v := range *san.DNSNames {
			r.DNSNames[i] = ToPtr(v)
		}
	}
	if san.Emails != nil && len(*san.Emails) > 0 {
		if r != nil {
			r = new(azcertificates.SubjectAlternativeNames)
		}
		r.Emails = make([]*string, len(*san.Emails))
		for i, v := range *san.Emails {
			r.Emails[i] = ToPtr(v)
		}
	}
	if san.UserPrincipalNames != nil && len(*san.UserPrincipalNames) > 0 {
		if r != nil {
			r = new(azcertificates.SubjectAlternativeNames)
		}
		r.UserPrincipalNames = make([]*string, len(*san.UserPrincipalNames))
		for i, v := range *san.UserPrincipalNames {
			r.UserPrincipalNames[i] = ToPtr(v)
		}
	}
	return r
}

func (p *PolicyCertRequestDocSection) ToKeyvaultCreateCertificateParameters(namespaceID uuid.UUID) (r azcertificates.CreateCertificateParameters) {

	x509Properties := azcertificates.X509CertificateProperties{
		Subject:                 ToPtr(p.Subject.ToPkixName().String()),
		ValidityInMonths:        ToPtr(int32(p.ValidityInMonths)),
		SubjectAlternativeNames: p.SubjectAlternativeNames.ToAzCertificatesSubjectAlternativeNames(),
		EnhancedKeyUsage:        make([]*string, 0, 2),
	}
	if p.Usage == UsageServerAndClient || p.Usage == UsageServerOnly {
		x509Properties.EnhancedKeyUsage = append(x509Properties.EnhancedKeyUsage, to.Ptr("1.3.6.1.5.5.7.3.1"))
	}
	if p.Usage == UsageServerAndClient || p.Usage == UsageClientOnly {
		x509Properties.EnhancedKeyUsage = append(x509Properties.EnhancedKeyUsage, to.Ptr("1.3.6.1.5.5.7.3.2"))
	}

	keyProperties := p.KeyProperties.ToAzCertificatesKeyProperties()
	if p.Usage == UsageRootCA || p.Usage == UsageIntCA {
		keyProperties.Exportable = to.Ptr(false)
		x509Properties.KeyUsage = []*azcertificates.KeyUsageType{
			ToPtr(azcertificates.KeyUsageTypeDigitalSignature),
			ToPtr(azcertificates.KeyUsageTypeKeyCertSign),
			ToPtr(azcertificates.KeyUsageTypeCRLSign),
		}
	}

	r.CertificatePolicy = &azcertificates.CertificatePolicy{
		Attributes: &azcertificates.CertificateAttributes{
			Enabled: to.Ptr(true),
		},
		KeyProperties:             &keyProperties,
		X509CertificateProperties: &x509Properties,
		SecretProperties: &azcertificates.SecretProperties{
			ContentType: to.Ptr("application/x-pem-file"),
		},
	}

	return
}

func (p *PolicyCertRequestDocSection) Fillx509(c *x509.Certificate, certId uuid.UUID, namespaceID uuid.UUID, csr *x509.CertificateRequest) error {
	// use certificate ID
	serialNumber := big.NewInt(0)
	serialNumber = serialNumber.SetBytes(certId[:])
	c.SerialNumber = serialNumber
	c.Subject = csr.Subject
	c.NotBefore = time.Now()
	c.NotAfter = time.Now().AddDate(0, int(p.ValidityInMonths), 0)
	c.Extensions = csr.Extensions
	c.EmailAddresses = csr.EmailAddresses
	c.DNSNames = csr.DNSNames
	c.IPAddresses = csr.IPAddresses
	c.URIs = csr.URIs

	if IsRootCANamespace(namespaceID) {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLen = 1
		c.MaxPathLenZero = false
		c.BasicConstraintsValid = true
	} else {
		return fmt.Errorf("namespace not supported yet: %s", namespaceID.String())
	}

	switch p.KeyProperties.KeyType {
	case KtyRSA:
		c.SignatureAlgorithm = x509.SHA384WithRSA
	case KtyEC:
		c.SignatureAlgorithm = x509.ECDSAWithSHA384
	default:
		return fmt.Errorf("unsupported key type: %s", p.KeyProperties.KeyType)
	}
	return nil
}

func toPublicKey(key *azkeys.JSONWebKey) (crypto.PublicKey, error) {
	if *key.Kty == azkeys.KeyTypeRSA {
		k := &rsa.PublicKey{}

		// N = modulus
		if len(key.N) == 0 {
			return nil, errors.New("property N is empty")
		}
		k.N = &big.Int{}
		k.N.SetBytes(key.N)

		// e = public exponent
		if len(key.E) == 0 {
			return nil, errors.New("property e is empty")
		}
		k.E = int(big.NewInt(0).SetBytes(key.E).Uint64())
		return k, nil
	} else if *key.Kty == azkeys.KeyTypeEC {
		k := &ecdsa.PublicKey{}

		switch *key.Crv {
		case azkeys.CurveNameP256:
			k.Curve = elliptic.P256()
		case azkeys.CurveNameP384:
			k.Curve = elliptic.P384()
		default:
			return nil, fmt.Errorf("unsupported curve: %s", *key.Crv)
		}

		if len(key.X) == 0 {
			return nil, errors.New("property X is empty")
		}
		k.X = &big.Int{}
		k.X.SetBytes(key.X)

		if len(key.Y) == 0 {
			return nil, errors.New("property Y is empty")
		}
		k.Y = &big.Int{}
		k.Y.SetBytes(key.Y)

		return k, nil
	}
	return nil, fmt.Errorf("unsupported key type: %s", *key.Kty)
}

func (p *PolicyCertRequestDocSection) action(ctx *gin.Context, s *adminServer, namespaceID uuid.UUID, policyDoc *PolicyDoc) (resultDoc *PolicyStateDoc, err error) {

	policyID := policyDoc.GetUUID()

	// request for new certifiate
	// read policy state
	log.Info().Msgf("Start CertRequest action for policy %s/%s", namespaceID, policyID)
	certID := uuid.New()
	log.Info().Msgf("Certificate ID %s", certID)
	keyName := p.KeyStorePath
	log.Info().Msgf("Certificate keyStorePath: %s", keyName)

	azCreateCertParams := p.ToKeyvaultCreateCertificateParameters(namespaceID)
	azcCcResp, err := s.AzCertificatesClient().CreateCertificate(ctx, keyName, azCreateCertParams, nil)
	if err != nil {
		return
	}
	log.Info().Msgf("Created certificate in KeyVault: %s", *azcCcResp.ID)
	csr, err := x509.ParseCertificateRequest(azcCcResp.CSR)
	if err != nil {
		return
	}

	certificate := x509.Certificate{}
	p.Fillx509(&certificate, certID, namespaceID, csr)

	// TODO non root CA signers
	signerKeyResp, err := s.AzKeysClient().GetKey(ctx, keyName, "", nil)
	if err != nil {
		return
	}
	signerKey := signerKeyResp.Key
	log.Info().Msgf("Using signer key: %s", *signerKey.KID)
	signerPubkey, err := toPublicKey(signerKey)
	if err != nil {
		return
	}
	signer := keyVaultSigner{
		ctx:        ctx,
		keysClient: s.AzKeysClient(),
		kid:        signerKey.KID,
		publicKey:  signerPubkey,
	}

	// Sign cert
	certSigned, err := x509.CreateCertificate(nil, &certificate, &certificate, csr.PublicKey, &signer)
	if err != nil {
		return
	}
	azcMcResp, err := s.AzCertificatesClient().MergeCertificate(ctx, keyName, azcertificates.MergeCertificateParameters{
		X509Certificates: [][]byte{certSigned},
	}, nil)
	if err != nil {
		return
	}
	log.Info().Msgf("Certificate signed by: %s, merged to keyvault: %s", *signerKey.KID, *azcMcResp.ID)

	// record certdoc
	certDoc := CertDoc{
		BaseDoc: kmsdoc.BaseDoc{
			ID:          kmsdoc.NewKmsDocID(kmsdoc.DocTypeCert, certID),
			NamespaceID: namespaceID,
		},
		PolicyID: policyID,
		Expires:  certificate.NotAfter,
		Usage:    p.Usage,
		CID:      string(*azcMcResp.ID),
		KID:      string(*azcMcResp.KID),
		SID:      string(*azcMcResp.SID),
	}
	err = kmsdoc.AzCosmosUpsert(ctx, s.azCosmosContainerClientCerts, &certDoc)
	if err != nil {
		return
	}
	log.Info().Msgf("Created certificate: %s/%s", namespaceID, certDoc.ID.String())

	// record latest cert doc
	certDocL := certDoc
	certDocL.ID = kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, policyID)
	certDocL.AliasID = &certDoc.ID
	err = kmsdoc.AzCosmosUpsert(ctx, s.azCosmosContainerClientCerts, &certDocL)
	if err != nil {
		return
	}
	log.Info().Msgf("Set certificate as latest for policy(%s): %s/%s", policyID, namespaceID, certDoc.ID.String())

	// record policy state
	resultDoc = &PolicyStateDoc{
		BaseDoc: kmsdoc.BaseDoc{
			ID:          kmsdoc.NewKmsDocID(kmsdoc.DocTypePolicyState, policyID),
			NamespaceID: namespaceID,
		},
		PolicyType: PolicyTypeCertRequest,
		Status:     PolicyStateStatusSuccess,
		Message:    fmt.Sprintf("certificate issued: %s/%s[%s]", namespaceID, certID, policyID),
		CertRequest: &PolicyStateCertRequestDocSection{
			LastCertCUID:    certDoc.ID,
			LastCertIssued:  certificate.NotBefore,
			LastCertExpires: certificate.NotAfter,
			LastAction:      PolicyCertRequestActionIssue,
		},
	}
	err = kmsdoc.AzCosmosUpsert(ctx, s.azCosmosContainerClientCerts, resultDoc)
	if err != nil {
		return
	}
	log.Info().Msgf("CertRequest completed for %s/%s", namespaceID, policyID)
	return
}

func (s *PolicyCertRequestDocSection) ToCertificateRequestPolicyParameters() *CertificateRequestPolicyParameters {
	if s == nil {
		return nil
	}
	return &CertificateRequestPolicyParameters{
		IssuerNamespaceID:       s.IssuerNamespaceID,
		KeyProperties:           &s.KeyProperties,
		KeyStorePath:            s.KeyStorePath,
		LifetimeTrigger:         s.LifetimeTrigger,
		Subject:                 s.Subject,
		SubjectAlternativeNames: s.SubjectAlternativeNames,
		Usage:                   s.Usage,
		ValidityInMonths:        ToPtr(s.ValidityInMonths),
	}
}

func (s *PolicyStateCertRequestDocSection) ToPolicyStateCertRequest() *PolicyStateCertRequest {
	if s == nil {
		return nil
	}
	return &PolicyStateCertRequest{
		LastCertID:      s.LastCertCUID.GetUUID(),
		LastCertIssued:  s.LastCertIssued,
		LastCertExpires: s.LastCertExpires,
		LastAction:      string(s.LastAction),
	}
}
