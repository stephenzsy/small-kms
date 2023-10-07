package admin

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

var (
	ErrCertificateTemplateVariable = errors.New("certificate template variable field is invalid")
)

func validateTemplateIdentifiers(nsID uuid.UUID, templateID uuid.UUID, name string, allowNamespaceDefault, allowEntraClientCredsDefault bool) (string, bool) {
	if templateID.Version() == 4 {
		// only allow non default prefixed name for user specified template ID
		return name, (name != "" && !strings.HasPrefix(name, "default"))
	} else if templateID.Version() == 5 {
		if allowNamespaceDefault && templateID == common.GetCanonicalCertificateTemplateID(nsID, common.DefaultCertTemplateName_GlobalDefault) {
			return string(common.DefaultCertTemplateName_GlobalDefault), true
		}
		if allowEntraClientCredsDefault && templateID == common.GetCanonicalCertificateTemplateID(nsID, common.DefaultCertTemplateName_ServicePrincipalClientCredential) {
			return string(common.DefaultCertTemplateName_ServicePrincipalClientCredential), true
		}
	}
	return "invalid", false
}
