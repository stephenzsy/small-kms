package managedapp

import (
	"net/http"

	"github.com/google/uuid"
	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type server struct {
	api.APIServer
}

// ListAgentServerAzureRoleAssignments implements ServerInterface.
func (s *server) ListAgentServerAzureRoleAssignments(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)

	return s.apiListAgentConfigServerRoleAssignments(c)
}

// SyncSystemApp implements ServerInterface.
func (s *server) SyncSystemApp(ec echo.Context, systemAppName SystemAppName) error {
	c := ec.(ctx.RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	return apiSyncSystemApp(c, systemAppName)
}

// GetSystemApp implements ServerInterface.
func (s *server) GetSystemApp(ec echo.Context, systemAppName SystemAppName) error {
	c := ec.(ctx.RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	return apiGetSystemApp(c, systemAppName)
}

// GetAgentConfigServer implements ServerInterface.
func (s *server) GetAgentConfigServer(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if namespaceIdentifier != base.StringIdentifier("me") && auth.AuthorizeAdminOnly(c) {
		// ok
	} else if authedNamespaceId, ok := auth.AuthorizeApplicationMe(c, namespaceIdentifier.UUID(), namespaceIdentifier == base.StringIdentifier("me")); !ok {
		return c.JSON(http.StatusForbidden, map[string]string{"message": "unauthorized"})
	} else {
		namespaceIdentifier = base.UUIDIdentifier(authedNamespaceId)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)

	return apiGetAgentConfigServer(c)
}

func (s *server) PutAgentConfigServer(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)

	param := &AgentConfigServerFields{}
	if err := c.Bind(param); err != nil {
		return err
	}

	return s.apiPutAgentConfigServer(c, param)
}

// SyncManagedApp implements ServerInterface.
func (s *server) SyncManagedApp(ec echo.Context, managedAppId uuid.UUID) error {
	c := ec.(ctx.RequestContext)
	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	return apiSyncManagedApp(c, managedAppId)
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

	result, err := listManagedApps(c)
	if err != nil {
		return err
	}
	if result == nil {
		result = []*ManagedApp{}
	}
	return c.JSON(200, result)
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
		apiServer,
	}
}
