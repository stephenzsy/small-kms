package bootstrap

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	agentclient "github.com/stephenzsy/small-kms/backend/agent/client"
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
	if baseUrl := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, "API_BASE_URL", ""); baseUrl == "" {
		return errors.New("missing API_URL_BASE")
	} else if clientID := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzClientID, ""); clientID == "" {
		return errors.New("missing APP_AZURE_CLIENT_ID")
	} else if tenantID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzTenantID, ""); tenantID == "" {
		return errors.New("missing APP_AZURE_TENANT_ID")
	} else if certPath := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, common.IdentityEnvVarNameAzClientCertPath, ""); certPath == "" {
		return errors.New("missing APP_AZURE_CLIENT_CERTIFICATE_PATH")
	} else if cert, key, err := parseCertificateKeyPair(certPath); err != nil {
		return err
	} else if cred, err = azidentity.NewClientCertificateCredential(tenantID, clientID, []*x509.Certificate{cert}, key, nil); err != nil {
		return err
	} else if apiAuthScope := common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixApp, "API_AUTH_SCOPE", ""); apiAuthScope == "" {
		return errors.New("missing APP_API_AUTH_SCOPE")
	} else if client, err = agentclient.NewClientWithResponses(baseUrl,
		agentclient.WithRequestEditorFn(common.ToAzTokenCredentialRequestEditorFn(cred, policy.TokenRequestOptions{
			Scopes: []string{apiAuthScope},
		}))); err != nil {
		return err
	}

	resp, err := client.GetAgentConfigServerWithResponse(c, base.NamespaceKindServicePrincipal, base.StringIdentifier("me"))
	if err != nil {
		return err
	}

	agentConfigServer := resp.JSON200
	//	log.Debug().Any("value", agentConfigServer).Msgf("GetAgentConfigServer: %d", resp.StatusCode())

	if err := dockerPullImage(c, agentConfigServer.AzureACRImageRef, cred, tenantID); err != nil {
		return err
	}

	// pull docker
	return nil
}
