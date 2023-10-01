package common

import (
	"fmt"

	"github.com/google/uuid"
)

func getCanonicalUUID(namespaceID uuid.UUID, typeName, name string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(fmt.Sprintf("https://example.com/%s/%s/%s", namespaceID, typeName, name)))
}

func GetCanonicalCertificateTemplateID(namespaceID uuid.UUID, templateName string) uuid.UUID {
	return getCanonicalUUID(namespaceID, "certificate-templates", templateName)
}

func GetCanonicalNamespaceRelationID(namespaceID uuid.UUID, relationName string) uuid.UUID {
	return getCanonicalUUID(namespaceID, "rel", relationName)
}
