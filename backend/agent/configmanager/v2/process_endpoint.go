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
	VerifyJWKs []cloudkey.JsonWebKey `json:"verifyJwks"`
	Version    string                `json:"version"`
}

type endpointProcessor struct {
	cm                   ConfigManager
	currentConfiguration *AgentEndpointConfiguration
}

func (p *endpointProcessor) init(c context.Context) error {
	logger := log.Ctx(c)
	f := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileEndpoint)
	if exists, err := f.Exists(); err != nil {
		logger.Error().Err(err).Msg("failed to check if endpoint config exists")
	} else if exists {
		config := &AgentEndpointConfiguration{}
		if err := f.ReadJSON(config); err != nil {
			logger.Error().Err(err).Msg("failed to read endpoint config")
		} else {
			p.currentConfiguration = config
		}
	}
	return nil
}

func (p *endpointProcessor) processEndpoint(c context.Context, ref *agentmodels.AgentConfigRef) error {
	requireEnroll, err := p.tlsCertRequireEnroll(c)
	if err != nil {
		return err
	}
	currentVersion := ""
	if p.currentConfiguration != nil {
		currentVersion = p.currentConfiguration.Version
	}
	if currentVersion == ref.Version && !requireEnroll {
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

	if requireEnroll {
		_, certFile, err := enrollCert(c, p.cm, endpointConfig.TlsCertificatePolicyId)
		if err != nil {
			return err
		}
		if err := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileServerCert).LinkToAbsolutePath(certFile); err != nil {
			return err
		}
	}
	if ref.Version == currentVersion {
		return nil
	}
	logger := log.Ctx(c)
	logger.Info().Str("current", currentVersion).Str("new", ref.Version).Msg("endpoint config version changed, updating")
	config := &AgentEndpointConfiguration{
		VerifyJWKs: make([]cloudkey.JsonWebKey, len(endpointConfig.JwtVerifyKeyIds)),
		Version:    ref.Version,
	}
	for _, keyID := range endpointConfig.JwtVerifyKeyIds {
		if key, err := pullPublicJWK(c, p.cm, keyID); err != nil {
			return err
		} else {
			config.VerifyJWKs = append(config.VerifyJWKs, key.Jwk)
		}
	}
	versionedDir := p.cm.ConfigDir().Versioned(agentmodels.AgentConfigNameEndpoint, ref.Version)
	if err := versionedDir.EnsureExist(); err != nil {
		return err
	}
	if err := versionedDir.ConfigFile(configFileEndpoint).WriteJSON(config); err != nil {
		return err
	}
	if err := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileEndpoint).LinkToAbsolutePath(
		string(versionedDir.ConfigFile(configFileEndpoint))); err != nil {
		return err
	}

	p.currentConfiguration = config

	return nil
}

func (p *endpointProcessor) tlsCertRequireEnroll(c context.Context) (bool, error) {
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
