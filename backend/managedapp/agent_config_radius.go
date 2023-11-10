package managedapp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/agent/radius"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/secret"
)

// PutAgentConfigRadius implements ServerInterface.
func (s *server) PutAgentConfigRadius(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	param := AgentConfigRadiusFields{}
	if err := c.Bind(&param); err != nil {
		return err
	}

	c = ns.WithNSContext(c, nsKind, nsID)

	doc := new(AgentConfigRadiusDoc)
	doc.init(nsKind, nsID)
	docSvc := base.GetAzCosmosCRUDService(c)
	switch nsKind {
	case base.NamespaceKindSystem:
		if string(nsID) != "default" {
			return fmt.Errorf("%w: only default system namespace is supported", base.ErrResponseStatusBadRequest)
		}
		doc.GlobalRadiusServerACRImageRef = *param.AzureACRImageRef
		digest := md5.New()
		digest.Write([]byte(doc.GlobalRadiusServerACRImageRef))
		doc.Version = hex.EncodeToString(digest.Sum(nil))
		err := docSvc.Upsert(c, doc, nil)
		if err != nil {
			return err
		}
	case base.NamespaceKindServicePrincipal:
		var err error
		doc, err = s.patchAgentConfigRadius(c, nsKind, nsID, param)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%w: namespace kind %s is not supported", base.ErrResponseStatusBadRequest, nsKind)
	}

	m := &AgentConfigRadius{}
	doc.populateModel(m)
	return c.JSON(http.StatusOK, m)
}

// PatchAgentConfigRadius implements ServerInterface.
func (s *server) PatchAgentConfigRadius(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	if namespaceKind == base.NamespaceKindSystem {
		return fmt.Errorf("%w: patch to system config is not supported", base.ErrResponseStatusBadRequest)
	}
	param := AgentConfigRadiusFields{}
	if err := c.Bind(&param); err != nil {
		return err
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceId)
	doc, err := s.patchAgentConfigRadius(c, namespaceKind, namespaceId, param)
	if err != nil {
		return err
	}
	m := &AgentConfigRadius{}
	doc.populateModel(m)
	return c.JSON(http.StatusOK, m)
}

func (s *server) patchAgentConfigRadius(c ctx.RequestContext, namespaceKind base.NamespaceKind, namespaceId base.ID, p AgentConfigRadiusFields) (*AgentConfigRadiusDoc, error) {
	nsUUID, ok := namespaceId.AsUUID()
	if !ok {
		return nil, fmt.Errorf("%w: invalid namespace identifier", base.ErrResponseStatusBadRequest)
	}

	docSvc := base.GetAzCosmosCRUDService(c)
	globalDoc, err := readAgentConfigRadiusDocByLocator(c,
		base.NewDocLocator(base.NamespaceKindSystem,
			base.ID("default"), base.ResourceKindNamespaceConfig, base.ID(base.AgentConfigNameRadius)))
	if err != nil {
		return nil, err
	}
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
	doc.GlobalRadiusServerACRImageRef = globalDoc.GlobalRadiusServerACRImageRef
	globalDocVersionBytes, err := hex.DecodeString(globalDoc.Version)
	if err != nil {
		return nil, err
	}
	digest.Write(globalDocVersionBytes)

	if p.Clients != nil {
		doc.Clients = make([]radius.RadiusClientConfig, len(p.Clients))
		for i, pClient := range p.Clients {
			docClient := radius.RadiusClientConfig{
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
	doc.Version = hex.EncodeToString(digest.Sum(nil))

	err = docSvc.Upsert(c, doc, &azcosmos.ItemOptions{IfMatchEtag: doc.ETag})

	s.assignAgentRadiusRoles(c, nsUUID, doc)
	return doc, err
}
