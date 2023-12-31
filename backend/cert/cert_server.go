package cert

import (
	"net/http"

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

// ListKeyVaultRoleAssignments implements ServerInterface.
func (s *server) ListKeyVaultRoleAssignments(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier ID, resourceIdentifier ID, kvCategory AzureKeyvaultResourceCategory) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}
	return s.apiListKeyVaultRoleAssignments(c, resourceIdentifier, kvCategory)
}

func (s *server) withAdminAccessAndNamespaceCtx(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) (ctx.RequestContext, error) {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return c, base.ErrResponseStatusForbidden
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)
	return c, nil
}

// GetCertificateRuleIssuer implements ServerInterface.
func (s *server) GetCertificateRuleIssuer(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiGetCertRuleIssuer(c)
}

// GetCertificate implements ServerInterface.
func (s *server) GetCertificate(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, resourceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)
	namespaceUUID := auth.ResolveSelfNamespace(c, string(namespaceIdentifier))
	if !auth.AuthorizeSelfOrAdmin(c, namespaceUUID) && !auth.HasRole(c, auth.RoleValueAgentActiveHost) {
		s.RespondRequireAdmin(c)
	} else if !utils.IsUUIDNil(namespaceUUID) {
		namespaceIdentifier = base.IDFromUUID(namespaceUUID)
	}

	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)

	r, err := apiGetCertificate(c, resourceIdentifier, true)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
