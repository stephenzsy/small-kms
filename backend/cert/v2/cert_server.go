package cert

import (
	"fmt"
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

// GetCertificateRuleMsEntraClientCredential implements ServerInterface.
func (s *server) GetCertificateRuleMsEntraClientCredential(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiGetCertRuleMsEntraClientCredential(c)
}

// PutCertificateRuleMsEntraClientCredential implements ServerInterface.
func (s *server) PutCertificateRuleMsEntraClientCredential(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	params := new(CertificateRuleMsEntraClientCredential)
	if err := ec.Bind(params); err != nil {
		return err
	}

	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiPutCertRuleMsEntraClientCredentrial(c, params)
}

// PutCertificateRuleIssuer implements ServerInterface.
func (s *server) PutCertificateRuleIssuer(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	params := new(CertificateRuleIssuer)
	if err := ec.Bind(params); err != nil {
		return err
	}

	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiPutCertRuleIssuer(c, params)
}

func (s *server) withAdminAccessAndNamespaceCtx(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) (ctx.RequestContext, error) {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return c, s.RespondRequireAdmin(c)
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
	return c, nil
}

// GetCertificateRuleIssuer implements ServerInterface.
func (s *server) GetCertificateRuleIssuer(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiGetCertRuleIssuer(c)
}

// EnrollMsEntraClientCredential implements ServerInterface.
func (s *server) EnrollCertificate(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !namespaceIdentifier.IsUUID() {
		return fmt.Errorf("%w: invalid namespace identifier", base.ErrResponseStatusForbidden)
	}

	if !auth.AuthorizeApplicationOrAdmin(c, namespaceIdentifier.UUID()) {
		return s.RespondRequireAdmin(c)
	}

	params := new(EnrollCertificateRequest)
	if err := c.Bind(params); err != nil {
		return err
	}

	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)

	return enrollMsEntraClientCredCert(c, resourceIdentifier, params)
}

// GetCertificate implements ServerInterface.
func (s *server) GetCertificate(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	r, err := getCertificate(c, resourceIdentifier)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// DeleteCertificate implements ServerInterface.
func (s *server) DeleteCertificate(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	err = deleteCertificate(c, resourceIdentifier)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// ListCertificates implements ServerInterface.
func (s *server) ListCertificates(ec echo.Context, namespaceKind base.NamespaceKind,
	namespaceIdentifier base.Identifier, params ListCertificatesParams) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	l, err := listCertificates(c, params)
	if err != nil {
		return err
	}
	if l == nil {
		l = make([]*CertificateRef, 0)
	}
	return c.JSON(http.StatusOK, l)
}

// CreateCertificate implements ServerInterface.
func (s *server) CreateCertificate(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	r, err := createCertFromPolicy(c, resourceIdentifier, nil)
	if err != nil {
		return err
	}
	m := new(Certificate)
	r.PopulateModel(m)
	return c.JSON(http.StatusCreated, m)
}

// ListCertPolicies implements ServerInterface.
func (s *server) ListCertPolicies(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	l, err := listCertPolicies(c)
	if err != nil {
		return err
	}
	if l == nil {
		l = make([]*CertPolicyRef, 0)
	}
	return c.JSON(http.StatusOK, l)
}

// GetCertPolicy implements ServerInterface.
func (s *server) GetCertPolicy(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.Identifier, resourceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiGetCertPolicy(c, resourceIdentifier)
}

// PutCertPolicy implements ServerInterface.
func (s *server) PutCertPolicy(ec echo.Context,
	namespaceKind base.NamespaceKind,
	namespaceIdentifier base.Identifier,
	resourceIdentifier base.Identifier) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	params := new(CertPolicyParameters)
	if err := c.Bind(params); err != nil {
		return err
	}

	if err := ns.VerifyKeyVaultIdentifier(namespaceIdentifier); err != nil {
		return err
	}
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
