package managedapp

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentConfigServerDoc struct {
	AgentConfigDoc
	TLSCertificatePolicyID base.Identifier          `json:"tlsCertificatePolicyId"`
	TLSCertificateID       base.Identifier          `json:"tlsCertificateId"`
	JWTKeyCertPolicyID     base.DocFullIdentifier   `json:"jwtKeyCertPolicyId"`
	JWTKeyCertIDs          []base.DocFullIdentifier `json:"jwtKeyCertIds"`
}

func (d *AgentConfigServerDoc) init(nsKind base.NamespaceKind, nsIdentifier base.Identifier) {
	d.AgentConfigDoc.init(nsKind, nsIdentifier, base.AgentConfigNameServer)
}

func (d *AgentConfigServerDoc) PopulateModel(m *AgentConfigServer) {
	if d == nil || m == nil {
		return
	}
	d.AgentConfigDoc.PopulateModel(&m.AgentConfig)
	m.TlsCertificatePolicyId = d.TLSCertificatePolicyID
	m.TlsCertificateId = d.TLSCertificateID
	m.JwtKeyCertPolicyId = d.JWTKeyCertPolicyID
	m.JWTKeyCertIDs = d.JWTKeyCertIDs
}

func apiPutAgentConfigServer(c ctx.RequestContext, param *AgentConfigServerParameters) error {

	nsCtx := ns.GetNSContext(c)

	doc := &AgentConfigServerDoc{}
	doc.init(nsCtx.Kind(), nsCtx.Identifier())

	doc.TLSCertificatePolicyID = param.TlsCertificatePolicyId

	certIds, err := cert.QueryLatestCertificateIdsIssuedByPolicy(c,
		base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, param.TlsCertificatePolicyId), 1)
	if err != nil {
		return err
	}
	if len(certIds) == 0 {
		return fmt.Errorf("%w: no certificate issued by policy %s", base.ErrResponseStatusBadRequest, param.TlsCertificatePolicyId.String())
	}
	doc.TLSCertificateID = certIds[0]

	err = base.GetAzCosmosCRUDService(c).Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	doc.JWTKeyCertPolicyID = param.JwtKeyCertPolicyId
	jwtKeyCertIdentifiers, err := cert.QueryLatestCertificateIdsIssuedByPolicy(c,
		doc.JWTKeyCertPolicyID, 2)
	if err != nil {
		return err
	}
	doc.JWTKeyCertIDs = utils.MapSlice(jwtKeyCertIdentifiers, func(id base.Identifier) base.DocFullIdentifier {
		return base.NewDocFullIdentifier(doc.JWTKeyCertPolicyID.NamespaceKind(), doc.JWTKeyCertPolicyID.NamespaceIdentifier(), base.ResourceKindCert, id)
	})

	m := &AgentConfigServer{}
	doc.PopulateModel(m)
	return c.JSON(200, m)
}
