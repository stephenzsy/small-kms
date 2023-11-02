package cert

import (
	"context"
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func ReadCertDocByFullIdentifier(c context.Context, fullIdentifier base.DocFullIdentifier) (*CertDoc, error) {
	doc := new(CertDoc)
	err := base.GetAzCosmosCRUDService(c).Read(c, fullIdentifier, doc, nil)
	return doc, err
}

// wraps 404
func ApiReadCertDocByID(c context.Context, rID base.Identifier) (*CertDoc, error) {
	if !rID.IsUUID() {
		return nil, fmt.Errorf("%w: invalid resource identifier: %s", base.ErrResponseStatusBadRequest, rID.String())
	}

	nsCtx := ns.GetNSContext(c)
	doc, err := ReadCertDocByFullIdentifier(c, base.NewDocFullIdentifier(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindCert, rID))
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil, fmt.Errorf("%w: cert with id %s not found", base.ErrResponseStatusNotFound, rID.String())
		}
	}
	return doc, err
}

func apiGetCertificate(c context.Context, rID base.Identifier, isAdminOrSelf bool) (*Certificate, error) {
	doc, err := ApiReadCertDocByID(c, rID)
	m := new(Certificate)
	doc.PopulateModel(m)

	if isAdminOrSelf && doc.KeyVaultStore.SID != "" {
		m.KeyVaultSecretID = &doc.KeyVaultStore.SID
	}
	return m, err
}
