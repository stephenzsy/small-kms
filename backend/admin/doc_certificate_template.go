package admin

import (
	"context"
	"crypto/x509/pkix"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type CertificateTemplateDocKeyProperties struct {
	// signature algorithm
	Alg models.JwkAlg `json:"alg"`
	// Kty      KeyType       `json:"kty"`
	// KeySize  *KeySize      `json:"key_size,omitempty"`
	// Crv      *CurveName    `json:"crv,omitempty"`
	ReuseKey *bool `json:"reuse_key,omitempty"`
}

type CertificateTemplateDocLifeTimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

type CertificateTemplateDocSubject struct {
	// CertificateSubject
	cachedString *string
}

type CertificateTemplateDoc struct {
	kmsdoc.BaseDoc
	DisplayName       string    `json:"displayName"`
	IssuerNamespaceID uuid.UUID `json:"issuerNamespaceId"`
	// IssuerNameSpaceType     NamespaceTypeShortName                `json:"issuerNameSpaceType"`
	IssuerTemplateID        kmsdoc.KmsDocID                       `json:"issuerTemplateId"`
	KeyProperties           CertificateTemplateDocKeyProperties   `json:"keyProperties"`
	KeyStorePath            *string                               `json:"keyStorePath,omitempty"`
	Subject                 CertificateTemplateDocSubject         `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames   `json:"sans,omitempty"`
	Usage                   CertificateUsage                      `json:"usage"`
	ValidityInMonths        int32                                 `json:"validity_months"`
	LifetimeTrigger         CertificateTemplateDocLifeTimeTrigger `json:"lifetimeTrigger"`
	VariablesEnabeled       bool                                  `json:"variablesEnabeled"`
}

func (doc *CertificateTemplateDoc) IsActive() bool {
	return doc.Deleted == nil || doc.Deleted.IsZero()
}

func (doc *CertificateTemplateDoc) IssuerCertificateDocID() kmsdoc.KmsDocID {
	return kmsdoc.NewKmsDocID(kmsdoc.DocTypeLatestCertForTemplate, doc.IssuerTemplateID.GetUUID())
}

func (p *CertificateTemplateDocKeyProperties) setDefault() {
	p.Alg = models.AlgRS256
	// p.Kty = KeyTypeRSA
	// p.KeySize = ToPtr(KeySize2048)
	// p.Crv = nil
}

func (p *CertificateTemplateDocKeyProperties) setRSA(alg models.JwkAlg, keySize int) {
	p.Alg = alg
	// p.Kty = KeyTypeRSA
	// p.KeySize = &keySize
	// p.Crv = nil
}

func (p *CertificateTemplateDocKeyProperties) setECDSA(crv models.JwtCrv) {
	p.Alg = models.AlgES384
	// p.Kty = KeyTypeEC
	// p.Crv = &crv
	// p.KeySize = nil
	// if crv == CurveNameP256 {
	// 	p.Alg = models.AlgES256
	// }
}

func (s *CertificateTemplateDocSubject) pkixName() (name pkix.Name) {
	// name.CommonName = s.CN
	// if s.C != nil && len(*s.C) > 0 {
	// 	name.Country = []string{*s.C}
	// }
	// if s.O != nil && len(*s.O) > 0 {
	// 	name.Organization = []string{*s.O}
	// }
	// if s.OU != nil && len(*s.OU) > 0 {
	// 	name.OrganizationalUnit = []string{*s.OU}
	// }
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

func (p *CertificateTemplateDocKeyProperties) fromJwkProperties(input *models.JwkProperties) error {
	// if input == nil {
	// 	return nil
	// }
	// if input.Alg == nil {
	// 	return errors.New("alg is nil")
	// }
	// switch *input.Alg {
	// case models.AlgRS256,
	// 	models.AlgRS384,
	// 	models.AlgRS512:
	// 	if input.Kty != KeyTypeRSA {
	// 		return errors.New("alg is RSA but kty is not RSA")
	// 	}
	// 	if input.KeySize == nil {
	// 		p.setRSA(*input.Alg, KeySize2048)
	// 	} else {
	// 		p.setRSA(*input.Alg, *input.KeySize)
	// 	}
	// case models.AlgES256:
	// 	if input.Crv != nil && *input.Crv != CurveNameP256 {
	// 		return errors.New("alg is ES256 but crv is not P256")
	// 	}
	// 	p.setECDSA(CurveNameP256)
	// case models.AlgES384:
	// 	if input.Crv != nil && *input.Crv != CurveNameP256 {
	// 		return errors.New("alg is ES384 but crv is not P384")
	// 	}
	// 	p.setECDSA(CurveNameP384)
	// }
	return nil
}

func (p *CertificateTemplateDocKeyProperties) populateJwkProperties(o *models.JwkProperties) {
	// if p == nil {
	// 	return
	// }
	// o.Alg = utils.ToPtr(p.Alg)
	// o.Kty = p.Kty
	// o.KeySize = p.KeySize
	// o.Crv = p.Crv
}

func (t *CertificateTemplateDocLifeTimeTrigger) setDefault() {
	t.DaysBeforeExpiry = nil
	t.LifetimePercentage = ToPtr(int32(80))
}

func (doc *CertificateTemplateDoc) createAzCertificate(
	ctx context.Context,
	client *azcertificates.Client,
	issueToNamespaceID uuid.UUID,
	subject string) (azcertificates.CreateCertificateResponse, error) {
	params := azcertificates.CreateCertificateParameters{}
	x509Properties := azcertificates.X509CertificateProperties{
		Subject:          utils.ToPtr(subject),
		ValidityInMonths: ToPtr(int32(doc.ValidityInMonths)),
	}

	keyExportable := !isAllowedCaNamespace(issueToNamespaceID)
	keyProperties := doc.KeyProperties.getAzCertificatesKeyProperties(keyExportable)

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

	params.Tags = map[string]*string{
		"kms-access-principal-id": to.Ptr(issueToNamespaceID.String()),
	}

	return client.CreateCertificate(ctx, *doc.KeyStorePath, params, nil)
}

func (p *CertificateTemplateDocKeyProperties) getAzCertificatesKeyProperties(keyExportable bool,
) (r azcertificates.KeyProperties) {
	r.KeyType = ToPtr(azcertificates.KeyTypeRSA)
	r.KeySize = ToPtr(int32(2048))
	r.ReuseKey = p.ReuseKey
	// switch p.Kty {
	// case KeyTypeRSA:
	// 	if p.KeySize != nil {
	// 		switch *p.KeySize {
	// 		case KeySize3072:
	// 			r.KeySize = ToPtr(int32(3072))
	// 		case KeySize4096:
	// 			r.KeySize = ToPtr(int32(4096))
	// 		}
	// 	}
	// case KeyTypeEC:
	// 	r.KeyType = ToPtr(azcertificates.KeyTypeEC)
	// 	r.KeySize = nil
	// 	r.Curve = ToPtr(azcertificates.CurveNameP256)
	// 	if p.Crv != nil {
	// 		switch *p.Crv {
	// 		case CurveNameP384:
	// 			r.Curve = ToPtr(azcertificates.CurveNameP384)
	// 		}
	// 	}
	// }
	r.Exportable = to.Ptr(keyExportable)
	return
}
