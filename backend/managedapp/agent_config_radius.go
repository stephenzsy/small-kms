package managedapp

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

// GetAgentConfigRadius implements ServerInterface.
func (s *server) GetAgentConfigRadius(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	doc := &AgentConfigRadiusDoc{}
	if err := base.GetAzCosmosCRUDService(c).Read(c, base.NewDocLocator(nsKind,
		nsID, base.ResourceKindNamespaceConfig, base.ID(base.AgentConfigNameRadius)), doc, nil); err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: %s", base.ErrResponseStatusNotFound, base.AgentConfigNameRadius)
		}
		return err
	}
	m := &AgentConfigRadius{}
	doc.populateModel(m)
	return c.JSON(200, m)
}

// PutAgentConfigRadius implements ServerInterface.
func (s *server) PutAgentConfigRadius(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	param := new(AgentConfigRadiusFields)
	if err := c.Bind(param); err != nil {
		return err
	}

	digest := md5.New()
	doc := new(AgentConfigRadiusDoc)
	doc.init(nsKind, nsID)
	docSvc := base.GetAzCosmosCRUDService(c)
	switch nsKind {
	case base.NamespaceKindSystem:
		if string(nsID) != "default" {
			return fmt.Errorf("%w: only default system namespace is supported", base.ErrResponseStatusBadRequest)
		}
		doc.GlobalRadiusServerACRImageRef = *param.AzureACRImageRef
		digest.Write([]byte(doc.GlobalRadiusServerACRImageRef))
	case base.NamespaceKindServicePrincipal:
		globalDoc := new(AgentConfigRadiusDoc)
		if err := docSvc.Read(c,
			base.NewDocLocator(base.NamespaceKindSystem,
				base.ID("default"), base.ResourceKindNamespaceConfig, base.ID(base.AgentConfigNameRadius)), globalDoc, nil); err != nil {
			if errors.Is(err, base.ErrAzCosmosDocNotFound) {
				return fmt.Errorf("%w: %s", base.ErrResponseStatusNotFound, base.AgentConfigNameRadius)
			}
			return err
		}
		doc.GlobalRadiusServerACRImageRef = globalDoc.GlobalRadiusServerACRImageRef
		globalDocVersionBytes, err := hex.DecodeString(globalDoc.Version)
		if err != nil {
			return err
		}
		digest.Write(globalDocVersionBytes)
	default:
		return fmt.Errorf("%w: namespace kind %s is not supported", base.ErrResponseStatusBadRequest, nsKind)
	}
	doc.Version = hex.EncodeToString(digest.Sum(nil))
	err := docSvc.Upsert(c, doc, nil)
	if err != nil {
		return err
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
	panic("unimplemented")
}
