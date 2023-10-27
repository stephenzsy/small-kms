package cert

import (
	"context"

	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func NewCertificateID(certId shared.Identifier) shared.ResourceIdentifier {
	return shared.NewResourceIdentifier(shared.ResourceKindCert, certId)
}

func ReadCertDocByLocator(c context.Context, locator shared.ResourceLocator) (*CertDoc, error) {
	certDoc := &CertDoc{}
	err := kmsdoc.Read(c, locator, certDoc)
	return certDoc, err
}
