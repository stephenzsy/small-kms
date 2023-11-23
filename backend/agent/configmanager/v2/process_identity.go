package agentconfigmanager

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	now := time.Now()
	halfway := p.clientCertificate.NotBefore.Add(p.clientCertificate.NotAfter.Sub(p.clientCertificate.NotBefore) / 2)

	if now.After(halfway) {
		logger.Info().Time("now", now).Time("now is past", halfway).Msg("client certificate expiring, re-enrolling")
		var enrolledFileName string
		enrolledCert, _, err := agentutils.EnrollCertificate(c, p.cm.client, identityConfig.KeyCredentialCertificatePolicyId,
			func(cert *certmodels.Certificate) (*os.File, error) {
				enrolledFileName = p.cm.configDir.Certs().File(fmt.Sprintf("%s.pem", cert.ID))
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
		// create config file

		linkFileName := p.cm.configDir.Config(agentmodels.AgentConfigNameIdentity).ConfigFile(configFileClientCert, false)
		if _, err := os.Lstat(filepath.Dir(linkFileName)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				if err := os.MkdirAll(filepath.Dir(linkFileName), 0750); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		if _, err := os.Lstat(linkFileName); err == nil {
			// delete ink
			if err := os.Remove(linkFileName); err != nil {
				return err
			}
		}
		relpath, err := filepath.Rel(filepath.Dir(linkFileName), enrolledFileName)
		if err != nil {
			return err
		}
		logger.Debug().Str("relpath", relpath).Msg("create symlink")
		if err := os.Symlink(relpath, linkFileName); err != nil {
			return err
		}
		p.cm.configureClient(linkFileName)
	}

	return nil
}
