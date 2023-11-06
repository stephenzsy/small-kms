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
	"github.com/stephenzsy/small-kms/backend/utils"
)

type server struct {
	api.APIServer
}

// AddKeyVaultRoleAssignment implements ServerInterface.
func (*server) AddKeyVaultRoleAssignment(ctx echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier ID, resourceIdentifier ID, kvCategory AzureKeyvaultResourceCategory, params AddKeyVaultRoleAssignmentParams) error {
	return fmt.Errorf("%w: unimplemented", base.ErrResponseStatusBadRequest)
}

// ListKeyVaultRoleAssignments implements ServerInterface.
func (s *server) ListKeyVaultRoleAssignments(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier ID, resourceIdentifier ID, kvCategory AzureKeyvaultResourceCategory) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}
	return s.apiListKeyVaultRoleAssignments(c, resourceIdentifier, kvCategory)
}

// GetCertificateRuleMsEntraClientCredential implements ServerInterface.
func (s *server) GetCertificateRuleMsEntraClientCredential(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiGetCertRuleMsEntraClientCredential(c)
}

// PutCertificateRuleMsEntraClientCredential implements ServerInterface.
func (s *server) PutCertificateRuleMsEntraClientCredential(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
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
func (s *server) PutCertificateRuleIssuer(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) error {
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

func (s *server) withAdminAccessAndNamespaceCtx(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID) (ctx.RequestContext, error) {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return c, base.ErrResponseStatusForbidden
	}
	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)
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

	c = ns.WithDefaultNSContext(c, namespaceKind, namespaceIdentifier)

	r, err := apiGetCertificate(c, resourceIdentifier, true)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// DeleteCertificate implements ServerInterface.
func (s *server) DeleteCertificate(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, resourceIdentifier base.ID) error {
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
	namespaceIdentifier base.ID, params ListCertificatesParams) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiListCertificates(c, params)
}

// CreateCertificate implements ServerInterface.
func (s *server) CreateCertificate(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID, policyID base.ID) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, nsKind, nsID)
	if err != nil {
		return err
	}

	r, err := createCertFromPolicy(c, base.NewDocFullIdentifier(nsKind, nsID, base.ResourceKindCertPolicy, policyID), nil)
	if err != nil {
		return err
	}
	m := new(Certificate)
	r.PopulateModel(m)
	return c.JSON(http.StatusCreated, m)
}

// GetCertPolicy implements ServerInterface.
func (s *server) GetCertPolicy(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, resourceIdentifier base.ID) error {
	c, err := s.withAdminAccessAndNamespaceCtx(ec, namespaceKind, namespaceIdentifier)
	if err != nil {
		return err
	}

	return apiGetCertPolicy(c, resourceIdentifier)
}

// PutCertPolicy implements ServerInterface.
func (s *server) PutCertPolicy(ec echo.Context,
	namespaceKind base.NamespaceKind,
	namespaceIdentifier base.ID,
	resourceIdentifier base.ID) error {
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
