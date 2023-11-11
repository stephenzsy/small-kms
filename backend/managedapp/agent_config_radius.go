package managedapp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/labstack/echo/v4"
	frconfig "github.com/stephenzsy/small-kms/backend/agent/freeradiusconfig"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/secret"
)

// PatchAgentConfigRadius implements ServerInterface.
func (s *server) PatchAgentConfigRadius(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	param := AgentConfigRadiusFields{}
	if err := c.Bind(&param); err != nil {
		return err
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceId)
	doc, err := s.patchAgentConfigRadius(c, namespaceKind, namespaceId, param, namespaceKind == base.NamespaceKindServicePrincipal)
	if err != nil {
		return err
	}
	m := &AgentConfigRadius{}
	doc.populateModel(m)
	return c.JSON(http.StatusOK, m)
}

func (s *server) patchAgentConfigRadius(c ctx.RequestContext, namespaceKind base.NamespaceKind, namespaceId base.ID, p AgentConfigRadiusFields, assignRoles bool) (*AgentConfigRadiusDoc, error) {
	docSvc := base.GetAzCosmosCRUDService(c)
	doc, err := readAgentConfigRadiusDoc(c)
	digest := md5.New()
	if err != nil {
		if !errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, err
		}
		// create a new
		doc = &AgentConfigRadiusDoc{}
		doc.init(namespaceKind, namespaceId)
	}

	if p.Container != nil {
		doc.Container = *p.Container
	}
	doc.Container.digest(digest)

	if p.Clients != nil {
		doc.Clients = make([]frconfig.RadiusClientConfig, len(p.Clients))
		for i, pClient := range p.Clients {
			docClient := frconfig.RadiusClientConfig{
				Ipaddr:         pClient.Ipaddr,
				Name:           pClient.Name,
				SecretPolicyId: pClient.SecretPolicyId,
			}
			if pClient.SecretPolicyId != "" {
				queried, err := secret.QueryLatestSecretIDIssuedByPolicy(c, base.NewDocLocator(namespaceKind, namespaceId, base.ResourceKindSecretPolicy, pClient.SecretPolicyId), 1)
				if err != nil {
					return nil, err
				}
				if len(queried) > 0 {
					docClient.SecretId = queried[0]
				}
			}
			doc.Clients[i] = docClient
		}
	}
	for _, client := range doc.Clients {
		digest.Write([]byte(client.Name))
		digest.Write([]byte(client.Ipaddr))
		digest.Write([]byte(client.SecretId))
	}

	if p.EapTls != nil {
		doc.EapTls = *p.EapTls
		queried, err := cert.QueryLatestCertificateIdsIssuedByPolicy(
			c,
			base.NewDocLocator(namespaceKind, namespaceId, base.ResourceKindCertPolicy, doc.EapTls.CertPolicyId), 1)
		if err != nil {
			return nil, err
		}
		if len(queried) <= 0 {
			return nil, errors.New("no certificate issued by policy")
		}
		doc.EapTls.CertId = queried[0]
	}
	digest.Write([]byte(doc.EapTls.CertId))

	if p.DebugMode != nil {
		doc.DebugMode = p.DebugMode
		if *p.DebugMode {
			digest.Write([]byte("debug"))
		}
	}

	doc.Version = hex.EncodeToString(digest.Sum(nil))

	err = docSvc.Upsert(c, doc, &azcosmos.ItemOptions{IfMatchEtag: doc.ETag})

	if assignRoles {
		s.assignAgentRadiusRoles(c, namespaceId.UUID(), doc)
	}
	return doc, err
}
