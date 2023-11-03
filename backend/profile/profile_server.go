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

const (
	namespaceIDCA              = "ca"
	namespaceIDDirectoryObject = "directoryObject"
)

var (
	namespaceIdentifierCA              = base.StringIdentifier(namespaceIDCA)
	namespaceIdentifierDirectoryObject = base.StringIdentifier(namespaceIDDirectoryObject)
)

func getNamespaceIdentifier(profileResourceKind base.ResourceKind) (base.Identifier, error) {
	switch profileResourceKind {
	case base.ProfileResourceKindRootCA,
		base.ProfileResourceKindIntermediateCA:
		return namespaceIdentifierCA, nil
	case base.ProfileResourceKindServicePrincipal:
		return namespaceIdentifierDirectoryObject, nil
	}
	return base.Identifier{}, fmt.Errorf("%w: invalid profile kind: %s", base.ErrResponseStatusBadRequest, profileResourceKind)
}

// ListRootCAs implements ServerInterface.
func (s *server) ListProfiles(ec echo.Context, profileResourceKind base.ResourceKind) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	nsId, err := getNamespaceIdentifier(profileResourceKind)
	if err != nil {
		return err
	}
	c = ns.WithDefaultNSContext(c, base.NamespaceKindProfile, nsId)
	return apiListProfiles(c, profileResourceKind)
}

// GetRootCA implements ServerInterface.
func (s *server) GetProfile(ec echo.Context, profileResourceKind base.ResourceKind, namespaceIdentifier base.Identifier) error {
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
func (s *server) PutProfile(ec echo.Context, profileResourceKind base.ResourceKind, namespaceIdentifier base.Identifier) error {
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
	c = ns.WithDefaultNSContext(c, base.NamespaceKindProfile, nsId)

	r, err := putProfile(c, profileResourceKind, namespaceIdentifier, params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

// ImportProfile implements ServerInterface.
func (s *server) ImportProfile(ec echo.Context,
	profileResourceKind base.ResourceKind, namespaceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	if err := ns.VerifyKeyVaultIdentifier(namespaceIdentifier); err != nil {
		return err
	}

	nsId, err := getNamespaceIdentifier(profileResourceKind)
	if err != nil {
		return err
	}
	c = ns.WithDefaultNSContext(c, base.NamespaceKindProfile, nsId)

	r, err := importProfile(c, profileResourceKind, namespaceIdentifier)
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
