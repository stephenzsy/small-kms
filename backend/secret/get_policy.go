package secret

import (
	"errors"
	"fmt"
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// ListSecretPolicy implements ServerInterface.
func (s *server) GetSecretPolicy(ec echo.Context, namespaceKind base.NamespaceKind, namespaceIdentifier base.ID, policyID base.ID) error {
	c := ec.(ctx.RequestContext)
	if !authz.AuthorizeAdminOnly(c) {
		return base.ErrResponseStatusForbidden
	}

	c = ns.WithNSContext(c, namespaceKind, namespaceIdentifier)
	doc, err := readSecretPolicyDoc(c, policyID)
	if err != nil {
		return wrapSecretPolicyNotFoundError(err, policyID)
	}
	model := &SecretPolicy{}
	doc.PopulateModel(model)
	return c.JSON(http.StatusOK, model)
}

func readSecretPolicyDoc(c ctx.RequestContext, policyIdentifier base.ID) (*SecretPolicyDoc, error) {
	nsCtx := ns.GetNSContext(c)
	doc := &SecretPolicyDoc{}
	err := base.GetAzCosmosCRUDService(c).Read(c, base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindSecretPolicy, policyIdentifier), doc, nil)
	return doc, err
}

func wrapSecretPolicyNotFoundError(err error, policyID base.ID) error {
	if errors.Is(err, base.ErrAzCosmosDocNotFound) {
		return fmt.Errorf("%w, secret policy not found: %s", base.ErrResponseStatusNotFound, policyID)
	}
	return err
}
