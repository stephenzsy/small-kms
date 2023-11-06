package cert

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/stephenzsy/small-kms/backend/base"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func readCertPolicyDocByLocator(c context.Context, locator base.DocFullIdentifier) (*CertPolicyDoc, error) {
	doc := new(CertPolicyDoc)
	err := base.GetAzCosmosCRUDService(c).Read(c, locator, doc, nil)
	return doc, err
}

func ReadCertPolicyDoc(c context.Context, rID base.ID) (*CertPolicyDoc, error) {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return nil, fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)

	return readCertPolicyDocByLocator(c, base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.ID(), base.ResourceKindCertPolicy, rID))
}

func apiGetCertPolicy(c ctx.RequestContext, rID base.ID) error {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}
	doc, err := ReadCertPolicyDoc(c, rID)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: certificate policy not found: %s", base.ErrResponseStatusNotFound, rID)
		} else {
			return err
		}
	}
	m := new(CertPolicy)
	doc.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}
