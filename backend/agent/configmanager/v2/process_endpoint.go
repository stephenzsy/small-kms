package agentconfigmanager

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	certmodels "github.com/stephenzsy/small-kms/backend/models/cert"
)

type AgentEndpointConfiguration struct {
	VerifyJWKs        []cloudkey.JsonWebKey `json:"verifyJwks"`
	VerifyJwkID       string                `json:"verifyJwkId"`
	Version           string                `json:"version"`
	TLSCertificateID  string                `json:"tlsCertificateId"`
	AllowedImageRepos []string              `json:"allowedImageRepos"`
}

type endpointProcessor struct {
	cm     ConfigManager
	config AgentEndpointConfiguration
}

func (p *endpointProcessor) init(c context.Context) error {
	logger := log.Ctx(c)
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
	certExpiring, err := p.tlsCertificateExpiringSoon(c)
	if err != nil {
		return err
	}
	if p.config.Version == ref.Version && !certExpiring {
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

	logger := log.Ctx(c)

	if endpointConfig.TLSCertificateAutoEnroll {
		if certExpiring {
			certResp, certFile, err := enrollCert(c, p.cm.CryptoProvider(),
				p.cm, endpointConfig.TlsCertificatePolicyId)
			if err != nil {
				return err
			}
			if err := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileServerCert).LinkToAbsolutePath(certFile); err != nil {
				return err
			}
			p.config.TLSCertificateID = certResp.ID
			hasChange = true
		}
	} else if p.config.TLSCertificateID != endpointConfig.TLSCertificateID {
		// get certificate
		jwk, err := cloudkey.NewEphemeralECDHJwk(p.cm.CryptoProvider())
		if err != nil {
			return err
		}
		resp, err := p.cm.Client().GetCertificateSecretWithResponse(c,
			models.NamespaceProviderServicePrincipal,
			"me", endpointConfig.TLSCertificateID, certmodels.CertificateSecretRequest{
				Jwk: *jwk,
			})
		if err != nil {
			return err
		} else if resp.StatusCode() != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}
		jwe, err := cloudkey.NewJsonWebEncryption(resp.JSON200.Payload)
		if err != nil {
			return err
		}
		if pemBytes, _, err := jwe.Decrypt(func(*cloudkey.JoseHeader) (crypto.PrivateKey, error) {
			return jwk.PrivateKey().(*ecdsa.PrivateKey).ECDH()
		}); err != nil {
			return err
		} else {
			certFile, err := writeCert(c, p.cm, endpointConfig.TLSCertificateID, pemBytes)
			if err != nil {
				return err
			}
			if err := p.cm.ConfigDir().Active(agentmodels.AgentConfigNameEndpoint).ConfigFile(configFileServerCert).LinkToAbsolutePath(certFile); err != nil {
				return err
			}
			p.config.TLSCertificateID = endpointConfig.TLSCertificateID
			hasChange = true
		}
	}
	if ref.Version != p.config.Version {
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
		p.config.AllowedImageRepos = endpointConfig.AllowedImageRepos
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

func (p *endpointProcessor) tlsCertificateExpiringSoon(c context.Context) (bool, error) {
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
