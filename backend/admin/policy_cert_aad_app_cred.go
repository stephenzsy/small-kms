package admin

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/rs/zerolog/log"

	"github.com/stephenzsy/small-kms/backend/kmsdoc"
)

type PolicyCertAadAppCredDocSection struct {
	CertID kmsdoc.KmsDocID `json:"certId"`
}

func (t *PolicyCertAadAppCredDocSection) validateAndFillWithParameters(p *CertificateAadAppCredPolicyParameters) error {
	if p == nil {
		return errors.New("missing CertRequest property")
	}

	t.CertID = p.CertificateIdentifier.docID()

	return nil
}

type PolicyCertAadAppCredAction string

const (
	PolicyCertAadAppCredActionInstallCert PolicyCertAadAppCredAction = "install-client-credential-certificate"
)

type PolicyStateCertAadAppCredDocSection struct {
	CertInstalledCUID kmsdoc.KmsDocID            `json:"installedCertId"`
	LastAction        PolicyCertAadAppCredAction `json:"lastAction"`
}

/*
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
*/
func (p *PolicyCertAadAppCredDocSection) action(ctx *gin.Context, s *adminServer, namespaceID uuid.UUID, policyDoc *PolicyDoc) (resultDoc *PolicyStateDoc, err error) {

	policyID := policyDoc.GetUUID()
	dirDoc, err := s.getDirectoryObjectDoc(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	certDoc, err := s.getCertDoc(ctx, namespaceID, p.CertID)
	if err != nil {
		return nil, err
	}
	pemBlob, err := s.FetchCertificatePEMBlob(ctx, certDoc.CertStorePath)
	if err != nil {
		return nil, err
	}
	requestBody := graphmodels.NewApplication()
	keyCredential := graphmodels.NewKeyCredential()
	keyCredential.SetTypeEscaped(ToPtr("AsymmetricX509Cert"))
	keyCredential.SetUsage(ToPtr("Verify"))
	keyCredential.SetKey(pemBlob)

	requestBody.SetKeyCredentials([]graphmodels.KeyCredentialable{keyCredential})

	_, err = s.msGraphClient.Applications().ByApplicationId(dirDoc.ID.GetUUID().String()).Patch(ctx, requestBody, nil)
	if err != nil {
		return nil, err
	}

	certDocCUID := certDoc.GetCUID()

	log.Info().Msgf("Certificate installed %s to application %s", certDocCUID.String(), *dirDoc.AppID)

	// record policy state
	resultDoc = &PolicyStateDoc{
		BaseDoc: kmsdoc.BaseDoc{
			ID:          kmsdoc.NewKmsDocID(kmsdoc.DocTypePolicyState, policyID),
			NamespaceID: namespaceID,
		},
		PolicyType: PolicyTypeCertRequest,
		Status:     PolicyStateStatusSuccess,
		Message:    fmt.Sprintf("Certificate installed %s to application %s", certDocCUID.String(), *dirDoc.AppID),
		CertAadAppCred: &PolicyStateCertAadAppCredDocSection{
			CertInstalledCUID: certDocCUID,
			LastAction:        PolicyCertAadAppCredActionInstallCert,
		},
	}
	err = kmsdoc.AzCosmosUpsert(ctx, s.azCosmosContainerClientCerts, resultDoc)
	if err != nil {
		return
	}
	log.Info().Msgf("CertRequest completed for %s/%s", namespaceID, policyID)
	return
}

func (s *PolicyCertAadAppCredDocSection) toCertificateAadAppPolicyParameters() *CertificateAadAppCredPolicyParameters {
	if s == nil {
		return nil
	}
	return &CertificateAadAppCredPolicyParameters{
		CertificateIdentifier: docIDtoCertIdentifier(s.CertID),
	}
}
