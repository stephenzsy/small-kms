package agentconfig

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type AgentConfigDocument interface {
	kmsdoc.KmsDocument
	toModel(isAdmin bool) *shared.AgentConfiguration
}

type AgentConfigDoc struct {
	kmsdoc.BaseDoc

	Name        shared.AgentConfigName   `json:"name"`          // for index only
	BaseVersion kmsdoc.HexStringStroable `json:"configVersion"` // version of the config from put request
	Version     kmsdoc.HexStringStroable `json:"version"`       // version of the config after evaluation
}

type docConfigurator[D AgentConfigDocument] struct {
	preparePut func(
		c context.Context,
		nsID shared.NamespaceIdentifier,
		params shared.AgentConfigurationParameters) (D, error)
	eval func(
		c context.Context,
		doc D) (*azcosmos.PatchOperations, error)
	readDoc func(
		c context.Context,
		nsID shared.NamespaceIdentifier) (D, error)
}

func NewConfigDocLocator(nsID shared.NamespaceIdentifier, configName shared.AgentConfigName) shared.ResourceLocator {
	return shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindAgentConfig,
		shared.StringIdentifier(configName)))
}

func (d *AgentConfigDoc) initLocator(nsID shared.NamespaceIdentifier, configName shared.AgentConfigName) {
	d.NamespaceID = nsID
	d.ID = shared.NewResourceIdentifier(shared.ResourceKindAgentConfig, shared.StringIdentifier(configName))
	d.Name = configName
}

func getTimeRefreshToken(t time.Time) string {
	return base64.RawURLEncoding.EncodeToString([]byte(t.Format(time.RFC3339)))
}
