package cert

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/stephenzsy/small-kms/backend/base"
)

func deleteCertificate(c context.Context, rID base.Identifier) error {
	doc, err := getCertDocByID(c, rID)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return nil
		}
		return err
	}

	if doc.Status != CertificateStatusPending {
		return errors.New("cannot delete certificate with status other than pending")
	}

	return base.GetAzCosmosCRUDService(c).Delete(c, doc.StorageNamespaceID, doc.StorageID, &azcosmos.ItemOptions{
		IfMatchEtag: doc.ETag,
	})
}
