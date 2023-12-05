package agentadmin

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/key/v2"
	"github.com/stephenzsy/small-kms/backend/models"
	agentmodels "github.com/stephenzsy/small-kms/backend/models/agent"
	"github.com/stephenzsy/small-kms/backend/resdoc"
)

type agentConfigDocEndpoint struct {
	AgentConfigDoc
	TLSCertificatePolicyID   string   `json:"tlsCertificatePolicyId"`
	TLSCertificateAutoEnroll bool     `json:"tlsCertificateAutoEnroll"`
	TLSCertificateID         string   `json:"tlsCertificateId,omitempty"`
	JWTVerifyKeyPolicyID     string   `json:"jwtVerifyKeyPolicyId"`
	JWTVerifyKeyIDs          []string `json:"jwtVerifyKeyIds"`
}

func (d *agentConfigDocEndpoint) ToModel() (m agentmodels.AgentConfigEndpoint) {
	m.Name = agentmodels.AgentConfigNameEndpoint
	m.Updated = d.Timestamp.Time
	m.Version = hex.EncodeToString(d.Version)
	m.TlsCertificatePolicyId = d.TLSCertificatePolicyID
	m.JwtVerifyKeyPolicyId = d.JWTVerifyKeyPolicyID
	m.JwtVerifyKeyIds = d.JWTVerifyKeyIDs
	m.TLSCertificateAutoEnroll = d.TLSCertificateAutoEnroll
	m.TLSCertificateID = d.TLSCertificateID
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
		TLSCertificatePolicyID:   req.TlsCertificatePolicyId,
		JWTVerifyKeyPolicyID:     req.JwtVerifyKeyPolicyId,
		TLSCertificateAutoEnroll: req.TLSCertificateAutoEnroll,
	}

	versiond := md5.New()

	certPolicy, err := cert.GetCertificatePolicyInternal(c, models.NamespaceProviderServicePrincipal, namespaceId, req.TlsCertificatePolicyId)
	if err != nil {
		return err
	}
	versiond.Write([]byte(doc.TLSCertificatePolicyID))
	versiond.Write(certPolicy.Version)

	if !doc.TLSCertificateAutoEnroll {
		doc.TLSCertificateID, err = certPolicy.GetLatestIssuedCertificateID(c)
		if err != nil {
			return err
		}
		versiond.Write([]byte(doc.TLSCertificateID))
	}

	keyPolicy, err := key.GetKeyPolicyInternal(c, models.NamespaceProviderServicePrincipal, namespaceId, req.JwtVerifyKeyPolicyId)
	if err != nil {
		return err
	}
	versiond.Write([]byte(doc.JWTVerifyKeyPolicyID))
	versiond.Write(keyPolicy.Version)

	doc.JWTVerifyKeyIDs, err = key.ListLatestActiveKeysByPolicyInternal(c, models.NamespaceProviderServicePrincipal, namespaceId,
		keyPolicy.Identifier())
	if err != nil {
		return err
	}
	for _, keyId := range doc.JWTVerifyKeyIDs {
		versiond.Write([]byte(keyId))
	}

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

	doc, err := getAgentConfigEndpointInternal(c, namespaceId)
	if err != nil {
		return err
	}

	return c.JSON(200, doc.ToModel())
}

func getAgentConfigEndpointInternal(c context.Context, namespaceId string) (*agentConfigDocEndpoint, error) {
	docSvc := resdoc.GetDocService(c)
	doc := &agentConfigDocEndpoint{}
	err := docSvc.Read(c, resdoc.NewDocIdentifier(models.NamespaceProviderServicePrincipal,
		namespaceId, models.ResourceProviderAgentConfig, string(agentmodels.AgentConfigNameEndpoint)), doc, nil)
	if err != nil {
		if errors.Is(err, resdoc.ErrAzCosmosDocNotFound) {
			return doc, base.ErrResponseStatusNotFound
		}
	}
	return doc, err
}
