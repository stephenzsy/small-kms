package profile

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type server struct {
	api.APIServer
}

var (
	namespaceIdentifierCA              = base.ID("ca")
	namespaceIdentifierDirectoryObject = base.ID("directoryObject")
)

func getNamespaceIdentifier(profileResourceKind base.ResourceKind) (base.ID, error) {
	switch profileResourceKind {
	case base.ProfileResourceKindRootCA,
		base.ProfileResourceKindIntermediateCA:
		return namespaceIdentifierCA, nil
	case base.ProfileResourceKindServicePrincipal,
		base.ProfileResourceKindUser,
		base.ProfileResourceKindGroup:
		return namespaceIdentifierDirectoryObject, nil
	}
	return base.ID(""), fmt.Errorf("%w: invalid profile kind: %s", base.ErrResponseStatusBadRequest, profileResourceKind)
}

// GetRootCA implements ServerInterface.
func (s *server) GetProfile(ec echo.Context, profileResourceKind base.ResourceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	// nsId, err := getNamespaceIdentifier(profileResourceKind)
	// if err != nil {
	// 	return err
	// }
	// c = ns.WithDefaultNSContext(c, base.NamespaceKindProfile, nsId)
	return ec.JSON(400, map[string]string{"message": "not implemented"})
}

// PutRootCA implements ServerInterface.
func (s *server) PutProfile(ec echo.Context, profileResourceKind base.ResourceKind, namespaceIdentifier base.ID) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	params := new(ProfileParameters)
	if err := c.Bind(params); err != nil {
		return err
	}

	if err := ns.VerifyKeyVaultIdentifier(namespaceIdentifier); err != nil {
		return err
	}

	nsId, err := getNamespaceIdentifier(profileResourceKind)
	if err != nil {
		return err
	}
	c = ns.WithNSContext(c, base.NamespaceKindProfile, nsId)

	r, err := putProfile(c, profileResourceKind, namespaceIdentifier, params)
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
