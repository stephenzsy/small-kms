package cert

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func getCertPolicy(c context.Context, rID base.Identifier) (*CertPolicyDoc, error) {

	if ns.VerifyKeyVaultIdentifier(rID) != nil {
		return nil, fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)
	doc := new(CertPolicyDoc)

	slocator := base.GetDefaultStorageLocator(c, nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCertPolicy, rID)

	err := base.GetAzCosmosCRUDService(c).Read(c, slocator.NID, slocator.RID, doc, nil)
	return doc, err
}
