package agentconfigmanager

import (
	"context"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type identityProcessor struct {
	cm                *configManager
	clientCertificate *x509.Certificate
}

func (p *identityProcessor) processIdentity(c context.Context, ref *agentmodels.AgentConfigRef) error {
	now := time.Now()
	halfway := p.clientCertificate.NotBefore.Add(p.clientCertificate.NotAfter.Sub(p.clientCertificate.NotBefore) / 2)

	if !now.After(halfway) {
		return nil
	}
	logger := log.Ctx(c)
	logger.Info().Time("now", now).Time("now is past", halfway).Msg("client certificate expiring, re-enrolling")

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

	// create config file
	cert, certFile, err := enrollCert(c, p.cm.cryptoProvider,
		p.cm, identityConfig.KeyCredentialCertificatePolicyId)
	if err != nil {
		return err
	}

	// add as entra key credential
	addEntraKeyResp, err := p.cm.client.AddMsEntraKeyCredentialWithResponse(c,
		models.NamespaceProviderServicePrincipal, "me", cert.ID)
	if err != nil {
		return err
	} else if addEntraKeyResp.StatusCode() < 200 || addEntraKeyResp.StatusCode() >= 300 {
		return fmt.Errorf("unexpected status code: %d", addEntraKeyResp.StatusCode())
	}

	activeConfigDir := p.cm.configDir.Active(agentmodels.AgentConfigNameIdentity)
	if err := activeConfigDir.EnsureExist(); err != nil {
		return err
	}
	certLink := activeConfigDir.ConfigFile(configFileClientCert)
	if err := certLink.LinkToAbsolutePath(certFile); err != nil {
		return err
	}

	// swap client with new credentials
	p.cm.configureClient(string(certLink))

	return nil
}
