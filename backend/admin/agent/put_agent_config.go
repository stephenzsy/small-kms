package agentadmin

import (
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

// PutAgentConfig implements ServerInterface.
func (*AgentAdminServer) PutAgentConfig(ec echo.Context, namespaceId string) error {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	params := new(agentmodels.AgentConfigFields)
	if err := c.Bind(params); err != nil {
		return err
	}

	doc := &AgentConfigDoc{
		ResourceDoc: resdoc.ResourceDoc{
			PartitionKey: resdoc.PartitionKey{
				NamespaceID:       namespaceId,
				NamespaceProvider: models.NamespaceProviderServicePrincipal,
				ResourceProvider:  models.ResourceProviderAgentConfig,
			},
			ID: ".default",
		},

		KeyCrendentialsCertificatePolicyID: params.KeyCredentialsCertificatePolicyId,
	}

	// TODO: verify credential policy exists

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}
