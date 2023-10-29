package managedapp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert/v2"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type AgentConfigServerDoc struct {
	AgentConfigDoc
	GlobalKeyVaultEndpoint string                   `json:"globalKeyVaultEndpoint"`
	GlobalACRImageRef      string                   `json:"globalAcrImageRef"`
	TLSCertificatePolicyID base.Identifier          `json:"tlsCertificatePolicyId"`
	TLSCertificateID       base.Identifier          `json:"tlsCertificateId"`
	JWTKeyCertPolicyID     base.DocFullIdentifier   `json:"jwtKeyCertPolicyId"`
	JWTKeyCertIDs          []base.DocFullIdentifier `json:"jwtKeyCertIds"`
}

func (d *AgentConfigServerDoc) init(nsKind base.NamespaceKind, nsIdentifier base.Identifier) {
	d.AgentConfigDoc.init(nsKind, nsIdentifier, base.AgentConfigNameServer)
}

const envVarMessage = "For security reasons, these variables should be set manualy after review in the agent config server. Server must not auto update these variables."

func (d *AgentConfigServerDoc) populateModel(m *AgentConfigServer) {
	if d == nil || m == nil {
		return
	}
	d.AgentConfigDoc.PopulateModel(&m.AgentConfig)
	m.Env = AgentConfigServerEnv{
		Message:                             envVarMessage,
		EnvVarAzureKeyVaultResourceEndpoint: d.GlobalKeyVaultEndpoint,
		EnvVarAzureContainerRegistryImageRepository: strings.Split(d.GlobalACRImageRef, ":")[0],
	}
	m.TlsCertificatePolicyId = d.TLSCertificatePolicyID
	m.TlsCertificateId = d.TLSCertificateID
	m.JwtKeyCertPolicyId = d.JWTKeyCertPolicyID
	m.JWTKeyCertIDs = d.JWTKeyCertIDs
	m.RefreshAfter = time.Now().Add(24 * time.Hour).UTC()
	m.AzureACRImageRef = d.GlobalACRImageRef
}

func apiGetAgentConfigServer(c ctx.RequestContext) error {
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
	doc.populateModel(m)
	return c.JSON(200, m)
}

func (s *server) apiPutAgentConfigServer(c ctx.RequestContext, param *AgentConfigServerFields) error {

	nsCtx := ns.GetNSContext(c)

	doc := &AgentConfigServerDoc{}
	doc.init(nsCtx.Kind(), nsCtx.Identifier())

	docSvc := base.GetAzCosmosCRUDService(c)

	digest := md5.New()
	switch nsCtx.Kind() {
	case base.NamespaceKindSystem:
		if nsCtx.Identifier() != base.StringIdentifier("default") {
			return fmt.Errorf("%w: only default system namespace is supported", base.ErrResponseStatusBadRequest)
		}
		doc.GlobalKeyVaultEndpoint = s.GetAzKeyVaultEndpoint()
		digest.Write([]byte(doc.GlobalKeyVaultEndpoint))

		doc.GlobalACRImageRef = param.AzureACRImageRef
		digest.Write([]byte(doc.GlobalACRImageRef))
	case base.NamespaceKindServicePrincipal:

		globalDoc := &AgentConfigServerDoc{}
		if err := docSvc.Read(c,
			base.NewDocFullIdentifier(base.NamespaceKindSystem,
				base.StringIdentifier("default"), base.ResourceKindNamespaceConfig, base.StringIdentifier(string(base.AgentConfigNameServer))), globalDoc, nil); err != nil {
			if errors.Is(err, base.ErrAzCosmosDocNotFound) {
				return fmt.Errorf("%w: %s", base.ErrResponseStatusNotFound, base.AgentConfigNameServer)
			}
			return err
		}

		doc.GlobalKeyVaultEndpoint = globalDoc.GlobalKeyVaultEndpoint
		doc.GlobalACRImageRef = globalDoc.GlobalACRImageRef
		globalDocVersionBytes, err := hex.DecodeString(globalDoc.Version)
		if err != nil {
			return err
		}
		digest.Write(globalDocVersionBytes)

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
		digest.Write([]byte(doc.TLSCertificateID.String()))

		doc.JWTKeyCertPolicyID = param.JwtKeyCertPolicyId
		jwtKeyCertIdentifiers, err := cert.QueryLatestCertificateIdsIssuedByPolicy(c,
			doc.JWTKeyCertPolicyID, 2)
		if err != nil {
			return err
		}
		doc.JWTKeyCertIDs = utils.MapSlice(jwtKeyCertIdentifiers, func(id base.Identifier) base.DocFullIdentifier {
			fullIdentifier := base.NewDocFullIdentifier(doc.JWTKeyCertPolicyID.NamespaceKind(), doc.JWTKeyCertPolicyID.NamespaceIdentifier(), base.ResourceKindCert, id)
			digest.Write([]byte(fullIdentifier.String()))
			return fullIdentifier
		})
	default:
		return fmt.Errorf("%w: unsupported namespace kind: %s", base.ErrResponseStatusBadRequest, nsCtx.Kind())
	}
	doc.Version = hex.EncodeToString(digest.Sum(nil))

	err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
	}

	m := &AgentConfigServer{}
	doc.populateModel(m)
	return c.JSON(200, m)
}
