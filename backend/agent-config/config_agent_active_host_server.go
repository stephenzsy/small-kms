package agentconfig

import (
	"crypto/md5"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type AgentActiveServerDoc struct {
	AgentConfigDoc

	ControllerImageRef string `json:"controllerImageRef"`
}

func (d *AgentActiveServerDoc) toModel() *models.AgentConfiguration {
	if d == nil {
		return nil
	}
	m := models.AgentConfiguration{
		Version: d.Version.HexString(),
	}
	m.Config.FromAgentConfigurationAgentActiveHostBootstrap(models.AgentConfigurationAgentActiveHostBootstrap{
		ControllerContainer: models.AgentConfigurationActiveHostControllerContainer{
			ImageRefStr: d.ControllerImageRef,
		},
		Name: models.AgentConfigNameActiveServer,
	})
	return &m
}

func handleGetAgentActiveServer(c RequestContext, nsID shared.NamespaceIdentifier) (*models.AgentConfiguration, error) {
	locator := shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindAgentConfig, common.StringIdentifier(models.AgentConfigNameActiveHostBootstrap)))
	doc := AgentActiveServerDoc{}
	err := kmsdoc.Read(c, locator, &doc)
	return doc.toModel(), err
}
func handlePutAgentActiveServer(c RequestContext, nsID models.NamespaceID, params models.AgentConfigurationParameters) (*models.AgentConfiguration, error) {
	locator := shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindAgentConfig, common.StringIdentifier(models.AgentConfigNameActiveHostBootstrap)))
	p, err := params.AsAgentConfigurationAgentActiveHostBootstrap()
	if err != nil {
		return nil, fmt.Errorf("%w:invalid input", common.ErrStatusBadRequest)
	}
	digest := md5.New()
	doc := AgentActiveServerDoc{
		AgentConfigDoc: AgentConfigDoc{
			BaseDoc: kmsdoc.BaseDoc{
				ID:          locator.GetID(),
				NamespaceID: locator.GetNamespaceID(),
			},
			Name: string(models.AgentConfigNameActiveHostBootstrap),
		},
		ControllerImageRef: p.ControllerContainer.ImageRefStr,
	}
	digest.Write([]byte(doc.ControllerImageRef))
	hash := digest.Sum(nil)
	doc.Version = hash[:]
	err = kmsdoc.Upsert(c, &doc)
	return doc.toModel(), err
}
