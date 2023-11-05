package cert

import (
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cloudutils"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func (s *server) apiListKeyVaultRoleAssignments(c ctx.RequestContext, policyIdentifier base.Identifier, kvCategory AzureKeyvaultResourceCategory) error {
	switch kvCategory {
	case AzureKeyvaultResourceCategoryCertificates,
		AzureKeyvaultResourceCategoryKeys,
		AzureKeyvaultResourceCategorySecrets:
		// ok
	default:
		return fmt.Errorf("%w: invalid keyvault resource category", base.ErrResponseStatusBadRequest)
	}

	// verify policy exists
	policyDoc, err := ReadCertPolicyDoc(c, policyIdentifier)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: certificate policy not found: %s", base.ErrResponseStatusNotFound, policyIdentifier.String())
		}
		return err
	}

	nsCtx := ns.GetNSContext(c)
	keyStoreName := GetKeyStoreName(nsCtx.Kind(), nsCtx.Identifier(), policyDoc.ID)

	c, armRAClient, err := s.WithDelegatedARMAuthRoleAssignmentsClient(c)
	if err != nil {
		return err
	}

	assignee := nsCtx.Identifier().UUID()
	// if params.PrincipalID != nil {
	// 	assignee = *params.PrincipalID
	// }
	filterParam := fmt.Sprintf("assignedTo('{%s}')", assignee.String())
	subscriptionIDBuilder := &cloudutils.AzureSubscriptionResourceIDBuilder{
		SubscriptionID: s.GetAzSubscriptionID(),
	}
	scope := subscriptionIDBuilder.WithResourceGroup(s.GetResourceGroupName()).WithKeyVault(s.GetKeyVaultName(), string(kvCategory), keyStoreName).Build()
	log.Debug().Msgf("Lookup role assignments for scope: %s", scope)
	pager := armRAClient.NewListForScopePager(
		scope,
		&armauthorization.RoleAssignmentsClientListForScopeOptions{
			Filter: &filterParam,
		},
	)

	itemsPager := utils.NewMappedPager[[]*armauthorization.RoleAssignment, armauthorization.RoleAssignmentsClientListForScopeResponse](
		utils.NewPagerWithContext(pager, c),
		func(resp armauthorization.RoleAssignmentsClientListForScopeResponse) []*armauthorization.RoleAssignment {
			return resp.Value
		})
	resultPager := utils.NewSerializableItemsPager(
		utils.NewMappedItemsPager[*base.AzureRoleAssignment, *armauthorization.RoleAssignment](
			itemsPager, func(item *armauthorization.RoleAssignment) *base.AzureRoleAssignment {
				if item == nil {
					return nil
				}
				return &base.AzureRoleAssignment{
					ID:               item.ID,
					Name:             item.Name,
					PrincipalId:      item.Properties.PrincipalID,
					RoleDefinitionId: item.Properties.RoleDefinitionID,
				}
			}))

	return api.RespondPagerList(c, resultPager)
}
