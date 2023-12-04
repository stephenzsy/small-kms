package agentconfigmanager

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type AgentEndpointConfiguration struct {
	VerifyJWKs       []cloudkey.JsonWebKey `json:"verifyJwks"`
	VerifyJwkID      string                `json:"verifyJwkId"`
	Version          string                `json:"version"`
	TLSCertificateID string                `json:"tlsCertificateId"`
	cm               ConfigManager
}

func (c *AgentEndpointConfiguration) TLSCertificateBundleFile() string {
	return string(c.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileServerCert))
}

type endpointProcessor struct {
	cm     ConfigManager
	config AgentEndpointConfiguration
}

func (p *endpointProcessor) init(c context.Context) error {
	logger := log.Ctx(c)
	p.config.cm = p.cm
	f := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileEndpoint)
	if exists, err := f.Exists(); err != nil {
		logger.Error().Err(err).Msg("failed to check if endpoint config exists")
	} else if exists {
		if err := f.ReadJSON(&p.config); err != nil {
			logger.Error().Err(err).Msg("failed to read endpoint config")
		} else {
			p.cm.(*configManager).endpointConfigUpdate <- &p.config
		}
	}
	return nil
}

func (p *endpointProcessor) processEndpoint(c context.Context, ref *agentmodels.AgentConfigRef) error {
	requireNewCertificate, err := p.tlsCertRequireNewCertificate(c)
	if err != nil {
		return err
	}
	if p.config.Version == ref.Version && !requireNewCertificate {
		return nil
	}
	resp, err := p.cm.Client().GetAgentConfigWithResponse(c, "me", agentmodels.AgentConfigNameEndpoint)
	if err != nil {
		return nil
	} else if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}
	endpointConfig, err := resp.JSON200.AsAgentConfigEndpoint()
	if err != nil {
		return err
	}
	hasChange := false
	defer func() {
		if hasChange {
			p.cm.(*configManager).endpointConfigUpdate <- &p.config
		}
	}()

	if requireNewCertificate {

		certResp, certFile, err := enrollCert(c, p.cm, endpointConfig.TlsCertificatePolicyId)
		if err != nil {
			return err
		}
		if err := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileServerCert).LinkToAbsolutePath(certFile); err != nil {
			return err
		}
		p.config.TLSCertificateID = certResp.ID
		hasChange = true
	}
	if ref.Version != p.config.Version {
		logger := log.Ctx(c)
		logger.Info().Str("current", p.config.Version).Str("new", ref.Version).Msg("endpoint config version changed, updating")

		verifyJwks := make([]cloudkey.JsonWebKey, 0, len(endpointConfig.JwtVerifyKeyIds))
		verifyJwkID := ""
		for _, keyID := range endpointConfig.JwtVerifyKeyIds {
			if key, err := pullPublicJWK(c, p.cm, keyID); err != nil {
				return err
			} else {
				verifyJwks = append(verifyJwks, key.Jwk)
				if verifyJwkID == "" {
					verifyJwkID = keyID
				}
			}
		}

		p.config.VerifyJWKs = verifyJwks
		p.config.VerifyJwkID = verifyJwkID
		p.config.Version = ref.Version
		hasChange = true
	}

	if hasChange {

		versionedDir := p.cm.ConfigDir().Versioned(agentmodels.AgentConfigNameEndpoint, ref.Version)
		if err := versionedDir.EnsureExist(); err != nil {
			return err
		}
		if err := versionedDir.ConfigFile(configFileEndpoint).WriteJSON(p.config); err != nil {
			return err
		}
		if err := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileEndpoint).LinkToAbsolutePath(
			string(versionedDir.ConfigFile(configFileEndpoint))); err != nil {
			return err
		}
	}

	return nil
}

func (p *endpointProcessor) tlsCertRequireNewCertificate(c context.Context) (bool, error) {
	certFile := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileServerCert)
	if exists, err := certFile.Exists(); err != nil {
		return false, err
	} else if !exists {
		return true, nil
	}
	// parse pem
	certFilename := string(certFile)
	tlsCert, err := tls.LoadX509KeyPair(certFilename, certFilename)
	if err != nil {
		return true, err
	}
	cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return true, err
	}
	halfway := cert.NotBefore.Add(cert.NotAfter.Sub(cert.NotBefore) / 2)
	if time.Now().After(halfway) {
		return true, nil
	}
	return false, nil
}
