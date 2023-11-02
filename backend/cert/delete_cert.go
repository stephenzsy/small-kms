package cert

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/base"
)

func deleteCertificate(c context.Context, rID base.Identifier) error {
	doc, err := ReadCertDocByID(c, rID)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil
		}
		return err
	}

	docSvc := base.GetAzCosmosCRUDService(c)
	if doc.Status == CertificateStatusPending {
		return docSvc.Delete(c, doc, &azcosmos.ItemOptions{
			IfMatchEtag: doc.ETag,
		})
	} else if doc.Status == CertificateStatusIssued {
		// check if certificate is expired, we can delete expired certificate
		if doc.NotAfter.Time.Before(time.Now()) {
			return docSvc.Delete(c, doc, &azcosmos.ItemOptions{
				IfMatchEtag: doc.ETag,
			})
		}
		// soft delete
		return docSvc.SoftDelete(c, doc, &azcosmos.ItemOptions{
			IfMatchEtag: doc.ETag,
		})
	}
	return fmt.Errorf("%w: cannot delete certificate with status %s", base.ErrResponseStatusBadRequest, doc.Status)
}
