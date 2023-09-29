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
	PolicyDocSectionIssuerProperties
	KeyProperties           KeyProperties                       `json:"keyProperties"`
	KeyStorePath            string                              `json:"keyStorePath"`
	Subject                 CertificateSubject                  `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"subjectAlternativeNames,omitempty"`
	Usage                   CertificateUsage                    `json:"usage"`
	ValidityInMonths        int32                               `json:"validity_months"`
	LifetimeTrigger         *CertificateLifetimeTrigger         `json:"lifetimeTrigger,omitempty"`
}

func getDefaultKeyProperties(namespaceID uuid.UUID) (kp KeyProperties) {
	kp.Kty = KeyTypeRSA
	kp.KeySize = ToPtr(KeySize2048)
	if IsCANamespace(namespaceID) {
		kp.KeySize = ToPtr(KeySize4096)
	}
	if IsTestCA(namespaceID) {
		kp.Kty = KeyTypeEC
		kp.KeySize = nil
		kp.Crv = ToPtr(CurveNameP384)
	}
	return
}

func (t *PolicyCertRequestDocSection) validateAndFillWithParameters(p *CertificateRequestPolicyParameters, namespaceID uuid.UUID, dirProfile *DirectoryObjectDoc) error {
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

	if err := t.PolicyDocSectionIssuerProperties.validateAndFillWithCertRequestParameters(p); err != nil {
		return err
	}

	// key store path to store certificate in keyvault is required except for client only certificates
	if p.Usage == UsageClientOnly {
	} else if len(t.KeyStorePath) == 0 {
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

	t.KeyProperties = getDefaultKeyProperties(namespaceID)
	// keyspec
	if p.KeyProperties != nil {
		switch p.KeyProperties.Kty {
		case KeyTypeRSA:
			if t.KeyProperties.Kty != KeyTypeRSA {
				t.KeyProperties.Kty = KeyTypeRSA
				t.KeyProperties.Crv = nil
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
		case KeyTypeEC:
			if t.KeyProperties.Kty != KeyTypeEC {
				t.KeyProperties.Kty = KeyTypeEC
				t.KeyProperties.Crv = ToPtr(CurveNameP256)
				t.KeyProperties.KeySize = nil
			}
			if t.KeyProperties.Crv != nil {
				switch *p.KeyProperties.Crv {
				case CurveNameP256,
					CurveNameP384:
					t.KeyProperties.Crv = p.KeyProperties.Crv
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
	switch p.Kty {
	case KeyTypeRSA:
		if p.KeySize != nil {
			switch *p.KeySize {
			case KeySize3072:
				r.KeySize = ToPtr(int32(3072))
			case KeySize4096:
				r.KeySize = ToPtr(int32(4096))
			}
		}
	case KeyTypeEC:
		r.KeyType = ToPtr(azcertificates.KeyTypeEC)
		r.KeySize = nil
		r.Curve = ToPtr(azcertificates.CurveNameP256)
		if p.Crv != nil {
			switch *p.Crv {
			case CurveNameP384:
				r.Curve = ToPtr(azcertificates.CurveNameP384)
			}
		}
	}
	return

}

func (p *KeyProperties) ToAzKeysCreateKeyParameters() (r azkeys.CreateKeyParameters) {
	r.Kty = to.Ptr(azkeys.KeyTypeRSA)
	r.KeySize = to.Ptr(int32(2048))
	switch p.Kty {
	case KeyTypeRSA:
		if p.KeySize != nil {
			switch *p.KeySize {
			case KeySize3072:
				r.KeySize = ToPtr(int32(3072))
			case KeySize4096:
				r.KeySize = ToPtr(int32(4096))
			}
		}
	case KeyTypeEC:
		r.Kty = to.Ptr(azkeys.KeyTypeEC)
		r.KeySize = nil
		r.Curve = to.Ptr(azkeys.CurveNameP256)
		if p.Crv != nil {
			switch *p.Crv {
			case CurveNameP384:
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
	if len(san.DNSNames) > 0 {
		if r != nil {
			r = new(azcertificates.SubjectAlternativeNames)
		}
		r.DNSNames = to.SliceOfPtrs(san.DNSNames...)
	}
	if len(san.EmailAddresses) > 0 {
		if r != nil {
			r = new(azcertificates.SubjectAlternativeNames)
		}
		r.Emails = to.SliceOfPtrs(san.EmailAddresses...)
	}
	return r
}

func (p *PolicyCertRequestDocSection) ToKeyvaultCreateCertificateParameters(namespaceID uuid.UUID) (r azcertificates.CreateCertificateParameters) {

	x509Properties := azcertificates.X509CertificateProperties{
		ValidityInMonths:        ToPtr(int32(p.ValidityInMonths)),
		SubjectAlternativeNames: p.SubjectAlternativeNames.ToAzCertificatesSubjectAlternativeNames(),
	}

	keyProperties := p.KeyProperties.ToAzCertificatesKeyProperties()
	if p.Usage == UsageRootCA || p.Usage == UsageIntCA {
		keyProperties.Exportable = to.Ptr(false)
	} else {
		keyProperties.Exportable = to.Ptr(true)
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
		c.EmailAddresses = csr.EmailAddresses
		c.DNSNames = csr.DNSNames
		c.IPAddresses = csr.IPAddresses
		c.URIs = csr.URIs
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
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment
		if p.Usage == UsageServerAndClient || p.Usage == UsageServerOnly {
			c.ExtKeyUsage = append(c.ExtKeyUsage, x509.ExtKeyUsageServerAuth)
		}
		if p.Usage == UsageServerAndClient || p.Usage == UsageClientOnly {
			c.ExtKeyUsage = append(c.ExtKeyUsage, x509.ExtKeyUsageClientAuth)
		}
	}

	return
}

func (s *adminServer) getRootCASigner(ctx context.Context, keyStorePath string, expires time.Time, p *PolicyCertRequestDocSection) (*keyVaultSigner, string, error) {
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

type signerCertBundle struct {
	privateKey                  *keyVaultSigner
	certificate                 *x509.Certificate
	certificateChainPEMRaw      []byte
	additionalCertificateDERRaw []byte
	signerNamespaceID           uuid.UUID
	signerCertId                kmsdoc.KmsDocID
}

func (b *signerCertBundle) SignatureAlgorithm() x509.SignatureAlgorithm {
	switch b.privateKey.sigAlg {
	case azkeys.SignatureAlgorithmRS384:
		return x509.SHA384WithRSA
	case azkeys.SignatureAlgorithmES256:
		return x509.ECDSAWithSHA256
	case azkeys.SignatureAlgorithmES384:
		return x509.ECDSAWithSHA384
	}
	return x509.UnknownSignatureAlgorithm
}

func (s *adminServer) loadSignerCertificateBundle(ctx context.Context, signerNamespaceID uuid.UUID, signerPolicyID uuid.UUID) (*signerCertBundle, error) {
	// load certificate
	crtDoc, err := s.getCertDoc(ctx, signerNamespaceID, kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, signerPolicyID))
	if err != nil {
		return nil, err
	}

	resp, err := s.AzKeysClient().GetKey(ctx, "", "", nil)
	if err != nil {
		return nil, err
	}
	bundle := new(signerCertBundle)

	// signer info
	bundle.signerNamespaceID = signerNamespaceID
	if crtDoc.AliasID != nil {
		// doc is alias
		bundle.signerCertId = *crtDoc.AliasID
	} else {
		bundle.signerCertId = crtDoc.ID
	}
	bundle.privateKey, err = newKeyVaultSigner(ctx, s.AzKeysClient(), resp.Key)
	if err != nil {
		return nil, err
	}
	signerBlobName := crtDoc.CertStorePath
	if len(signerBlobName) == 0 {
		return nil, errors.New("issuer certificate missing")
	}
	bundle.certificateChainPEMRaw, err = s.FetchCertificatePEMBlob(ctx, signerBlobName)
	if err != nil {
		return nil, err
	}
	pemBlock, rest := pem.Decode(bundle.certificateChainPEMRaw)
	bundle.certificate, err = x509.ParseCertificate(pemBlock.Bytes)
	extraPemBlock, _ := pem.Decode(rest)
	if extraPemBlock != nil {
		bundle.additionalCertificateDERRaw = extraPemBlock.Bytes
	}
	return bundle, err
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
	var signerBundle *signerCertBundle
	if isRootNS {
		signerBundle = new(signerCertBundle)
		signerBundle.certificate = &certificate
		signerBundle.privateKey, _, err = s.getRootCASigner(ctx, keyName, certificate.NotAfter, p)
		if err != nil {
			return nil, err
		}
		certPubKey = signerBundle.privateKey.Public()
	} else {
		signerBundle, err = s.loadSignerCertificateBundle(ctx, policyDoc.CertRequest.IssuerNamespaceID, policyDoc.CertRequest.IssuerPolicyID)
		if err != nil {
			return nil, err
		}
	}
	certificate.SignatureAlgorithm = signerBundle.SignatureAlgorithm()

	// Sign cert
	certSigned, err := x509.CreateCertificate(nil, &certificate, signerBundle.certificate, certPubKey, signerBundle.privateKey)
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
		_, err = bb.Write(signerBundle.certificateChainPEMRaw)
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
	if !isRootNS && len(keyName) > 0 {
		mergeRequestCerts := [][]byte{certSigned, signerBundle.certificate.Raw}
		if len(signerBundle.additionalCertificateDERRaw) > 0 {
			mergeRequestCerts = append(mergeRequestCerts, signerBundle.additionalCertificateDERRaw)
		}
		azcMcResp, err := s.AzCertificatesClient().MergeCertificate(ctx, keyName, azcertificates.MergeCertificateParameters{
			X509Certificates: mergeRequestCerts,
		}, nil)
		if err != nil {
			return nil, err
		}
		cid = string(*azcMcResp.ID)
		log.Info().Msgf("Certificate merged to keyvault: %s", cid)
	}

	// record certdoc
	certDoc := CertDoc{
		BaseDoc: kmsdoc.BaseDoc{
			ID:          kmsdoc.NewKmsDocID(kmsdoc.DocTypeCert, certID),
			NamespaceID: namespaceID,
		},
		NotAfter: certificate.NotAfter,
		Usage:    p.Usage,

		CertStorePath: blobName,
		CommonName:    certificate.Subject.CommonName,
	}
	if isRootNS {
		certDoc.IssuerNamespaceID = namespaceID
		certDoc.IssuerCertificateID = certDoc.ID
	} else {
		certDoc.IssuerNamespaceID = signerBundle.signerNamespaceID
		certDoc.IssuerCertificateID = signerBundle.signerCertId
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

func (s *PolicyCertRequestDocSection) toCertificateRequestPolicyParameters() *CertificateRequestPolicyParameters {
	if s == nil {
		return nil
	}
	policyID := s.IssuerPolicyID.String()
	return &CertificateRequestPolicyParameters{
		IssuerNamespaceID:       s.IssuerNamespaceID,
		IssuerPolicyIdentifier:  &policyID,
		KeyProperties:           &s.KeyProperties,
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
