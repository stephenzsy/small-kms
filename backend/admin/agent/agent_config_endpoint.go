package agentadmin

import (
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type agentConfigDocEndpoint struct {
	AgentConfigDoc
	TLSCertificatePolicyID string `json:"tlsCertificatePolicyId"`
}

func (d *agentConfigDocEndpoint) ToModel() (m agentmodels.AgentConfigEndpoint) {
	m.Name = agentmodels.AgentConfigNameEndpoint
	m.Updated = d.Timestamp.Time
	m.Version = hex.EncodeToString(d.Version)
	m.TlsCertificatePolicyId = d.TLSCertificatePolicyID
	return m
}

func putAgentConfigEndpoint(c ctx.RequestContext, namespaceId string, param *agentmodels.CreateAgentConfigRequest) error {

	req, err := param.AsAgentConfigEndpointFields()
	if err != nil {
		return err
	}

	doc := &agentConfigDocEndpoint{
		AgentConfigDoc: AgentConfigDoc{
			ResourceDoc: resdoc.ResourceDoc{
				PartitionKey: resdoc.PartitionKey{
					NamespaceID:       namespaceId,
					NamespaceProvider: models.NamespaceProviderServicePrincipal,
					ResourceProvider:  models.ResourceProviderAgentConfig,
				},
				ID: string(agentmodels.AgentConfigNameEndpoint),
			},
		},
		TLSCertificatePolicyID: req.TlsCertificatePolicyId,
	}

	versiond := md5.New()

	policy, err := cert.GetCertificatePolicyInternal(c, models.NamespaceProviderServicePrincipal, namespaceId, req.TlsCertificatePolicyId)
	if err != nil {
		return err
	}
	versiond.Write([]byte(doc.TLSCertificatePolicyID))
	versiond.Write(policy.Version)

	doc.Version = versiond.Sum(nil)
	docSvc := resdoc.GetDocService(c)
	resp, err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	ops := azcosmos.PatchOperations{}
	ops.AppendSet("/items/endpoint", &AgentConfigBundleDocItem{
		Updated: doc.Timestamp.Time,
		Version: doc.Version,
	})
	_, err = docSvc.PatchByIdentifier(c, bundleDocIdentifier(namespaceId), ops, nil)
	if err != nil {
		return err
	}

	return c.JSON(resp.RawResponse.StatusCode, doc.ToModel())
}

func getAgentConfigEndpoint(c ctx.RequestContext, namespaceId string) error {

	docSvc := resdoc.GetDocService(c)
	doc := &agentConfigDocEndpoint{}
	if err := docSvc.Read(c, resdoc.NewDocIdentifier(models.NamespaceProviderServicePrincipal,
		namespaceId, models.ResourceProviderAgentConfig, string(agentmodels.AgentConfigNameEndpoint)), doc, nil); err != nil {
		err = resdoc.HandleAzCosmosError(err)
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return base.ErrResponseStatusNotFound
		}
		return err
	}

	return c.JSON(200, doc.ToModel())
}
