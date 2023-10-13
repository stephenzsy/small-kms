package cert

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

func CreateLinkedCertificate(c context.Context, targetCertDoc *CertDoc) (*CertDoc, error) {
	nsID := ns.GetNamespaceContext(c).GetID()
	targetLocator := targetCertDoc.GetLocator()
	certDocClone := *targetCertDoc
	certDocClone.NamespaceID = nsID
	certDocClone.ID = targetLocator.GetID()
	certDocClone.AliasTo = &targetLocator
	if targetCertDoc.AliasToETag == nil {
		certDocClone.AliasToETag = &targetCertDoc.ETag
	} else {
		certDocClone.AliasToETag = targetCertDoc.AliasToETag
	}
	err := kmsdoc.Upsert(c, &certDocClone)
	return &certDocClone, err
}
