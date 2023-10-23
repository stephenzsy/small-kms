package cert

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func putCertPolicy(c context.Context, rID base.Identifier, params *CertPolicyParameters) (*CertPolicy, error) {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return nil, fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)
	doc := new(CertPolicyDoc)

	doc.Init(nsCtx.Kind(), nsCtx.Identifier(), rID, params)
	err := base.GetAzCosmosCRUDService(c).Upsert(c, doc, nil)
	if err != nil {
		return nil, err
	}

	m := new(CertPolicy)
	doc.PopulateModel(m)
	return m, nil
}
