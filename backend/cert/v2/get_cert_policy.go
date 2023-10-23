package cert

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func getCertPolicy(c context.Context, rID base.Identifier) (*CertPolicy, error) {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return nil, fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)
	doc := new(CertPolicyDoc)

	nid, rid := base.GetDefaultStorageLocator(c, nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, rID)

	err := base.GetAzCosmosCRUDService(c).Read(c, nid, rid, doc, nil)
	if err != nil {
		return nil, err
	}

	m := new(CertPolicy)
	doc.PopulateModel(m)
	return m, nil
}
