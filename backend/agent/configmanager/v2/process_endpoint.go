package agentconfigmanager

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"time"

	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
)

type endpointProcessor struct {
	cm             ConfigManager
	currentVersion string
}

func (p *endpointProcessor) processEndpoint(c context.Context, ref *agentmodels.AgentConfigRef) error {
	requireEnroll, err := p.tlsCertRequireEnroll(c)
	if err != nil {
		return err
	}
	if ref.Version == p.currentVersion && !requireEnroll {
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

	// endpointConfig, err := resp.JSON200.AsAgentConfigEndpoint()
	// if err != nil {
	// 	return err
	// }

	// logger := log.Ctx(c)

	// // process tls cert
	// now := time.Now()

	// halfway := p.clientCertificate.NotBefore.Add(p.clientCertificate.NotAfter.Sub(p.clientCertificate.NotBefore) / 2)
	// if not.After(halfway) {

	// }

	// if now.After(halfway) {
	// 	logger.Info().Time("now", now).Time("now is past", halfway).Msg("client certificate expiring, re-enrolling")
	// 	var enrolledFileName string
	// 	enrolledCert, _, err := agentutils.EnrollCertificate(c, p.cm.client, endpointConfig.KeyCredentialCertificatePolicyId,
	// 		func(cert *certmodels.Certificate) (*os.File, error) {
	// 			enrolledFileName = p.cm.configDir.Certs().File(fmt.Sprintf("%s.pem", cert.ID))
	// 			return p.cm.configDir.Certs().OpenFile(fmt.Sprintf("%s.pem", cert.ID), os.O_CREATE|os.O_WRONLY, 0400, true)
	// 		}, false)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	addEntraKeyResp, err := p.cm.client.AddMsEntraKeyCredentialWithResponse(c,
	// 		models.NamespaceProviderServicePrincipal, "me", enrolledCert.ID)

	// 	if err != nil {
	// 		return err
	// 	} else if addEntraKeyResp.StatusCode() < 200 || addEntraKeyResp.StatusCode() >= 300 {
	// 		return fmt.Errorf("unexpected status code: %d", addEntraKeyResp.StatusCode())
	// 	}
	// 	// create config file

	// 	linkFileName := p.cm.configDir.Config(agentmodels.AgentConfigNameIdentity).ConfigFile(configFileClientCert, false)
	// 	if _, err := os.Lstat(filepath.Dir(linkFileName)); err != nil {
	// 		if errors.Is(err, os.ErrNotExist) {
	// 			if err := os.MkdirAll(filepath.Dir(linkFileName), 0750); err != nil {
	// 				return err
	// 			}
	// 		} else {
	// 			return err
	// 		}
	// 	}
	// 	if _, err := os.Lstat(linkFileName); err == nil {
	// 		// delete ink
	// 		if err := os.Remove(linkFileName); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	relpath, err := filepath.Rel(filepath.Dir(linkFileName), enrolledFileName)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	logger.Debug().Str("relpath", relpath).Msg("create symlink")
	// 	if err := os.Symlink(relpath, linkFileName); err != nil {
	// 		return err
	// 	}
	// 	p.cm.configureClient(linkFileName)
	// }

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
