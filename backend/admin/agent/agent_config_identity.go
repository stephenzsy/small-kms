package agentadmin

import (
	"errors"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

func putAgentConfigIdentity(c ctx.RequestContext, namespaceId string, param *agentmodels.CreateAgentConfigRequest) error {

	req, err := param.AsAgentConfigIdentityFields()
	if err != nil {
		return err
	}

	doc := &AgentConfigDocIdentity{
		AgentConfigDoc: AgentConfigDoc{
			ResourceDoc: resdoc.ResourceDoc{
				PartitionKey: resdoc.PartitionKey{
					NamespaceID:       namespaceId,
					NamespaceProvider: models.NamespaceProviderServicePrincipal,
					ResourceProvider:  models.ResourceProviderAgentConfig,
				},
				ID: string(agentmodels.AgentConfigNameIdentity),
			},
		},
		KeyCredentialsCertificatePolicyID: req.KeyCredentialCertificatePolicyId,
	}

	policy, err := cert.GetCertificatePolicyInternal(c, models.NamespaceProviderServicePrincipal, namespaceId, req.KeyCredentialCertificatePolicyId)
	if err != nil {
		return err
	}
	doc.Version = policy.Version

	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}

func getAgentConfigIdentity(c ctx.RequestContext, namespaceId string) error {

	docSvc := resdoc.GetDocService(c)
	doc := &AgentConfigDocIdentity{}
	if err := docSvc.Read(c, resdoc.NewDocIdentifier(models.NamespaceProviderServicePrincipal,
		namespaceId, models.ResourceProviderAgentConfig, string(agentmodels.AgentConfigNameIdentity)), doc, nil); err != nil {
		err = resdoc.HandleAzCosmosError(err)
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return base.ErrResponseStatusNotFound
		}
		return err
	}

	return c.JSON(200, doc.ToModel())
}
