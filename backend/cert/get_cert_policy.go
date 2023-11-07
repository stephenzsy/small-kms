package cert

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

// GetCertPolicy implements ServerInterface.
func (s *server) GetCertPolicy(ec echo.Context, nsKind base.NamespaceKind, nsID base.ID, policyID base.ID) error {
	c := ec.(ctx.RequestContext)
	c, _, err := s.allowGeneralNonAdminAuth(c, nsKind, nsID)
	if err != nil {
		return err
	}

	if ns.VerifyKeyVaultIdentifier(policyID) != nil {
		return fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}
	doc, err := ReadCertPolicyDoc(c, policyID)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: certificate policy not found: %s", base.ErrResponseStatusNotFound, policyID)
		} else {
			return err
		}
	}
	m := new(CertPolicy)
	doc.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}

func readCertPolicyDocByLocator(c context.Context, locator base.DocLocator) (*CertPolicyDoc, error) {
	doc := new(CertPolicyDoc)
	err := base.GetAzCosmosCRUDService(c).Read(c, locator, doc, nil)
	return doc, err
}

func ReadCertPolicyDoc(c context.Context, rID base.ID) (*CertPolicyDoc, error) {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return nil, fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)

	return readCertPolicyDocByLocator(c, base.NewDocLocator(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy, rID))
}
