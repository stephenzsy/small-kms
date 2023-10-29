package managedapp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

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

func (d *AgentConfigServerDoc) populateModel(s *server, m *AgentConfigServer) {
	if d == nil || m == nil {
		return
	}
	d.AgentConfigDoc.PopulateModel(&m.AgentConfig)
	m.Env = AgentConfigServerEnv{
		AZUREKEYVAULTRESOURCEENDPOINT: s.GetAzKeyVaultEndpoint(),
	}
	m.TlsCertificatePolicyId = d.TLSCertificatePolicyID
	m.TlsCertificateId = d.TLSCertificateID
	m.JwtKeyCertPolicyId = d.JWTKeyCertPolicyID
	m.JWTKeyCertIDs = d.JWTKeyCertIDs
	m.RefreshAfter = time.Now().Add(24 * time.Hour).UTC()
	m.ImageTag = s.GetBuildID()
}

func (s *server) apiGetAgentConfigServer(c ctx.RequestContext) error {
	nsCtx := ns.GetNSContext(c)

	doc := &AgentConfigServerDoc{}

	if err := base.GetAzCosmosCRUDService(c).Read(c, base.NewDocFullIdentifier(nsCtx.Kind(),
		nsCtx.Identifier(), base.ResourceKindNamespaceConfig, base.StringIdentifier(string(base.AgentConfigNameServer))), doc, nil); err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: %s", base.ErrResponseStatusNotFound, base.AgentConfigNameServer)
		}
		return err
	}

	m := &AgentConfigServer{}
	doc.populateModel(s, m)
	return c.JSON(200, m)
}

func (s *server) apiPutAgentConfigServer(c ctx.RequestContext, param *AgentConfigServerParameters) error {

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

	doc.JWTKeyCertPolicyID = param.JwtKeyCertPolicyId
	jwtKeyCertIdentifiers, err := cert.QueryLatestCertificateIdsIssuedByPolicy(c,
		doc.JWTKeyCertPolicyID, 2)
	if err != nil {
		return err
	}
	doc.JWTKeyCertIDs = utils.MapSlice(jwtKeyCertIdentifiers, func(id base.Identifier) base.DocFullIdentifier {
		return base.NewDocFullIdentifier(doc.JWTKeyCertPolicyID.NamespaceKind(), doc.JWTKeyCertPolicyID.NamespaceIdentifier(), base.ResourceKindCert, id)
	})

	digest := md5.New()
	digest.Write([]byte(doc.TLSCertificateID.String()))
	for _, certID := range doc.JWTKeyCertIDs {
		digest.Write([]byte(certID.String()))
	}
	doc.Version = hex.EncodeToString(digest.Sum(nil))

	err = base.GetAzCosmosCRUDService(c).Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	m := &AgentConfigServer{}
	doc.populateModel(s, m)
	return c.JSON(200, m)
}
