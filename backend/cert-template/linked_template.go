package certtemplate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
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

func createLinkedCertificate(c context.Context, target shared.ResourceLocator) (*CertificateTemplateDoc, error) {
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
	tDoc := &CertificateTemplateDoc{
		BaseDoc: kmsdoc.BaseDoc{
			NamespaceID: nsID,
			ID:          transformedIdentifier,
			AliasTo:     &target,
			AliasToETag: &doc.ETag,
		},
	}
	err = kmsdoc.Upsert(c, tDoc)
	return tDoc, err
}

func ApiCreateLinkedCertificateTemplate(c RequestContext, params models.CreateLinkedCertificateTemplateParameters) error {
	template, err := createLinkedCertificate(c, params.TargetTemplate)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, template.toModel())
}
