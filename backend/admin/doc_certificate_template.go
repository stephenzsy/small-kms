package admin

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"slices"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type CertificateTemplateDocKeyProperties struct {
	// signature algorithm
	Alg      JwkAlg     `json:"alg"`
	Kty      KeyType    `json:"kty"`
	KeySize  *KeySize   `json:"key_size,omitempty"`
	Crv      *CurveName `json:"crv,omitempty"`
	ReuseKey *bool      `json:"reuse_key,omitempty"`
}

type CertificateTemplateDocLifeTimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

type CertificateTemplateDocSubject struct {
	CertificateSubject
	cachedString *string
}

type CertificateTemplateDoc struct {
	kmsdoc.BaseDoc
	DisplayName             string                                `json:"displayName"`
	IssuerNamespaceID       uuid.UUID                             `json:"issuerNamespaceId"`
	IssuerNameSpaceType     NamespaceTypeShortName                `json:"issuerNameSpaceType"`
	IssuerTemplateID        kmsdoc.KmsDocID                       `json:"issuerTemplateId"`
	KeyProperties           CertificateTemplateDocKeyProperties   `json:"keyProperties"`
	KeyStorePath            *string                               `json:"keyStorePath,omitempty"`
	Subject                 CertificateTemplateDocSubject         `json:"subject"`
	SubjectAlternativeNames *SANsSanitized                        `json:"sans,omitempty"`
	Usage                   CertificateUsage                      `json:"usage"`
	ValidityInMonths        int32                                 `json:"validity_months"`
	LifetimeTrigger         CertificateTemplateDocLifeTimeTrigger `json:"lifetimeTrigger"`
}

func (doc *CertificateTemplateDoc) IsActive() bool {
	return doc.Deleted == nil || doc.Deleted.IsZero()
}

func (doc *CertificateTemplateDoc) IssuerCertificateDocID() kmsdoc.KmsDocID {
	return kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForPolicy, doc.IssuerTemplateID.GetUUID())
}

func (s *adminServer) readCertificateTemplateDoc(ctx context.Context, nsID uuid.UUID, templateID uuid.UUID) (*CertificateTemplateDoc, error) {
	doc := new(CertificateTemplateDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.AzCosmosContainerClient(), nsID,
		kmsdoc.NewKmsDocID(kmsdoc.DocTypeCertTemplate, templateID), doc)
	return doc, err
}

func (p *CertificateTemplateDocKeyProperties) setDefault() {
	p.Alg = AlgRS256
	p.Kty = KeyTypeRSA
	p.KeySize = ToPtr(KeySize2048)
	p.Crv = nil
}

func (p *CertificateTemplateDocKeyProperties) setRSA(alg JwkAlg, keySize KeySize) {
	p.Alg = alg
	p.Kty = KeyTypeRSA
	p.KeySize = &keySize
	p.Crv = nil
}

func (p *CertificateTemplateDocKeyProperties) setECDSA(crv CurveName) {
	p.Alg = AlgES384
	p.Kty = KeyTypeEC
	p.Crv = &crv
	p.KeySize = nil
	if crv == CurveNameP256 {
		p.Alg = AlgES256
	}
}

func (s *CertificateTemplateDocSubject) pkixName() (name pkix.Name) {
	name.CommonName = s.CN
	if s.C != nil && len(*s.C) > 0 {
		name.Country = []string{*s.C}
	}
	if s.O != nil && len(*s.O) > 0 {
		name.Organization = []string{*s.O}
	}
	if s.OU != nil && len(*s.OU) > 0 {
		name.OrganizationalUnit = []string{*s.OU}
	}
	return
}

func (s *CertificateTemplateDocSubject) String() string {
	if s == nil {
		return ""
	}
	if s.cachedString != nil {
		return *s.cachedString
	}
	name := s.pkixName()
	str := name.String()
	s.cachedString = &str
	return str
}

func (p *CertificateTemplateDocKeyProperties) fromInput(input *JwkProperties) error {
	if input == nil {
		return nil
	}
	if input.Alg == nil {
		return errors.New("alg is nil")
	}
	switch *input.Alg {
	case AlgRS256,
		AlgRS384,
		AlgRS512:
		if input.Kty != KeyTypeRSA {
			return errors.New("alg is RSA but kty is not RSA")
		}
		if input.KeySize == nil {
			p.setRSA(*input.Alg, KeySize2048)
		} else {
			p.setRSA(*input.Alg, *input.KeySize)
		}
	case AlgES256:
		if input.Crv != nil && *input.Crv != CurveNameP256 {
			return errors.New("alg is ES256 but crv is not P256")
		}
		p.setECDSA(CurveNameP256)
	case AlgES384:
		if input.Crv != nil && *input.Crv != CurveNameP256 {
			return errors.New("alg is ES384 but crv is not P384")
		}
		p.setECDSA(CurveNameP384)
	}
	return nil
}

func (t *CertificateTemplateDocLifeTimeTrigger) setDefault() {
	t.DaysBeforeExpiry = nil
	t.LifetimePercentage = ToPtr(int32(80))
}

/*
func (t *CertificateTemplateDocLifeTimeTrigger) setDisabled() {
	t.DaysBeforeExpiry = ToPtr(int32(0))
	t.LifetimePercentage = nil
}
*/

func (t *CertificateTemplateDocLifeTimeTrigger) fromInput(input *CertificateLifetimeTrigger, validityInMonths int32) error {
	if input == nil {
		return nil
	}
	if input.DaysBeforeExpiry != nil {
		if *input.DaysBeforeExpiry < 0 || *input.DaysBeforeExpiry > validityInMonths*15 {
			return errors.New("days_before_expiry must be between 0 and validity_months * 15")
		}
		t.DaysBeforeExpiry = input.DaysBeforeExpiry
		t.LifetimePercentage = nil
		return nil
	}
	if input.LifetimePercentage != nil {
		if *input.LifetimePercentage < 50 || *input.LifetimePercentage > 100 {
			return errors.New("lifetime_percentage must be between 50 and 100")
		}
		t.DaysBeforeExpiry = nil
		t.LifetimePercentage = input.LifetimePercentage
	}
	return nil
}

func (doc *CertificateTemplateDoc) toCertificateTemplate(nsType NamespaceTypeShortName) *CertificateTemplate {
	if doc == nil {
		return nil
	}
	o := new(CertificateTemplate)
	baseDocPopulateRefWithMetadata(&doc.BaseDoc, &o.Ref, nsType)
	o.DisplayName = doc.DisplayName
	o.Ref.Type = RefTypeCertificateTemplate
	o.Issuer = CertificateIssuer{
		NamespaceID:   doc.IssuerNamespaceID,
		NamespaceType: doc.IssuerNameSpaceType,
		TemplateID:    ToPtr(doc.IssuerTemplateID.GetUUID()),
	}
	o.KeyProperties = &JwkProperties{
		Alg:     ToPtr(doc.KeyProperties.Alg),
		Kty:     doc.KeyProperties.Kty,
		KeySize: doc.KeyProperties.KeySize,
		Crv:     doc.KeyProperties.Crv,
	}
	o.ReuseKey = doc.KeyProperties.ReuseKey
	o.KeyStorePath = doc.KeyStorePath
	o.LifetimeTrigger = &CertificateLifetimeTrigger{
		DaysBeforeExpiry:   doc.LifetimeTrigger.DaysBeforeExpiry,
		LifetimePercentage: doc.LifetimeTrigger.LifetimePercentage,
	}
	o.Subject = doc.Subject.CertificateSubject
	if doc.SubjectAlternativeNames != nil {
		o.SubjectAlternativeNames = &doc.SubjectAlternativeNames.CertificateSubjectAlternativeNames
	}
	o.Usage = doc.Usage
	o.ValidityInMonths = ToPtr(doc.ValidityInMonths)
	return o
}

func (doc *CertificateTemplateDoc) createAzKey(ctx context.Context, client *azkeys.Client, nsType NamespaceTypeShortName, cert *x509.Certificate) (r azkeys.KeyBundle, err error) {
	params := azkeys.CreateKeyParameters{
		KeyOps: []*azkeys.KeyOperation{to.Ptr(azkeys.KeyOperationSign), to.Ptr(azkeys.KeyOperationVerify)},
		KeyAttributes: &azkeys.KeyAttributes{
			Enabled: to.Ptr(true),
		},
	}
	kp := doc.KeyProperties
	switch kp.Kty {
	case KeyTypeRSA:
		params.Kty = to.Ptr(azkeys.KeyTypeRSA)
		if doc.KeyProperties.KeySize == nil {
			return r, fmt.Errorf("key size null for RSA key")
		}
		switch *doc.KeyProperties.KeySize {
		case KeySize2048:
			params.KeySize = to.Ptr(int32(KeySize2048))
		case KeySize3072:
			params.KeySize = to.Ptr(int32(KeySize3072))
		case KeySize4096:
			params.KeySize = to.Ptr(int32(KeySize4096))
		default:
			return r, fmt.Errorf("unsupported key size %d", *doc.KeyProperties.KeySize)
		}
	case KeyTypeEC:
		params.Kty = to.Ptr(azkeys.KeyTypeEC)
		if doc.KeyProperties.Crv == nil {
			return r, fmt.Errorf("curve null for EC key")
		}
		switch *doc.KeyProperties.Crv {
		case CurveNameP256:
			params.Curve = to.Ptr(azkeys.CurveNameP256)
		case CurveNameP384:
			params.Curve = to.Ptr(azkeys.CurveNameP384)
		default:
			return r, fmt.Errorf("unsupported curve %s", *doc.KeyProperties.Crv)
		}
	default:
		return r, fmt.Errorf("unsupported key type %s", doc.KeyProperties.Kty)
	}

	if doc.KeyStorePath == nil || len(*doc.KeyStorePath) <= 0 {
		return r, fmt.Errorf("nil key name")
	}

	switch nsType {
	case NSTypeRootCA,
		NSTypeIntCA:
		params.KeyAttributes.Exportable = to.Ptr(false)
	default:
		params.KeyAttributes.Exportable = to.Ptr(true)
	}

	if kp.ReuseKey != nil && *kp.ReuseKey {
		// try get certificate
		resp, err := client.GetKey(ctx, *doc.KeyStorePath, "", nil)
		if err != nil {
			return resp.KeyBundle, err
		}
		// verify key does not expire before certifiate
		if resp.Attributes.Expires != nil && resp.Attributes.Expires.Before(cert.NotAfter) {
			goto createKey
		}
		key := resp.Key
		// verify key parameters
		switch *key.Kty {
		case azkeys.KeyTypeEC:
			if kp.Kty != KeyTypeEC {
				goto createKey
			}
			switch *key.Crv {
			case azkeys.CurveNameP256:
				if *kp.Crv != CurveNameP256 {
					goto createKey
				}
			case azkeys.CurveNameP384:
				if *kp.Crv != CurveNameP384 {
					goto createKey
				}
			}
		case azkeys.KeyTypeRSA:
			if kp.Kty != KeyTypeRSA || len(key.N)*8 != int(*kp.KeySize) {
				goto createKey
			}
		default:
			goto createKey
		}
		// verify key ops
		if !slices.ContainsFunc(key.KeyOps, func(op *azkeys.KeyOperation) bool {
			return *op == azkeys.KeyOperationSign
		}) {
			goto createKey
		}
		return resp.KeyBundle, err
	} else {
		params.KeyAttributes.Expires = to.Ptr(cert.NotAfter)
	}

createKey:
	resp, err := client.CreateKey(ctx, *doc.KeyStorePath, params, nil)
	return resp.KeyBundle, err
}

func (doc *CertificateTemplateDoc) createAzCertificate(ctx context.Context, client *azcertificates.Client,
	nsType NamespaceTypeShortName) (azcertificates.CreateCertificateResponse, error) {
	params := azcertificates.CreateCertificateParameters{}
	x509Properties := azcertificates.X509CertificateProperties{
		Subject:          ToPtr(doc.Subject.pkixName().String()),
		ValidityInMonths: ToPtr(int32(doc.ValidityInMonths)),
	}

	keyProperties := doc.KeyProperties.getAzCertificatesKeyProperties()

	if nsType == NSTypeIntCA {
		keyProperties.Exportable = to.Ptr(false)
	} else {
		keyProperties.Exportable = to.Ptr(true)
	}

	params.CertificatePolicy = &azcertificates.CertificatePolicy{
		Attributes: &azcertificates.CertificateAttributes{
			Enabled: to.Ptr(true),
		},
		KeyProperties:             &keyProperties,
		X509CertificateProperties: &x509Properties,
		SecretProperties: &azcertificates.SecretProperties{
			ContentType: to.Ptr("application/x-pem-file"),
		},
	}

	switch nsType {
	case NSTypeServicePrincipal:
		params.Tags = map[string]*string{
			"kms-access-principal-id": to.Ptr(doc.NamespaceID.String()),
		}
	}

	return client.CreateCertificate(ctx, *doc.KeyStorePath, params, nil)
}

func (p *CertificateTemplateDocKeyProperties) getAzCertificatesKeyProperties() (r azcertificates.KeyProperties) {
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

func (s *adminServer) listCertificateTemplateDoc(ctx context.Context, nsID uuid.UUID) ([]*CertificateTemplateDoc, error) {
	partitionKey := azcosmos.NewPartitionKeyString(nsID.String())
	pager := s.AzCosmosContainerClient().NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.displayName FROM c
WHERE c.namespaceId = @namespaceId
  AND c.type = @type`,
		partitionKey, &azcosmos.QueryOptions{
			QueryParameters: []azcosmos.QueryParameter{
				{Name: "@namespaceId", Value: nsID.String()},
				{Name: "@type", Value: kmsdoc.DocTypeNameCertTemplate},
			},
		})

	return PagerToList[CertificateTemplateDoc](ctx, pager)
}
