package cert

import (
	"context"
	"errors"
	"fmt"

	"github.com/stephenzsy/small-kms/backend/base"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

const (
	RelNameIssuerCert base.RelName = "issuer-cert"
)

type CertLinkRelDoc struct {
	base.BaseDoc
	RelName base.RelName `json:"relName"`
}

func getLinkDocIdentifierForPolicyID(policyID base.Identifier) base.Identifier {
	return base.StringIdentifier[string](fmt.Sprintf("%s:%s", policyID.String(), RelNameIssuerCert))
}

func (d *CertLinkRelDoc) Init(
	nsKind base.NamespaceKind,
	nsID base.Identifier,
	policyID base.Identifier,
) {
	d.NamespaceKind = nsKind
	d.NamespaceIdentifier = nsID
	d.ResourceKind = base.ResourceKindRel
	d.ResourceIdentifier = getLinkDocIdentifierForPolicyID(policyID)

	d.RelName = RelNameIssuerCert
}

func setIssuerCertRel(c context.Context, targetCert *CertDoc, policyID base.Identifier) error {
	nsCtx := ns.GetNSContext(c)
	relDoc := &CertLinkRelDoc{}
	relDoc.Init(
		nsCtx.Kind(),
		nsCtx.Identifier(),
		policyID,
	)
	linkLocator := base.GetDefaultStorageLocator(nsCtx.Kind(), nsCtx.Identifier(), base.ResourceKindRel, relDoc.ResourceIdentifier)
	docService := base.GetAzCosmosCRUDService(c)
	err := docService.Read(c, linkLocator.NID, linkLocator.RID, relDoc, nil)
	if err != nil {
		if !errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return err
		}
	}
	return base.UpsertRelDoc1To1(c, RelNameIssuerCert, targetCert, relDoc)
}

func setIssuerCert(c context.Context, policyID base.Identifier, params *PolicyIssuerCertRequest) error {
	nsCtx := ns.GetNSContext(c)
	switch nsCtx.Kind() {
	case base.NamespaceKindRootCA,
		base.NamespaceKindIntermediateCA:
		// ok
	default:
		return fmt.Errorf("%w: invalid namespace kind to set cert issuer: %s", base.ErrResponseStatusBadRequest, nsCtx.Kind())
	}
	certDoc, err := getCertDocByID(c, params.IssuerId)
	if err != nil {
		return err
	}
	return setIssuerCertRel(c, certDoc, policyID)
}

func getLinkDoc(c context.Context, nsKind base.NamespaceKind, nsId base.Identifier, policyId base.Identifier) (*CertLinkRelDoc, error) {
	doc := &CertLinkRelDoc{}
	docService := base.GetAzCosmosCRUDService(c)
	linkDocLocator := base.GetDefaultStorageLocator(nsKind, nsId, base.ResourceKindRel, base.StringIdentifier[string](fmt.Sprintf("%s:%s", policyId.String(), RelNameIssuerCert)))
	err := docService.Read(c, linkDocLocator.NID, linkDocLocator.RID, doc, nil)
	return doc, err
}

func (d *CertLinkRelDoc) getLinkedToCertDoc(c context.Context) (*CertDoc, error) {
	if d.Relations == nil || d.Relations.NamedTo == nil {
		return nil, fmt.Errorf("%w: no certificate found", base.ErrResponseStatusBadRequest)
	}
	if certDocLocator, ok := d.Relations.NamedTo[RelNameIssuerCert]; ok {
		return getCertDocBySLocator(c, certDocLocator)
	}
	return nil, fmt.Errorf("%w: no certificate found", base.ErrResponseStatusBadRequest)
}
