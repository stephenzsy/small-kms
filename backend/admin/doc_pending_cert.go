package admin

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyDocSectionIssuerProperties struct {
	IssuerNamespaceID uuid.UUID `json:"issuerNamespaceId"`
	IssuerPolicyID    uuid.UUID `json:"issuerPolicyId"`
}

type PendingCertDoc struct {
	kmsdoc.BaseDoc
	Expires             time.Time                           `json:"exp"` // indicates pending status expires
	TemplateNamespaceID uuid.UUID                           `json:"templateNamespaceId"`
	TemplateID          kmsdoc.KmsDocID                     `json:"templateId"`
	JWT                 [3]string                           `json:"jwt"`         // base64url encoded JWT segments
	RequesterID         uuid.UUID                           `json:"requesterId"` // must be matched to issue certificate
	KeyProperties       CertificateTemplateDocKeyProperties `json:"keyProperties"`

	Issued time.Time `json:"issued"`
}

func encodeJwtJsonSegment[D any](jsonObj D) (string, error) {
	marshalled, err := json.Marshal(jsonObj)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(marshalled), nil
}

func newPendingCertDoc(
	certID uuid.UUID,
	claimsEncoded string,
	templateDoc *CertificateTemplateDoc,
	issueToNamespaceID uuid.UUID,
	requesterID uuid.UUID) PendingCertDoc {
	return PendingCertDoc{
		BaseDoc: kmsdoc.BaseDoc{
			ID:          kmsdoc.NewKmsDocID(kmsdoc.DocTypePendingCert, certID),
			NamespaceID: issueToNamespaceID,
		},
		Expires:             time.Now().Add(time.Hour * 24 * 7),
		JWT:                 [3]string{"", claimsEncoded, ""},
		TemplateNamespaceID: templateDoc.NamespaceID,
		TemplateID:          templateDoc.ID,
		RequesterID:         requesterID,
		KeyProperties:       templateDoc.KeyProperties,
	}
}

func (doc *PendingCertDoc) toReceipt(nsType NamespaceTypeShortName) *CertificateEnrollmentReceipt {
	r := CertificateEnrollmentReceipt{}
	baseDocPopulateRefWithMetadata(&doc.BaseDoc, &r.Ref)
	r.Ref.Type = RefTypeCertificateEnrollReceipt
	r.Expires = doc.Expires
	r.TemplateID = doc.TemplateID.GetUUID()
	r.TemplateNamespaceID = doc.TemplateNamespaceID
	r.JwtClaims = doc.JWT[1]
	r.RequesterID = doc.RequesterID
	doc.KeyProperties.populateJwkProperties(&r.KeyProperties)
	return &r
}
