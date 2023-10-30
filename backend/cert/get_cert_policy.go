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

func ReadCertPolicyDoc(c context.Context, rID base.Identifier) (*CertPolicyDoc, error) {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return nil, fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)
	doc := new(CertPolicyDoc)

	slocator := base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, rID)

	err := base.GetAzCosmosCRUDService(c).Read(c, slocator, doc, nil)
	return doc, err
}

func apiGetCertPolicy(c ctx.RequestContext, rID base.Identifier) error {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}
	doc, err := ReadCertPolicyDoc(c, rID)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: certificate policy not found: %s", base.ErrResponseStatusNotFound, rID.String())
		} else {
			return err
		}
	}
	m := new(CertPolicy)
	doc.PopulateModel(m)
	return c.JSON(http.StatusOK, m)
}
