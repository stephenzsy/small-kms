package bootstrap

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	agentclient "github.com/stephenzsy/small-kms/backend/agent/client"
	agentcommon "github.com/stephenzsy/small-kms/backend/agent/common"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/common"
)

func parseCertificateKeyPair(filename string) (cert *x509.Certificate, key crypto.PrivateKey, err error) {
	bad := func(e error) (*x509.Certificate, crypto.PrivateKey, error) {
		return nil, nil, e
	}
	if fileContent, err := os.ReadFile(filename); err != nil {
		return bad(err)
	} else {
		for block, rest := pem.Decode(fileContent); block != nil; block, rest = pem.Decode(rest) {
			if block.Type == "CERTIFICATE" {
				if cert == nil {
					if cert, err = x509.ParseCertificate(block.Bytes); err != nil {
						return bad(err)
					}
				}
			} else if block.Type == "PRIVATE KEY" {
				if key, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
					return bad(err)
				}
			}
		}
		return cert, key, nil
	}
}

func (*ServicePrincipalBootstraper) BootstarpActiveServer(c context.Context) error {
	var client *agentclient.ClientWithResponses
	var cred azcore.TokenCredential
	var tenantID string
	envSvc := common.NewEnvService()
	if baseUrl, ok := envSvc.RequireNonWhitespace(agentcommon.EnvKeyAPIBaseURL, common.IdentityEnvVarPrefixApp); !ok {
		return envSvc.ErrMissing(agentcommon.EnvKeyAPIBaseURL)
	} else if clientID, ok := envSvc.RequireNonWhitespace(common.EnvKeyAzClientID, common.IdentityEnvVarPrefixAgent); !ok {
		return envSvc.ErrMissing(common.EnvKeyAzClientID)
	} else if tenantID, ok = envSvc.RequireNonWhitespace(common.EnvKeyAzTenantID, common.IdentityEnvVarPrefixAgent); !ok {
		return envSvc.ErrMissing(common.EnvKeyAzTenantID)
	} else if certPath, ok := envSvc.RequireAbsPath(common.EnvKeyAzClientCertPath, common.IdentityEnvVarPrefixAgent); !ok {
		return envSvc.ErrMissing(common.EnvKeyAzClientCertPath)
	} else if cert, key, err := parseCertificateKeyPair(certPath); err != nil {
		return err
	} else if cred, err = azidentity.NewClientCertificateCredential(tenantID, clientID, []*x509.Certificate{cert}, key, nil); err != nil {
		return err
	} else if apiAuthScope, ok := envSvc.RequireNonWhitespace(agentcommon.EnvKeyAPIAuthScope, common.IdentityEnvVarPrefixApp); !ok {
		return envSvc.ErrMissing(agentcommon.EnvKeyAPIAuthScope)
	} else if client, err = agentclient.NewClientWithResponses(baseUrl,
		agentclient.WithRequestEditorFn(common.ToAzTokenCredentialRequestEditorFn(cred, policy.TokenRequestOptions{
			Scopes: []string{apiAuthScope},
		}))); err != nil {
		return err
	}

	resp, err := client.GetAgentConfigServerWithResponse(c, base.NamespaceKindServicePrincipal, base.ID("me"))
	if err != nil {
		return err
	}

	agentConfigServer := resp.JSON200
	//	log.Debug().Any("value", agentConfigServer).Msgf("GetAgentConfigServer: %d", resp.StatusCode())

	if err := agentcommon.DockerPullImage(c, agentConfigServer.AzureACRImageRef, cred, tenantID); err != nil {
		return err
	}

	// pull docker
	return nil
}
