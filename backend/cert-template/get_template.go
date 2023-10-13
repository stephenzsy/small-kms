package certtemplate

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getCertificateTemplateDocLocator(nsID shared.NamespaceIdentifier, templateID shared.Identifier) shared.ResourceLocator {
	return shared.NewResourceLocator(nsID, shared.NewResourceIdentifier(shared.ResourceKindCertTemplate, templateID))
}

func getDirectCertificateTemplateDoc(c context.Context, locator shared.ResourceLocator) (doc *CertificateTemplateDoc, err error) {
	if locator.GetID().Kind() != shared.ResourceKindCertTemplate {
		return nil, fmt.Errorf("invalid resource type: %s, expected: %s", locator.GetID().Kind(), shared.ResourceKindCertTemplate)
	}
	doc = new(CertificateTemplateDoc)
	err = kmsdoc.Read(c, locator, doc)
	return
}

func GetCertificateTemplateDoc(c context.Context,
	locator shared.ResourceLocator) (doc *CertificateTemplateDoc, err error) {
	if doc, err = getDirectCertificateTemplateDoc(c, locator); err == nil && doc.ID.Identifier().IsUUID() && doc.ID.Identifier().UUID().Version() == 5 {
		if doc.AliasTo == nil {
			return nil, fmt.Errorf("%w: invalid template", common.ErrStatusBadRequest)
		}
		return getDirectCertificateTemplateDoc(c, *doc.AliasTo)
	} else {
		return doc, err
	}
}

// PutCertificateTemplate implements CertificateTemplateService.
func GetCertificateTemplate(c RequestContext,
) (*models.CertificateTemplateComposed, error) {

	templateLocator := GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c)
	doc, err := GetCertificateTemplateDoc(c, templateLocator)
	if err != nil {
		return nil, err
	}

	return doc.toModel(), nil
}

func ListKeyVaultRoleAssignments(c RequestContext) ([]*models.AzureRoleAssignment, error) {
	doc, err := GetCertificateTemplateDoc(c, GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c))
	if err != nil {
		return nil, err
	}
	if doc.KeyStorePath == nil || *doc.KeyStorePath == "" {
		return nil, fmt.Errorf("%w: key store path is empty", common.ErrStatusBadRequest)
	}
	delegatedClientProvider := common.GetAdminServerRequestClientProvider(c)
	raClient, err := delegatedClientProvider.ArmRoleAssignmentsClient()
	if err != nil {
		return nil, err
	}
	nsID := ns.GetNamespaceContext(c).GetID()
	filterParam := fmt.Sprintf("assignedTo('{%s}')", nsID.Identifier().UUID().String())
	scope := delegatedClientProvider.GetKeyvaultCertificateResourceScopeID(*doc.KeyStorePath, "secrets")
	log.Info().Msgf("Lookup role assignments for scope: %s", scope)
	pager := raClient.NewListForScopePager(
		scope,
		&armauthorization.RoleAssignmentsClientListForScopeOptions{
			Filter: &filterParam,
		},
	)

	itemsPager := utils.NewMappedPager[[]*armauthorization.RoleAssignment, armauthorization.RoleAssignmentsClientListForScopeResponse](
		pager,
		func(resp armauthorization.RoleAssignmentsClientListForScopeResponse) []*armauthorization.RoleAssignment {
			return resp.Value
		})
	return utils.PagerAllItems[*models.AzureRoleAssignment](
		utils.NewMappedItemsPager[*models.AzureRoleAssignment, *armauthorization.RoleAssignment](
			itemsPager, func(item *armauthorization.RoleAssignment) *models.AzureRoleAssignment {
				if item == nil {
					return nil
				}
				return &models.AzureRoleAssignment{
					Id:               item.ID,
					Name:             item.Name,
					PrincipalId:      item.Properties.PrincipalID,
					RoleDefinitionId: item.Properties.RoleDefinitionID,
				}
			}), c)
}

func DeleteKeyVaultRoleAssignment(c RequestContext, roleAssignmentID string) error {
	doc, err := GetCertificateTemplateDoc(c, GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c))
	if err != nil {
		return err
	}
	if doc.KeyStorePath == nil || *doc.KeyStorePath == "" {
		return fmt.Errorf("%w: key store path is empty", common.ErrStatusBadRequest)
	}
	delegatedClientProvider := common.GetAdminServerRequestClientProvider(c)
	scope := delegatedClientProvider.GetKeyvaultCertificateResourceScopeID(*doc.KeyStorePath, "secrets")
	raClient, err := delegatedClientProvider.ArmRoleAssignmentsClient()
	if err != nil {
		return err
	}
	resp, err := raClient.DeleteByID(c, fmt.Sprintf("%s/providers/Microsoft.Authorization/roleAssignments/%s", scope, roleAssignmentID), nil)
	log.Info().Msgf("Delete role assignment: %s, resp: %v, err: %v", roleAssignmentID, resp, err)
	return err
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

func AddKeyVaultRoleAssignment(c RequestContext, roleDefID uuid.UUID) (*models.AzureRoleAssignment, error) {
	doc, err := GetCertificateTemplateDoc(c, GetCertificateTemplateContext(c).GetCertificateTemplateLocator(c))
	if err != nil {
		return nil, err
	}
	if doc.KeyStorePath == nil || *doc.KeyStorePath == "" {
		return nil, fmt.Errorf("%w: key store path is empty", common.ErrStatusBadRequest)
	}
	delegatedClientProvider := common.GetAdminServerRequestClientProvider(c)
	scope := delegatedClientProvider.GetKeyvaultCertificateResourceScopeID(*doc.KeyStorePath, "secrets")
	raClient, err := delegatedClientProvider.ArmRoleAssignmentsClient()
	if err != nil {
		return nil, err
	}
	nsID := ns.GetNamespaceContext(c).GetID()
	roleAssignmentID := uuid.NewString()
	resp, err := raClient.Create(c, scope, roleAssignmentID, armauthorization.RoleAssignmentCreateParameters{
		Properties: &armauthorization.RoleAssignmentProperties{
			RoleDefinitionID: utils.ToPtr(fmt.Sprintf("%s/providers/Microsoft.Authorization/roleDefinitions/%s", scope, roleDefID.String())),
			PrincipalID:      nsID.Identifier().StringPtr(),
		},
	}, nil)
	if err != nil {
		return nil, err
	} else {
		log.Info().Msgf("Added role assignment: %s, resp: %v, err: %v", roleAssignmentID, resp, err)
	}
	return &models.AzureRoleAssignment{
		Id:               resp.ID,
		Name:             resp.Name,
		PrincipalId:      resp.Properties.PrincipalID,
		RoleDefinitionId: resp.Properties.RoleDefinitionID,
	}, nil
}
