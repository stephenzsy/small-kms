package cert

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/base"
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
	policyDoc, err := readCertPolicyDoc(c, policyIdentifier)
	if err != nil {
		if errors.Is(err, base.ErrAzCosmosDocNotFound) {
			return fmt.Errorf("%w: certificate policy not found: %s", base.ErrResponseStatusNotFound, policyIdentifier.String())
		}
		return err
	}

	nsCtx := ns.GetNSContext(c)
	keyStoreName := getKeyStoreName(nsCtx.Kind(), nsCtx.Identifier(), policyDoc)

	c, armRAClient, err := s.WithDelegatedARMAuthRoleAssignmentsClient(c)
	if err != nil {
		return err
	}

	assignee := nsCtx.Identifier().UUID()
	// if params.PrincipalID != nil {
	// 	assignee = *params.PrincipalID
	// }
	filterParam := fmt.Sprintf("assignedTo('{%s}')", assignee.String())
	scope := s.GetKeyvaultCertificateResourceScopeID(keyStoreName, string(kvCategory))
	log.Debug().Msgf("Lookup role assignments for scope: %s", scope)
	pager := armRAClient.NewListForScopePager(
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
	allItems, err := utils.PagerAllItems[*base.AzureRoleAssignment](
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
			}), c)
	if err != nil {
		return err
	}
	if allItems == nil {
		allItems = make([]*base.AzureRoleAssignment, 0)
	}
	return c.JSON(http.StatusOK, allItems)
}
