package managedapp

import (
	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type server struct {
	api.APIServer
}

// GetAgentInstance implements ServerInterface.
func (s *server) GetAgentInstance(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, resourceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)

	return apiGetAgentInstance(c, resourceIdentifier)
}

// ListAgentInstances implements ServerInterface.
func (s *server) ListAgentInstances(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)
	return apiListAgentInstances(c)
}

// PutAgentConfigServerInstance implements ServerInterface.
func (s *server) PutAgentInstance(ec echo.Context, namespaceKind base.NamespaceKind, nsID base.ID, instanceId base.ID) error {
	c := ec.(ctx.RequestContext)
	nsUUID := auth.ResolveSelfNamespace(c, string(nsID))
	if !auth.AuthorizeSelfOrAdmin(c, nsUUID) {
		s.RespondRequireAdmin(c)
	} else if !utils.IsUUIDNil(nsUUID) {
		nsID = base.IDFromUUID(nsUUID)
	}
	c = ns.WithNSContext(c, namespaceKind, nsID)
	fields := AgentInstanceFields{}
	if err := c.Bind(&fields); err != nil {
		return err
	}
	return apiPutAgentInstance(c, instanceId, fields)
}

// GetAgentConfigServer implements ServerInterface.
func (s *server) GetAgentConfigServer(ec echo.Context, namespaceKind base.NamespaceKind, nsID base.ID) error {
	c := ec.(ctx.RequestContext)

	nsUUID := auth.ResolveSelfNamespace(c, string(nsID))
	if !auth.AuthorizeSelfOrAdmin(c, nsUUID) {
		s.RespondRequireAdmin(c)
	} else if !utils.IsUUIDNil(nsUUID) {
		nsID = base.IDFromUUID(nsUUID)
	}
	c = ns.WithNSContext(c, namespaceKind, nsID)

	return apiGetAgentConfigServer(c)
}

func (s *server) PutAgentConfigServer(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)

	param := &AgentConfigServerFields{}
	if err := c.Bind(param); err != nil {
		return err
	}

	return s.apiPutAgentConfigServer(c, param)
}

// GetManagedApp implements ServerInterface.
func (s *server) GetManagedApp(ec echo.Context, managedAppId uuid.UUID) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	doc, err := getManagedApp(c, managedAppId)
	if err != nil {
		return err
	}

	m := &ManagedApp{}
	doc.PopulateModel(m)
	return c.JSON(200, doc)
}

// ListManagedApps implements ServerInterface.
func (s *server) ListManagedApps(ec echo.Context) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	return apiListManagedApps(c)
}

// CreateManagedApp implements ServerInterface.
func (s *server) CreateManagedApp(ec echo.Context) error {
	c := ec.(ctx.RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	param := &ManagedAppParameters{}
	if err := c.Bind(param); err != nil {
		return err
	}

	result, err := createManagedApp(c, param)
	if err != nil {
		return err
	}

	m := &ManagedApp{}
	result.PopulateModel(m)
	return c.JSON(200, m)
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		APIServer: apiServer,
	}
}
