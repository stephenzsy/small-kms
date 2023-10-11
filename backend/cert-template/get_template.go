package certtemplate

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func getCertificateTemplateDocLocator(nsID models.NamespaceID, templateID common.Identifier) models.ResourceLocator {
	return common.NewLocator(nsID, common.NewIdentifierWithKind(models.ResourceKindCertTemplate, templateID))
}

func GetCertificateTemplateDoc(c RequestContext,
	locator models.ResourceLocator) (doc *CertificateTemplateDoc, err error) {

	if locator.GetID().Kind() != models.ResourceKindCertTemplate {
		return nil, fmt.Errorf("invalid resource type: %s, expected: %s", locator.GetID().Kind(), models.ResourceKindCertTemplate)
	}

	doc = new(CertificateTemplateDoc)
	err = kmsdoc.Read(c, locator, doc)
	return
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

func ListKeyVaultRoleAssignments(c RequestContext) ([]*models.KeyVaultRoleAssignment, error) {
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
	templateContext := GetCertificateTemplateContext(c)
	nsID := templateContext.GetCertificateTemplateLocator(c).GetNamespaceID()
	filterParam := fmt.Sprintf("assignedTo('{%s}')", nsID.Identifier().UUID().String())
	scope := delegatedClientProvider.GetKeyvaultCertificateResourceScopeID(*doc.KeyStorePath)
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
	return utils.PagerAllItems[*models.KeyVaultRoleAssignment](
		utils.NewMappedItemsPager[*models.KeyVaultRoleAssignment, *armauthorization.RoleAssignment](
			itemsPager, func(item *armauthorization.RoleAssignment) *models.KeyVaultRoleAssignment {
				if item == nil {
					return nil
				}
				return &models.KeyVaultRoleAssignment{
					Id:               *item.ID,
					PrincipalId:      *item.Properties.PrincipalID,
					RoleDefinitionId: *item.Properties.RoleDefinitionID,
				}
			}), c)
}
