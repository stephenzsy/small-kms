package admin

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type CertificateTemplateDocKeyProperties struct {
	// signature algorithm
	Alg     JwkAlg     `json:"alg"`
	Kty     KeyType    `json:"kty"`
	KeySize *KeySize   `json:"key_size,omitempty"`
	Crv     *CurveName `json:"crv,omitempty"`
}

type CertificateTemplateDocLifeTimeTrigger struct {
	Disabled           bool   `json:"disabled"`
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

type CertificateTemplateDoc struct {
	kmsdoc.BaseDoc
	DisplayName             string                                `json:"displayName"`
	IssuerNamespaceID       uuid.UUID                             `json:"issuerNamespaceId"`
	IssuerNameSpaceType     NamespaceTypeShortName                `json:"issuerNameSpaceType"`
	IssuerTemplateID        kmsdoc.KmsDocID                       `json:"issuerTemplateId"`
	KeyProperties           CertificateTemplateDocKeyProperties   `json:"keyProperties"`
	KeyStorePath            *string                               `json:"keyStorePath,omitempty"`
	Subject                 CertificateSubject                    `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames   `json:"subjectAlternativeNames,omitempty"`
	Usage                   CertificateUsage                      `json:"usage"`
	ValidityInMonths        int32                                 `json:"validity_months"`
	LifetimeTrigger         CertificateTemplateDocLifeTimeTrigger `json:"lifetimeTrigger"`
}

func (s *adminServer) readCertificateTemplateDoc(ctx context.Context, nsID uuid.UUID, templateID uuid.UUID) (*CertificateTemplateDoc, error) {
	doc := new(CertificateTemplateDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.azCosmosContainerClientCerts, nsID,
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

func (p *CertificateTemplateDocKeyProperties) fromInput(input *JwkKeyProperties) error {
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
	t.Disabled = false
	t.DaysBeforeExpiry = nil
	t.LifetimePercentage = ToPtr(int32(80))
}

func (t *CertificateTemplateDocLifeTimeTrigger) setDisabled() {
	t.Disabled = true
	t.DaysBeforeExpiry = nil
	t.LifetimePercentage = nil
}

func (t *CertificateTemplateDocLifeTimeTrigger) fromInput(input *CertificateLifetimeTrigger, validityInMonths int32) error {
	if input == nil {
		return nil
	}
	if input.Disabled != nil && *input.Disabled {
		t.setDisabled()
		return nil
	}
	if input.DaysBeforeExpiry != nil {
		if *input.DaysBeforeExpiry < 0 || *input.DaysBeforeExpiry > validityInMonths*15 {
			return errors.New("days_before_expiry must be between 0 and validity_months * 15")
		}
		t.Disabled = false
		t.DaysBeforeExpiry = input.DaysBeforeExpiry
		t.LifetimePercentage = nil
		return nil
	}
	if input.LifetimePercentage != nil {
		if *input.LifetimePercentage < 50 || *input.LifetimePercentage > 100 {
			return errors.New("lifetime_percentage must be between 50 and 100")
		}
		t.Disabled = false
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
	baseDocPopulateRef(&doc.BaseDoc, &o.Ref, nsType)
	o.Ref.DisplayName = doc.DisplayName
	o.Ref.Type = RefTypeCertificateTemplate
	o.Issuer = CertificateIssuer{
		NamespaceID:   doc.IssuerNamespaceID,
		NamespaceType: doc.IssuerNameSpaceType,
		TemplateID:    ToPtr(doc.IssuerTemplateID.GetUUID()),
	}
	o.KeyProperties = &JwkKeyProperties{
		Alg:     ToPtr(doc.KeyProperties.Alg),
		Kty:     doc.KeyProperties.Kty,
		KeySize: doc.KeyProperties.KeySize,
		Crv:     doc.KeyProperties.Crv,
	}
	o.KeyStorePath = doc.KeyStorePath
	o.LifetimeTrigger = &CertificateLifetimeTrigger{
		Disabled:           ToPtr(doc.LifetimeTrigger.Disabled),
		DaysBeforeExpiry:   doc.LifetimeTrigger.DaysBeforeExpiry,
		LifetimePercentage: doc.LifetimeTrigger.LifetimePercentage,
	}
	o.Subject = doc.Subject
	o.SubjectAlternativeNames = doc.SubjectAlternativeNames
	o.Usage = doc.Usage
	o.ValidityInMonths = ToPtr(doc.ValidityInMonths)
	return o
}

func (s *adminServer) listCertificateTemplateDoc(ctx context.Context, nsID uuid.UUID) ([]*CertificateTemplateDoc, error) {
	partitionKey := azcosmos.NewPartitionKeyString(nsID.String())
	pager := s.azCosmosContainerClientCerts.NewQueryItemsPager(`SELECT `+kmsdoc.GetBaseDocQueryColumns("c")+`,c.displayName FROM c
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
