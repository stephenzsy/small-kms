package agentcommon

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
	"github.com/stephenzsy/small-kms/backend/managedapp"
)

type AgentEnv struct {
	common.EnvService
	mode managedapp.AgentMode

	certCred        *azidentity.ClientCertificateCredential
	agentClient     *agentclient.ClientWithResponses
	azSecretsClient *azsecrets.Client
}

func NewAgentEnv(envService common.EnvService, mode managedapp.AgentMode) (env *AgentEnv, err error) {
	env = &AgentEnv{
		EnvService: envService,
		mode:       mode,
	}
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
	if tenantID, ok := ae.RequireNonWhitespace(common.EnvKeyAzTenantID, common.IdentityEnvVarPrefixAgent); !ok {
		return nil, ae.ErrMissing(common.EnvKeyAzTenantID)
	} else if clientID, ok := ae.RequireNonWhitespace(common.EnvKeyAzClientID, common.IdentityEnvVarPrefixAgent); !ok {
		return nil, ae.ErrMissing(common.EnvKeyAzClientID)
	} else if clientCertPath, ok := ae.RequireAbsPath(common.EnvKeyAzClientCertPath, common.IdentityEnvVarPrefixAgent); !ok {
		return nil, ae.ErrMissing(common.EnvKeyAzClientCertPath)
	} else if cert, key, err := parseCertificateKeyPair(clientCertPath); err != nil {
		return nil, err
	} else if cred, err := azidentity.NewClientCertificateCredential(
		tenantID,
		clientID,
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
	if apiBaseURL, ok := ae.RequireNonWhitespace(EnvKeyAPIBaseURL, common.IdentityEnvVarPrefixApp); !ok {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, EnvKeyAPIBaseURL)
	} else if apiAuthScope, ok := ae.RequireNonWhitespace(EnvKeyAPIAuthScope, common.IdentityEnvVarPrefixApp); !ok {
		return nil, fmt.Errorf("%w: %s", common.ErrMissingEnvVar, EnvKeyAPIAuthScope)
	} else if agentClient, err := agentclient.NewClientWithResponses(
		apiBaseURL,
		agentclient.WithRequestEditorFn(common.ToAzTokenCredentialRequestEditorFn(
			certCred, policy.TokenRequestOptions{
				Scopes: []string{apiAuthScope},
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
	if endpoint, ok := ae.RequireNonWhitespace(common.EnvKeyAzKeyvaultResourceEndpoint, common.IdentityEnvVarPrefixAgent); !ok {
		return nil, ae.ErrMissing(common.EnvKeyAzKeyvaultResourceEndpoint)
	} else if azSecretsClient, err := azsecrets.NewClient(endpoint, creds, nil); err != nil {
		return nil, err
	} else {
		ae.azSecretsClient = azSecretsClient
		return ae.azSecretsClient, nil
	}
}
