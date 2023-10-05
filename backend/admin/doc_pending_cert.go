package admin

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyDocSectionIssuerProperties struct {
	IssuerNamespaceID uuid.UUID `json:"issuerNamespaceId"`
	IssuerPolicyID    uuid.UUID `json:"issuerPolicyId"`
}

type PendingCertDoc struct {
	kmsdoc.BaseDoc
	Expires                 time.Time                           `json:"exp"`                 // indicates pending status expires
	TemplateNamespaceID     uuid.UUID                           `json:"templateNamespaceId"` // this is informational only
	TemplateID              kmsdoc.KmsDocID                     `json:"templateId"`          // this is informational only
	JWT                     [3]string                           `json:"jwt"`                 // base64url encoded JWT segments
	RequesterID             uuid.UUID                           `json:"requesterId"`         // must be matched to issue certificate
	KeyProperties           CertificateTemplateDocKeyProperties `json:"keyProperties"`
	Usage                   CertificateUsage                    `json:"usage"`
	Subject                 CertificateTemplateDocSubject       `json:"subject"`
	SubjectAlternativeNames *CertificateSubjectAlternativeNames `json:"sans,omitempty"`
	NotBefore               time.Time                           `json:"notBefore"`
	NotAfter                time.Time                           `json:"notAfter"`
	IssuerNamespaceID       uuid.UUID                           `json:"issuerNamespaceId"`
	IssuerTemplateID        kmsdoc.KmsDocID                     `json:"issuerTemplateId"`

	Issued time.Time `json:"issued"`
}

func (s *adminServer) readPendingCertDoc(ctx context.Context, nsID uuid.UUID, docID kmsdoc.KmsDocID) (*PendingCertDoc, error) {
	doc := new(PendingCertDoc)
	err := kmsdoc.AzCosmosRead(ctx, s.AzCosmosContainerClient(), nsID, docID, doc)
	return doc, common.WrapAzRsNotFoundErr(err, fmt.Sprintf("%s:cert:%s", nsID, docID))
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
	cert *x509.Certificate,
	claimsEncoded string,
	templateDoc *CertificateTemplateDoc,
	issueToNamespaceID uuid.UUID,
	requesterID uuid.UUID) PendingCertDoc {
	// fisrtOrNil := func(s []string) *string {
	// 	if len(s) < 1 {
	// 		return nil
	// 	}
	// 	return &s[0]
	// }
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
		Usage:               templateDoc.Usage,
		Subject:             CertificateTemplateDocSubject{
			// CertificateSubject: CertificateSubject{
			// 	CN: cert.Subject.CommonName,
			// 	OU: fisrtOrNil(cert.Subject.OrganizationalUnit),
			// 	O:  fisrtOrNil(cert.Subject.Organization),
			// 	C:  fisrtOrNil(cert.Subject.Country),
			// },
		},
		SubjectAlternativeNames: certificateSubjectAlternativeNamesToDoc(cert),
		NotBefore:               cert.NotBefore,
		NotAfter:                cert.NotAfter,
		IssuerNamespaceID:       templateDoc.IssuerNamespaceID,
		IssuerTemplateID:        templateDoc.IssuerTemplateID,
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
