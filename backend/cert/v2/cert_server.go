package cert

import (
	"net/http"

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

// CreateCertificate implements ServerInterface.
func (s *server) CreateCertificate(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	r, err := createCertFromPolicy(c, resourceIdentifier)
	if err != nil {
		return err
	}
	m := new(Certificate)
	r.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}

// ListCertPolicies implements ServerInterface.
func (s *server) ListCertPolicies(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	l, err := listCertPolicies(c, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, l)
}

// GetCertPolicy implements ServerInterface.
func (s *server) GetCertPolicy(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	r, err := getCertPolicy(c, resourceIdentifier)
	if err != nil {
		return err
	}
	m := new(CertPolicy)
	r.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}

// PutCertPolicy implements ServerInterface.
func (s *server) PutCertPolicy(ec echo.Context,
	namespaceKind base.NamespaceKind,
	namespaceIdentifier base.Identifier,
	resourceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	params := new(CertPolicyParameters)
	if err := c.Bind(params); err != nil {
		return err
	}

	if err := ns.VerifyKeyVaultIdentifier(namespaceIdentifier); err != nil {
		return err
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	r, err := putCertPolicy(c, resourceIdentifier, params)
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
