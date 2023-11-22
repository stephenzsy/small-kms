package agentconfigmanager

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	agentutils "github.com/stephenzsy/small-kms/backend/agent/utils"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
)

type identityProcessor struct {
	cm                *configManager
	currentVersion    string
	clientCertificate *x509.Certificate
}

func (p *identityProcessor) processIdentity(c context.Context, ref *agentmodels.AgentConfigRef) error {
	if ref.Version == p.currentVersion {
		return nil
	}
	resp, err := p.cm.client.GetAgentConfigWithResponse(c, "me", agentmodels.AgentConfigNameIdentity)
	if err != nil {
		return nil
	} else if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	identityConfig, err := resp.JSON200.AsAgentConfigIdentity()
	if err != nil {
		return err
	}

	logger := log.Ctx(c)
	logger.Debug().
		Time("notBefore", p.clientCertificate.NotBefore).
		Time("notAfter", p.clientCertificate.NotAfter).
		Msg("check identity")

	if time.Now().After(p.clientCertificate.NotBefore.Add(time.Hour)) {
		logger.Info().Msg("client certificate expiring, re-enrolling")
		enrolledCert, _, err := agentutils.EnrollCertificate(c, p.cm.client, identityConfig.KeyCredentialCertificatePolicyId,
			func(cert *certmodels.Certificate) (*os.File, error) {
				return p.cm.configDir.Certs().OpenFile(fmt.Sprintf("%s.pem", cert.ID), os.O_CREATE|os.O_WRONLY, 0400, true)
			}, false)
		if err != nil {
			return err
		}

		addEntraKeyResp, err := p.cm.client.AddMsEntraKeyCredentialWithResponse(c,
			models.NamespaceProviderServicePrincipal, "me", enrolledCert.ID)

		if err != nil {
			return err
		} else if addEntraKeyResp.StatusCode() < 200 || addEntraKeyResp.StatusCode() >= 300 {
			return fmt.Errorf("unexpected status code: %d", addEntraKeyResp.StatusCode())
		}
	}

	return nil
}
