package admin

import (
	"bytes"
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/stephenzsy/small-kms/backend/common"
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
	t.KeyStorePath = p.KeyStorePath
	if p.Usage == UsageClientOnly {
		t.KeyStorePath = ""
	} else if len(t.KeyStorePath) == 0 {
		t.KeyStorePath = strings.TrimSpace(p.KeyStorePath)
		if len(t.KeyStorePath) == 0 {
			return fmt.Errorf("missing KeyStorePath for usage: %s", p.Usage)
		}
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
	LastCertExpires time.Time               `json:"lastCertExpires"`
	LastAction      PolicyCertRequestAction `json:"lastAction"`
}

func (p *PolicyCertRequestDocSection) evaluateForAction(ctx context.Context, s *adminServer, namespaceID uuid.UUID, policyDoc *PolicyDoc, forceFlag *bool) (
	shouldTrigger bool, ps *PolicyStateDoc, msg string, err error) {
	shouldTrigger = false
	msg = "unknown"
	if forceFlag != nil && *forceFlag {
		shouldTrigger = true
		msg = "forced"
		return
	}
	// read policy state
	ps, err = s.GetPolicyStateDoc(ctx, namespaceID, policyDoc.GetUUID())
	if err != nil {
		if common.IsAzNotFound(err) {
			shouldTrigger = true
			msg = "no previous run"
			err = nil
			return
		} else {
			msg = "error reaing state"
			return
		}
	}
	if p.LifetimeTrigger == nil {
		msg = "no renewal configured"
		return
	}
	if p.LifetimeTrigger.DaysBeforeExpiry != nil {
		testExpireAfter := time.Now().AddDate(0, 0, int(*p.LifetimeTrigger.DaysBeforeExpiry))
		if ps.CertRequest.LastCertExpires.Before(testExpireAfter) {
			shouldTrigger = true
			msg = fmt.Sprintf("renew before %d days till expiry", *p.LifetimeTrigger.DaysBeforeExpiry)
			return
		}
	} else if p.LifetimeTrigger.LifetimePercentage != nil {
		testCutoff := ps.CertRequest.LastCertIssued.Add(ps.CertRequest.LastCertExpires.Sub(ps.CertRequest.LastCertIssued) *
			time.Duration(*p.LifetimeTrigger.LifetimePercentage) / 100)
		if testCutoff.Before(time.Now()) {
			shouldTrigger = true
			msg = fmt.Sprintf("renew after lifetime percentage %d%%", *p.LifetimeTrigger.LifetimePercentage)
			return
		}
	}
	msg = "no renewal needed"
	return
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
		r.KeySize = nil
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

func (p *KeyProperties) ToAzKeysCreateKeyParameters() (r azkeys.CreateKeyParameters) {
	r.Kty = to.Ptr(azkeys.KeyTypeRSA)
	r.KeySize = to.Ptr(int32(2048))
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
		r.Kty = to.Ptr(azkeys.KeyTypeEC)
		r.KeySize = nil
		r.Curve = to.Ptr(azkeys.CurveNameP256)
		if p.CurveName != nil {
			switch *p.CurveName {
			case EcCurveP384:
				r.Curve = to.Ptr(azkeys.CurveNameP384)
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
		r.DNSNames = to.SliceOfPtrs(*san.DNSNames...)
	}
	if san.Emails != nil && len(*san.Emails) > 0 {
		if r != nil {
			r = new(azcertificates.SubjectAlternativeNames)
		}
		r.Emails = to.SliceOfPtrs(*san.Emails...)
	}
	if san.UserPrincipalNames != nil && len(*san.UserPrincipalNames) > 0 {
		if r != nil {
			r = new(azcertificates.SubjectAlternativeNames)
		}
		r.UserPrincipalNames = to.SliceOfPtrs(*san.UserPrincipalNames...)
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

func prepareCertificate(p *PolicyCertRequestDocSection, namespaceID uuid.UUID, certId uuid.UUID, csr *x509.CertificateRequest) (c x509.Certificate, err error) {
	// use certificate ID
	serialNumber := big.NewInt(0)
	serialNumber = serialNumber.SetBytes(certId[:])
	c.SerialNumber = serialNumber
	c.NotBefore = time.Now()
	c.NotAfter = time.Now().AddDate(0, int(p.ValidityInMonths), 0)

	if csr != nil {
		c.Subject = csr.Subject
		if !IsCANamespace(namespaceID) {
			c.Extensions = csr.Extensions
		}
		c.EmailAddresses = csr.EmailAddresses
		c.DNSNames = csr.DNSNames
		c.IPAddresses = csr.IPAddresses
		c.URIs = csr.URIs
	} else {
		c.Subject = p.Subject.ToPkixName()
	}

	if IsRootCANamespace(namespaceID) {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLen = 1
		c.MaxPathLenZero = false
		c.BasicConstraintsValid = true
	} else if IsIntCANamespace(namespaceID) {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLenZero = true
		c.BasicConstraintsValid = true
	} else {
		err = fmt.Errorf("namespace not supported yet: %s", namespaceID.String())
		return
	}

	return
}

func (s *adminServer) getRootCASigner(ctx context.Context, keyStorePath string, expires time.Time, p *PolicyCertRequestDocSection) (crypto.Signer, string, error) {
	var signerKey *azkeys.JSONWebKey = nil
	reuseKey := false
	if p.KeyProperties.ReuseKey != nil {
		reuseKey = *p.KeyProperties.ReuseKey
	}
	if reuseKey {
		getKeyResp, err := s.AzKeysClient().GetKey(ctx, keyStorePath, "", nil)
		if err != nil {
			if !common.IsAzNotFound(err) {
				return nil, "", err
			}
		} else if getKeyResp.Attributes.Expires != nil && getKeyResp.Attributes.Expires.After(expires) {
			signerKey = getKeyResp.Key
		}
		log.Info().Msgf("Using existing existing key: %s", *signerKey.KID)
	}

	if signerKey == nil {
		// creating key
		params := p.KeyProperties.ToAzKeysCreateKeyParameters()
		params.KeyOps = to.SliceOfPtrs(azkeys.KeyOperationSign, azkeys.KeyOperationVerify)
		if !reuseKey {
			params.KeyAttributes = &azkeys.KeyAttributes{
				Expires: &expires,
			}
		}
		resp, err := s.AzKeysClient().CreateKey(ctx, keyStorePath, params, nil)
		if err != nil {
			return nil, "", err
		}
		signerKey = resp.Key
	}

	signer, err := newKeyVaultSigner(ctx, s.AzKeysClient(), signerKey)
	return signer, string(*signerKey.KID), err
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

	isRootNS := IsRootCANamespace(namespaceID)

	// prepare certificate
	var csr *x509.CertificateRequest
	var certPubKey crypto.PublicKey
	var kid string
	if !isRootNS {
		azCreateCertParams := p.ToKeyvaultCreateCertificateParameters(namespaceID)
		azcCcResp, err := s.AzCertificatesClient().CreateCertificate(ctx, keyName, azCreateCertParams, nil)
		if err != nil {
			return nil, err
		}
		log.Info().Msgf("Created certificate in KeyVault: %s", *azcCcResp.ID)
		csr, err = x509.ParseCertificateRequest(azcCcResp.CSR)
		if err != nil {
			return nil, err
		}
		certPubKey = csr.PublicKey
	}
	certificate, err := prepareCertificate(p, namespaceID, certID, csr)
	if err != nil {
		return
	}

	// get signer, and signer certificate
	var signerCert *x509.Certificate
	var signer crypto.Signer
	var signerPemBytes []byte
	if isRootNS {
		signer, kid, err = s.getRootCASigner(ctx, keyName, certificate.NotAfter, p)
		if err != nil {
			return nil, err
		}
		certPubKey = signer.Public()
		signerCert = &certificate
	} else if IsIntCANamespace(namespaceID) {
		// load certificate
		crtDoc, err := s.GetLatestCertDocForPolicy(ctx, policyID, policyID)
		if err != nil {
			return nil, err
		}
		if len(crtDoc.KID) == 0 {
			return nil, errors.New("issuer certificate key not found")
		}
		keyId := azkeys.ID(crtDoc.KID)
		resp, err := s.AzKeysClient().GetKey(ctx, keyId.Name(), keyId.Version(), nil)
		if err != nil {
			return nil, err
		}
		signer, err = newKeyVaultSigner(ctx, s.AzKeysClient(), resp.Key)
		if err != nil {
			return nil, err
		}
		signerBlobName := crtDoc.CertStorePath
		if len(signerBlobName) == 0 {
			return nil, errors.New("issuer certificate missing")
		}
		signerPemBytes, err := s.FetchCertificatePEMBlob(ctx, signerBlobName)
		if err != nil {
			return nil, err
		}
		pemBlock, _ := pem.Decode(signerPemBytes)
		signerCert, err = x509.ParseCertificate(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}
	} else {
		// TODO
		return nil, fmt.Errorf("namespace not supported yet: %s", namespaceID.String())
	}

	switch p.KeyProperties.KeyType {
	case KtyRSA:
		certificate.SignatureAlgorithm = x509.SHA384WithRSA
	case KtyEC:
		switch *p.KeyProperties.CurveName {
		case EcCurveP384:
			certificate.SignatureAlgorithm = x509.ECDSAWithSHA384
		case EcCurveP256:
			certificate.SignatureAlgorithm = x509.ECDSAWithSHA256
		default:
			return nil, fmt.Errorf("unsupported curve: %s", *p.KeyProperties.CurveName)
		}
	default:
		return nil, fmt.Errorf("unsupported key type: %s", p.KeyProperties.KeyType)
	}

	// Sign cert
	certSigned, err := x509.CreateCertificate(nil, &certificate, signerCert, certPubKey, signer)
	if err != nil {
		return
	}

	log.Info().Msgf("Certificate signed and validated, prepare to upload")
	// encode to pem
	pemBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certSigned,
	}
	bb := bytes.Buffer{}
	err = pem.Encode(&bb, &pemBlock)
	if err != nil {
		return
	}
	if !isRootNS {
		// attach chain
		_, err = bb.Write(signerPemBytes)
		if err != nil {
			return nil, err
		}
	}
	// upload to certificate storage only
	blobName := fmt.Sprintf("%s/%s.pem", namespaceID, certID)
	_, err = s.azBlobContainerClient.NewBlockBlobClient(blobName).UploadBuffer(ctx, bb.Bytes(), &blockblob.UploadBufferOptions{
		HTTPHeaders: &blob.HTTPHeaders{
			BlobContentType: to.Ptr("application/x-pem-file"),
		},
	})
	if err != nil {
		return
	}
	log.Info().Msgf("Certificate uploaded to blob")

	var cid string
	if !isRootNS && IsIntCANamespace(namespaceID) && len(keyName) > 0 {
		azcMcResp, err := s.AzCertificatesClient().MergeCertificate(ctx, keyName, azcertificates.MergeCertificateParameters{
			X509Certificates: [][]byte{certSigned, signerCert.Raw},
		}, nil)
		if err != nil {
			return nil, err
		}
		cid = string(*azcMcResp.ID)
		kid = string(*azcMcResp.KID)
		log.Info().Msgf("Certificate merged to keyvault: %s", cid)
	}

	// record certdoc
	certDoc := CertDoc{
		BaseDoc: kmsdoc.BaseDoc{
			ID:          kmsdoc.NewKmsDocID(kmsdoc.DocTypeCert, certID),
			NamespaceID: namespaceID,
		},
		PolicyID: policyID,
		Expires:  certificate.NotAfter,
		Usage:    p.Usage,
		CID:      cid,
		KID:      kid,
		//		SID:           string(*azcMcResp.SID),
		CertStorePath: blobName,
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
