package cert

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func getCertDocBySLocator(c context.Context, slocator base.SLocator) (*CertDoc, error) {
	doc := new(CertDoc)
	err := base.GetAzCosmosCRUDService(c).Read(c, slocator.NID, slocator.RID, doc, nil)
	return doc, err
}

func getCertDocByID(c context.Context, rID base.Identifier) (*CertDoc, error) {
	if !rID.IsUUID() {
		return nil, fmt.Errorf("%w: invalid resource identifier: %s", base.ErrResponseStatusBadRequest, rID.String())
	}

	nsCtx := ns.GetNSContext(c)
	slocator := base.SLocator{
		NID: base.GetDefaultStorageNamespaceID(nsCtx.Kind(), nsCtx.Identifier()),
		RID: rID.UUID(),
	}

	return getCertDocBySLocator(c, slocator)
}

func getCertificate(c context.Context, rID base.Identifier) (*Certificate, error) {
	doc, err := getCertDocByID(c, rID)
	m := new(Certificate)
	doc.PopulateModel(m)
	return m, err
}
