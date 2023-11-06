package secret

import (
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

func (s *server) withAdminAccessAndNamespaceCtx(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) (ctx.RequestContext, error) {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return c, s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	return c, nil
}

// ListSecretPolicies implements ServerInterface.
func (s *server) ListSecretPolicies(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)
	c, err := s.withAdminAccessAndNamespaceCtx(c, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiListSecretPolicies(c)
}

// ListSecretPolicy implements ServerInterface.
func (s *server) GetSecretPolicy(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, resourceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)
	c, err := s.withAdminAccessAndNamespaceCtx(c, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiGetSecretPolicy(c, resourceIdentifier)
}

// PutSecretPolicy implements ServerInterface.
func (s *server) PutSecretPolicy(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, resourceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)
	c, err := s.withAdminAccessAndNamespaceCtx(c, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	params := SecretPolicyParameters{}
	if err := c.Bind(&params); err != nil {
		return err
	}
	return apiPutSecretPolicy(c, resourceIdentifier, params)
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
