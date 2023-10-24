package cert

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type RequestContext = ctx.RequestContext

func deleteCertificate(c RequestContext, certDoc *CertDoc) error {
	if certDoc.Status != CertStatusIssued {
		return kmsdoc.Delete(c, certDoc)
	} else if certDoc.Owner != nil && !certDoc.Owner.IsNilOrEmpty() {
		c := ctx.Elevate(c)
		err := kmsdoc.Delete(c, certDoc)
		if err != nil {
			return err
		}

		patchOps := azcosmos.PatchOperations{}
		patchOps.AppendRemove(fmt.Sprintf("%s/%s", kmsdoc.PatchPathOwns, certDoc.NamespaceID))
		return kmsdoc.Patch(c, certDoc, patchOps, nil)
	} else {
		c := ctx.Elevate(c)
		patchOps := azcosmos.PatchOperations{}
		if len(certDoc.Owns) > 0 {
			for _, certLocator := range certDoc.Owns {
				err := kmsdoc.DeleteByRef(c, certLocator)
				if err != nil {
					return err
				}
			}
			patchOps.AppendRemove(kmsdoc.PatchPathOwns)
		}
		// disable in keyvault
		kid := certDoc.CertSpec.KID
		if kid != "" {
			if certDoc.NamespaceID.Kind() == shared.NamespaceKindCaRoot {
				kid := azkeys.ID(kid)
				// disable key
				_, err := common.GetAdminServerClientProvider(c).AzKeysClient().UpdateKey(c, kid.Name(), kid.Version(), azkeys.UpdateKeyParameters{
					KeyAttributes: &azkeys.KeyAttributes{
						Enabled: to.Ptr(false),
					},
				}, nil)
				if err != nil {
					return err
				}
			} else {
				kid := azcertificates.ID(kid)
				_, err := common.GetAdminServerClientProvider(c).AzCertificatesClient().UpdateCertificate(c, kid.Name(), kid.Version(), azcertificates.UpdateCertificateParameters{
					CertificateAttributes: &azcertificates.CertificateAttributes{
						Enabled: to.Ptr(false),
					},
				}, nil)
				if err != nil {
					return err
				}
			}
		}

		patchOps.AppendSet(kmsdoc.PatchPathDeleted, time.Now().UTC().Format(time.RFC3339))
		err := kmsdoc.Patch(c, certDoc, patchOps, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func ApiDeleteCertificate(c RequestContext, certificateId shared.Identifier) error {
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

	err = deleteCertificate(c, certDoc)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
