package profile

import (
	"fmt"
	"net/http"
	"regexp"

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

const namespaceIDCA = "ca"

var namespaceIdentifierCA = base.StringIdentifier(namespaceIDCA)

// ListRootCAs implements ServerInterface.
func (s *server) ListRootCAs(ec echo.Context) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	result, err := listProfiles(c, namespaceIdentifierCA, base.ProfileResourceKindRootCA)
	if err != nil {
		return err
	}
	if result == nil {
		result = []*ProfileRef{}
	}
	return c.JSON(http.StatusOK, result)
}

// GetRootCA implements ServerInterface.
func (s *server) GetRootCA(ec echo.Context, namespaceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	return ec.JSON(400, map[string]string{"message": "not implemented"})
}

// PutRootCA implements ServerInterface.
func (s *server) PutRootCA(ec echo.Context, namespaceIdentifier base.Identifier) error {
	c := ec.(ctx.RequestContext)

	if !auth.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}

	params := new(ProfileParameters)
	if err := c.Bind(params); err != nil {
		return err
	}

	if err := verifyProfileName(namespaceIdentifier); err != nil {
		return err
	}
	c = ns.WithDefaultNSContext(c, base.NamespaceKindProfile, namespaceIdentifierCA)

	r, err := putProfile(c, base.ProfileResourceKindRootCA, namespaceIdentifier, params)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, r)
}

var profileNameRegex = regexp.MustCompile(`^[0-9A-Za-z\-]+$`)

func verifyProfileName(identifier base.Identifier) error {
	if identifier.IsUUID() {
		return nil
	}
	str := identifier.String()
	if len(str) < 1 || len(str) > 96 || !profileNameRegex.MatchString(str) {
		return fmt.Errorf("%w: invalid identifier, %s", base.ErrResponseStatusBadRequest, str)
	}

	return nil
}

var _ ServerInterface = (*server)(nil)

func NewServer(apiServer api.APIServer) *server {
	return &server{
		apiServer,
	}
}
