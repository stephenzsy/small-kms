package cert

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func getCertificate(c context.Context, rID base.Identifier) (*Certificate, error) {

	if !rID.IsUUID() {
		return nil, fmt.Errorf("%w: invalid resource identifier", base.ErrResponseStatusBadRequest)
	}

	nsCtx := ns.GetNSContext(c)
	doc := new(CertDoc)

	slocator := base.SLocator{
		NID: base.GetDefaultStorageNamespaceID(c, nsCtx.Kind(), nsCtx.Identifier()),
		RID: rID.UUID(),
	}

	err := base.GetAzCosmosCRUDService(c).Read(c, slocator.NID, slocator.RID, doc, nil)
	m := new(Certificate)
	doc.PopulateModel(m)
	return m, err
}
