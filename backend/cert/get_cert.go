package cert

import (
	"context"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func getCertDocBySLocator(c context.Context, slocator base.DocFullIdentifier) (*CertDoc, error) {
	doc := new(CertDoc)
	err := base.GetAzCosmosCRUDService(c).Read(c, slocator, doc, nil)
	return doc, err
}

func getCertDocByID(c context.Context, rID base.Identifier) (*CertDoc, error) {
	if !rID.IsUUID() {
		return nil, fmt.Errorf("%w: invalid resource identifier: %s", base.ErrResponseStatusBadRequest, rID.String())
	}

	nsCtx := ns.GetNSContext(c)
	return getCertDocBySLocator(c, base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCert, rID))
}

func apiGetCertificate(c context.Context, rID base.Identifier, isAdminOrSelf bool) (*Certificate, error) {
	doc, err := getCertDocByID(c, rID)
	m := new(Certificate)
	doc.PopulateModel(m)

	if isAdminOrSelf && doc.KeyVaultStore.SID != "" {
		m.KeyVaultSecretID = &doc.KeyVaultStore.SID
	}
	return m, err
}
