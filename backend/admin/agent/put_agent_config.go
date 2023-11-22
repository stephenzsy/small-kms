package agentadmin

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// PutAgentConfig implements admin.ServerInterface.
func (*AgentAdminServer) PutAgentConfig(ec echo.Context, namespaceId string, configName agentmodels.AgentConfigName) error {

	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	params := new(agentmodels.CreateAgentConfigRequest)
	if err := c.Bind(params); err != nil {
		return err
	}

	if err := createBundleDocIfNotExist(c, namespaceId); err != nil {
		return err
	}

	switch configName {
	case agentmodels.AgentConfigNameIdentity:
		return putAgentConfigIdentity(c, namespaceId, params)
	}
	return base.ErrResponseStatusNotFound

}

func bundleDocIdentifier(nsID string) resdoc.DocIdentifier {
	return resdoc.NewDocIdentifier(models.NamespaceProviderServicePrincipal, nsID, models.ResourceProviderAgentConfig, "bundle")
}

func createBundleDocIfNotExist(c ctx.RequestContext, nsID string) error {
	doc := &AgentConfigBundleDoc{}
	docSvc := resdoc.GetDocService(c)
	docIdentifier := bundleDocIdentifier(nsID)
	if err := docSvc.Read(c, docIdentifier, doc, nil); err != nil {
		if !errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return err
		}
	} else {
		return nil
	}
	doc.PartitionKey = docIdentifier.PartitionKey
	doc.ID = docIdentifier.ID
	doc.Items = make(map[agentmodels.AgentConfigName]*AgentConfigBundleDocItem)

	_, err := docSvc.Upsert(c, doc, nil)
	return err
}
