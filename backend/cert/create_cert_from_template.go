package cert

import (
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
)

func shouldCreateNewCertificate(c RequestContext, templateDoc *ct.CertificateTemplateDoc) (*CertDoc, bool, error) {
	// certDoc, err := getLatestCertificateByTemplateDoc(c, templateDoc.GetLocator())
	// if err != nil {
	// 	if !errors.Is(err, common.ErrStatusNotFound) {
	// 		return nil, false, err
	// 	}
	// 	return nil, true, nil
	// }

	// if certDoc.Deleted != nil {
	// 	log.Info().Msg("should create new cert: latest certificate is deleted")
	// 	return certDoc, true, nil
	// }
	// if certDoc.NotAfter.Time().Before(time.Now()) {
	// 	log.Info().Msg("should create new cert: latest certificate is expired")
	// 	return certDoc, true, nil
	// }
	// if !slices.Equal(certDoc.TemplateDigest, templateDoc.Digest) {
	// 	log.Info().Msg("should create new cert: template digest mismatch")
	// 	return certDoc, true, nil
	// }

	// if templateDoc.LifetimeTrigger.DaysBeforeExpiry != nil {
	// 	cutOff := certDoc.NotAfter.Time().AddDate(0, 0, 0-int(*templateDoc.LifetimeTrigger.DaysBeforeExpiry))
	// 	shouldRenew := cutOff.Before(time.Now())
	// 	log.Info().Msgf("should create new cert: eval days before expiry: %v", shouldRenew)
	// 	return certDoc, shouldRenew, nil
	// }

	// if templateDoc.LifetimeTrigger.LifetimePercentage != nil {
	// 	duration := certDoc.NotAfter.Time().Sub(certDoc.NotBefore.Time())
	// 	cutOff := certDoc.NotBefore.Time().Add(duration * time.Duration(*templateDoc.LifetimeTrigger.LifetimePercentage) / 100)
	// 	shouldRenew := cutOff.Before(time.Now())
	// 	log.Info().Msgf("should create new cert: eval lifetime percentage: %v", shouldRenew)
	// 	return certDoc, shouldRenew, nil
	// }

	// log.Info().Msg("should not create new cert: current is within range")
	// return certDoc, false, nil
	return nil, false, nil
}
