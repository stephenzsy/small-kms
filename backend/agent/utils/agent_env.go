package agentutils

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	agentclient "github.com/stephenzsy/small-kms/backend/agent/client"
	"github.com/stephenzsy/small-kms/backend/common"
)

type AgentEnv struct {
	tenantID                   string
	clientID                   string
	clientCertPath             string
	apiBaseURL                 string
	apiAuthScope               string
	azKeyVaultResourceEndpoint string

	certCred        *azidentity.ClientCertificateCredential
	agentClient     *agentclient.ClientWithResponses
	azSecretsClient *azsecrets.Client
}

func NewAgentEnv() (env AgentEnv, err error) {
	env.tenantID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, common.IdentityEnvVarNameAzTenantID, "")
	env.clientID = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, common.IdentityEnvVarNameAzClientID, "")
	env.clientCertPath = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, common.IdentityEnvVarNameAzClientCertPath, "")
	env.apiBaseURL = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, "API_BASE_URL", "")
	env.apiAuthScope = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, "API_AUTH_SCOPE", "")
	env.azKeyVaultResourceEndpoint = common.LookupPrefixedEnvWithDefault(common.IdentityEnvVarPrefixAgent, "AZURE_KEYVAULT_RESOURCEENDPOINT", "")
	return env, nil
}

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

func (ae *AgentEnv) CertCred() (*azidentity.ClientCertificateCredential, error) {
	if ae.certCred != nil {
		return ae.certCred, nil
	}
	if ae.tenantID == "" {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, "AZURE_TENANT_ID")
	}
	if ae.clientID == "" {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, "AZURE_CLIENT_ID")
	}
	if ae.clientCertPath == "" {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, "AZURE_CLIENT_CERTIFICATE_PATH")
	}
	if cert, key, err := parseCertificateKeyPair(ae.clientCertPath); err != nil {
		return nil, err
	} else if cred, err := azidentity.NewClientCertificateCredential(
		ae.tenantID,
		ae.clientID,
		[]*x509.Certificate{cert}, key, nil); err != nil {
		return nil, err
	} else {
		ae.certCred = cred
		return ae.certCred, nil
	}
}

func (ae *AgentEnv) AgentClient() (*agentclient.ClientWithResponses, error) {
	if ae.agentClient != nil {
		return ae.agentClient, nil
	}
	certCred, err := ae.CertCred()
	if err != nil {
		return nil, err
	}
	if ae.apiBaseURL == "" {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, "API_BASE_URL")
	}
	if ae.apiAuthScope == "" {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, "API_AUTH_SCOPE")
	}
	if agentClient, err := agentclient.NewClientWithResponses(
		ae.apiBaseURL,
		agentclient.WithRequestEditorFn(common.ToAzTokenCredentialRequestEditorFn(
			certCred, policy.TokenRequestOptions{
				Scopes: []string{ae.apiAuthScope},
			}))); err != nil {
		return nil, err
	} else {
		ae.agentClient = agentClient
		return ae.agentClient, nil
	}
}

func (ae *AgentEnv) AzSecretsClient() (*azsecrets.Client, error) {
	if ae.azSecretsClient != nil {
		return ae.azSecretsClient, nil
	}
	creds, err := ae.CertCred()
	if err != nil {
		return nil, err
	}
	if ae.azKeyVaultResourceEndpoint == "" {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, "AZURE_KEYVAULT_RESOURCE_ENDPOINT")
	}
	if azSecretsClient, err := azsecrets.NewClient(ae.azKeyVaultResourceEndpoint, creds, nil); err != nil {
		return nil, err
	} else {
		ae.azSecretsClient = azSecretsClient
		return ae.azSecretsClient, nil
	}
}
