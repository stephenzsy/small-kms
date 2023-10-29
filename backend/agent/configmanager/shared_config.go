package cm

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	agentclient "github.com/stephenzsy/small-kms/backend/agent-client"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type sharedConfig struct {
	client             agentclient.ClientWithResponsesInterface
	azSecretsClient    *azsecrets.Client
	configDir          string
	versionedConfigDir string
}

func (sc *sharedConfig) AgentClient() agentclient.ClientWithResponsesInterface {
	return sc.client
}

func (sc *sharedConfig) AzSecretesClient() *azsecrets.Client {
	return sc.azSecretsClient
}

func (sc *sharedConfig) init(
	buildID string,
	client agentclient.ClientWithResponsesInterface,
	azSecretsClient *azsecrets.Client,
	configDir string,
) error {
	sc.client = client
	sc.azSecretsClient = azSecretsClient

	// ensure config dir
	if _, err := os.Stat(configDir); err != nil {
		return err
	}
	sc.configDir = configDir
	sc.versionedConfigDir = filepath.Join(configDir, "versioned")
	if err := os.MkdirAll(sc.versionedConfigDir, 0700); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	return nil
}

var meNamespaceIdIdentifier = shared.StringIdentifier("me")

const (
	TaskNameLoad     = "load" // load from file
	TaskNameFetch    = "fetch"
	TaskNameActivate = "activate"
	TaskNameConfirm  = "confirm"
)
