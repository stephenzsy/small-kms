package agentconfig

import (
	"context"
	"crypto/md5"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type AgentActiveHostBootstrapDoc struct {
	AgentConfigDoc

	ControllerImageRef string `json:"controllerImageRef"`
}

func (d *AgentActiveHostBootstrapDoc) toModel(bool) *shared.AgentConfiguration {
	if d == nil {
		return nil
	}
	m := shared.AgentConfiguration{
		Version: d.Version.HexString(),
	}
	m.Config.FromAgentConfigurationAgentActiveHostBootstrap(shared.AgentConfigurationAgentActiveHostBootstrap{
		ControllerContainer: shared.AgentConfigurationActiveHostControllerContainer{
			ImageRefStr: d.ControllerImageRef,
		},
		Name: shared.AgentConfigNameActiveHostBootstrap,
	})
	return &m
}

var _ AgentConfigDocument = (*AgentActiveHostBootstrapDoc)(nil)

func newAgentActiveHostBootStrapConfigurator() *docConfigurator[AgentConfigDocument] {
	return &docConfigurator[AgentConfigDocument]{
		preparePut: func(
			c context.Context,
			nsID shared.NamespaceIdentifier, params shared.AgentConfigurationParameters) (AgentConfigDocument, error) {
			p, err := params.AsAgentConfigurationAgentActiveHostBootstrap()
			if err != nil {
				return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
			}

			d := AgentActiveHostBootstrapDoc{
				ControllerImageRef: p.ControllerContainer.ImageRefStr,
			}
			d.initLocator(nsID, shared.AgentConfigNameActiveHostBootstrap)
			digester := md5.New()
			digester.Write([]byte(d.ControllerImageRef))
			d.BaseVersion = digester.Sum(nil)
			d.Version = d.BaseVersion
			return &d, nil
		},
		eval: func(_ context.Context, doc AgentConfigDocument) (*azcosmos.PatchOperations, error) {
			return nil, nil
		},
		readDoc: func(c context.Context, nsID shared.NamespaceIdentifier) (AgentConfigDocument, error) {
			d := AgentActiveHostBootstrapDoc{}
			err := kmsdoc.Read(c, NewConfigDocLocator(nsID, shared.AgentConfigNameActiveHostBootstrap), &d)
			return &d, err
		},
	}
}
