package cert

import (
	"fmt"

	"github.com/rs/zerolog/log"
	ct "github.com/stephenzsy/small-kms/backend/cert-template"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type RequestContext = common.RequestContext

/*
type certServiceContextKey string

const (
	certContext certServiceContextKey = "certContext"
)
*/
// IssueCertificateFromTemplate implements CertificateService.
func IssueCertificateFromTemplate(
	c RequestContext,
	params models.IssueCertificateFromTemplateParams) (*shared.CertificateInfo, error) {

	nsID := ns.GetNamespaceContext(c).GetID()
	ctc := ct.GetCertificateTemplateContext(c)

	// verify template
	tmplDoc, err := ctc.GetCertificateTemplateDoc(c)
	if err != nil {
		return nil, err
	}
	if utils.IsTimeNotNilOrZero(tmplDoc.Deleted) {
		return nil, fmt.Errorf("%w: template not found or no longer active", common.ErrStatusNotFound)
	}
	if tmplDoc.Owner != nil {
		return nil, fmt.Errorf("%w: cannot create certificate for linked template", common.ErrStatusBadRequest)
	}

	newCert := false
	var certDoc *CertDoc
	if params.Force != nil && *params.Force {
		log.Info().Msg("force issue certificate")
		newCert = true
	} else {
		certDoc, newCert, err = shouldCreateNewCertificate(c, tmplDoc)
		if err != nil {
			return nil, err
		}
	}

	if newCert {
		certDoc, err := prepareNewCertDoc(nsID, tmplDoc)
		if err != nil {
			return nil, err
		}

		_, certDoc, err = issueCertificate(c, certDoc)
		if err != nil {
			return nil, err
		}
		return certDoc.toModel(), common.ErrStatus2xxCreated
	}
	return certDoc.toModel(), nil
}
