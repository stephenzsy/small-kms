package common

import (
	"fmt"

	"github.com/google/uuid"
)

func GetCanonicalCertificateTemplateID(msGraphODataType string, objectID uuid.UUID, templateName string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("https://example.com/%s/certificate-templates/%s", msGraphODataType, templateName)))
}
