package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *adminServer) shouldCreateCertificateForTemplate(ctx context.Context, nsID uuid.UUID, templateDoc *CertificateTemplateDoc, certDoc *CertDoc) (renewReason string) {
	// load existing certificate
	if !certDoc.IsActive() {
		return "existing certificate does not exist or is not active"
	}

	// verify template matches certificate metadata
	if certDoc.TemplateID != templateDoc.ID {
		return "template mismatch"
	}
	if certDoc.IssuerNamespaceID != templateDoc.IssuerNamespaceID {
		return "issuer namespace mismatch"
	}
	if certDoc.SubjectBase != templateDoc.Subject.String() {
		return "subject mismatch"
	}
	// if certDoc.KeyInfo.Alg == nil || *certDoc.KeyInfo.Alg != templateDoc.KeyProperties.Alg ||
	// 	certDoc.KeyInfo.Kty == models.KeyTypeRSA && (certDoc.KeyInfo.KeySize == nil || templateDoc.KeyProperties.KeySize == nil ||
	// 		*certDoc.KeyInfo.KeySize != *templateDoc.KeyProperties.KeySize) {
	// 	return "alg or key mismatch"
	// }
	if certDoc.Usage != templateDoc.Usage {
		return "usage mismatch"
	}

	// verify life time
	if templateDoc.LifetimeTrigger.DaysBeforeExpiry != nil {
		daysBeforeExpiry := *templateDoc.LifetimeTrigger.DaysBeforeExpiry
		if daysBeforeExpiry > 0 && time.Now().AddDate(0, 0, int(daysBeforeExpiry)).
			After(certDoc.NotAfter) {
			return "within days before expiry"
		}
	} else if templateDoc.LifetimeTrigger.LifetimePercentage != nil {
		p := *templateDoc.LifetimeTrigger.LifetimePercentage
		if time.Now().
			After(certDoc.NotBefore.
				Add(certDoc.NotAfter.Sub(certDoc.NotBefore) * time.Duration(p) / 100)) {
			return "outside lifetime percentage"
		}
	}
	return
}
