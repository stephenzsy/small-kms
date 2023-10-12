package cert

import (
	"fmt"

	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func DeleteCertificate(c RequestContext, certificateId shared.Identifier) error {
	if !certificateId.IsUUID() {
		return fmt.Errorf("%w: invalid certificate ID for delete: %s", common.ErrStatusBadRequest, certificateId)
	}

	nsID := ns.GetNamespaceContext(c).GetID()
	certDocLocator := shared.NewResourceLocator(nsID, NewCertificateID(certificateId))

	certDoc := &CertDoc{}
	err := kmsdoc.Read(c, certDocLocator, certDoc)
	if err != nil {
		if err == common.ErrStatusNotFound {
			return nil
		}
		return err
	}

	if certDoc.Status == CertStatusIssued {
		return fmt.Errorf("%w: cannot delete issued certificate", common.ErrStatusBadRequest)
	}

	return kmsdoc.DeleteByRef(c, certDocLocator)
}
