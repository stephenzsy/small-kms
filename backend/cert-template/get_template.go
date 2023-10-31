package certtemplate

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/shared"
)

func getCertificateTemplateDocLocator(nsID shared.NamespaceIdentifier, templateID shared.Identifier) shared.ResourceLocator {
	return shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindCertTemplate, templateID))
}

func DeleteKeyVaultRoleAssignment(c RequestContext, roleAssignmentID string) error {
	// doc, err := GetCertificateTemplateDoc(c, GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c))
	// if err != nil {
	// 	return err
	// }
	// if doc.KeyStorePath == nil || *doc.KeyStorePath == "" {
	// 	return fmt.Errorf("%w: key store path is empty", common.ErrStatusBadRequest)
	// }
	// delegatedClientProvider := common.GetAdminServerRequestClientProvider(c)
	// scope := delegatedClientProvider.GetKeyvaultCertificateResourceScopeID(*doc.KeyStorePath, "secrets")
	// raClient, err := delegatedClientProvider.ArmRoleAssignmentsClient()
	// if err != nil {
	// 	return err
	// }
	// resp, err := raClient.DeleteByID(c, fmt.Sprintf("%s/providers/Microsoft.Authorization/roleAssignments/%s", scope, roleAssignmentID), nil)
	// log.Info().Msgf("Delete role assignment: %s, resp: %v, err: %v", roleAssignmentID, resp, err)
	// return err
	return nil
}

var roleAssignmentCategories = map[uuid.UUID]string{
	uuid.MustParse("4633458b-17de-408a-b874-0445c86b69e6"): "secrets",
}

func ValidateRoleDefnitionIDForAdd(inputId string) (uuid.UUID, error) {
	id, err := uuid.Parse(inputId)
	if err != nil {
		return id, err
	}
	if category, ok := roleAssignmentCategories[id]; !ok {
		return id, fmt.Errorf("%w: role assignment category: %s", common.ErrStatusBadRequest, category)
	}
	return id, nil
}
