package admin

import (
	"github.com/google/uuid"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyCertIssueDocSection struct {
	IssuerID            kmsdoc.KmsDocID    `json:"issuerID"`
	AllowedRequesters   []uuid.UUID        `json:"allowedRequesters"`
	AllowedUsages       []CertificateUsage `json:"allowedUsages"`
	MaxValidityInMonths int32              `json:"max_validity_months"`
}

/*
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
	LastCertExpires time.Time               `json:"lastCertExpired"`
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

	if IsRootCANamespace(namespaceID) {
		c.IsCA = true
		c.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign
		c.MaxPathLen = 1
		c.MaxPathLenZero = false
		c.BasicConstraintsValid = true
	} else {
		err = fmt.Errorf("namespace not supported yet: %s", namespaceID.String())
		return
	}

	if csr != nil {
		c.Subject = csr.Subject
		c.Extensions = csr.Extensions
		c.EmailAddresses = csr.EmailAddresses
		c.DNSNames = csr.DNSNames
		c.IPAddresses = csr.IPAddresses
		c.URIs = csr.URIs
	} else {
		c.Subject = p.Subject.ToPkixName()
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
*/
