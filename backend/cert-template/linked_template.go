package certtemplate

import (
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func createTemplateLinkIdentifier(target shared.ResourceLocator) shared.ResourceIdentifier {
	newIdentifierUuid := uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("https://example.com/%s", target.String())))
	return shared.NewResourceIdentifier(shared.ResourceKindCertTemplate,
		shared.UUIDIdentifier(newIdentifierUuid))
}

func createLinkedCertificateTemplate(c RequestContext, target shared.ResourceLocator, usage models.LinkedCertificateTemplateUsage) (*CertificateTemplateDoc, error) {
	if target.GetID().Identifier().IsUUID() && target.GetID().Identifier().UUID().Version() == 5 {
		return nil, fmt.Errorf("%w:cannot create a link to another link", common.ErrStatusBadRequest)
	}
	doc, err := getDirectCertificateTemplateDoc(c, target)
	if err != nil {
		return nil, err
	}
	if doc.AliasTo != nil {
		return nil, fmt.Errorf("%w:cannot create a link to another link", common.ErrStatusBadRequest)
	}

	nsID := ns.GetNamespaceContext(c).GetID()
	if nsID == target.GetNamespaceID() {
		return nil, fmt.Errorf("%w:cannot create a link within the same namespace", common.ErrStatusBadRequest)
	}

	transformedIdentifier := createTemplateLinkIdentifier(target)
	tDoc := *doc
	tDoc.NamespaceID = nsID
	tDoc.ID = transformedIdentifier
	tDoc.Owner = &target
	tDoc.LinkProperties = &CertificateTemplateDocLink{
		Usage: usage,
	}
	tDoc.Owns = nil
	{
		c := ctx.Elevate(c)
		err = kmsdoc.Upsert(c, &tDoc)
		if err != nil {
			return nil, err
		}
		patchOps := azcosmos.PatchOperations{}
		if doc.Owns == nil {
			patchOps.AppendSet(kmsdoc.PatchPathOwns, map[shared.NamespaceIdentifier]shared.ResourceLocator{
				nsID: tDoc.GetLocator(),
			})
		} else {
			patchOps.AppendSet(fmt.Sprintf("%s/%s", kmsdoc.PatchPathOwns, nsID), tDoc.GetLocator())
		}
		err = kmsdoc.Patch(c, doc, patchOps, &azcosmos.ItemOptions{
			IfMatchEtag: &doc.ETag,
		})
		return &tDoc, err
	}
}

func ApiCreateLinkedCertificateTemplate(c RequestContext, params models.CreateLinkedCertificateTemplateParameters) error {
	switch params.Usage {
	case models.LinkedCertificateTemplateUsageClientAuthorization, models.LinkedCertificateTemplateUsageMemberDelegatedEnrollment:
		// ok
	default:
		return fmt.Errorf("%w: invalid usage: %s", common.ErrStatusBadRequest, params.Usage)
	}
	template, err := createLinkedCertificateTemplate(c, params.TargetTemplate, params.Usage)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, template.toModel())
}
