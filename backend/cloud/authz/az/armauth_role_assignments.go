package cloudauthzaz

import (
	"context"
	"fmt"
	"slices"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stephenzsy/small-kms/backend/utils"
)

func ListRoleAssignments(c context.Context, client *armauthorization.RoleAssignmentsClient, scope string, assignedTo uuid.UUID) utils.ItemsPager[*armauthorization.RoleAssignment] {

	filterParam := fmt.Sprintf("assignedTo('{%s}')", assignedTo.String())
	log.Ctx(c).Debug().Msgf("Lookup role assignments for scope: %s", scope)
	pager := client.NewListForScopePager(
		scope,
		&armauthorization.RoleAssignmentsClientListForScopeOptions{
			Filter: &filterParam,
		},
	)

	return utils.NewMappedPager[[]*armauthorization.RoleAssignment, armauthorization.RoleAssignmentsClientListForScopeResponse](
		utils.NewPagerWithContext(pager, c),
		func(resp armauthorization.RoleAssignmentsClientListForScopeResponse) []*armauthorization.RoleAssignment {
			return resp.Value
		})
}

type RoleAssignmentProvisioner struct {
	AssignedTo       uuid.UUID
	Scope            string
	RoleDefinitionID uuid.UUID
}

func (p *RoleAssignmentProvisioner) IsRoleAssigned(c context.Context, client *armauthorization.RoleAssignmentsClient, roleDefinitionResourceID string) (bool, error) {
	filterParam := fmt.Sprintf("assignedTo('{%s}')", p.AssignedTo.String())
	log.Ctx(c).Debug().Msgf("Lookup role assignments: ID: %s, scope: %s", p.AssignedTo.String(), p.Scope)
	pager := client.NewListForScopePager(
		p.Scope,
		&armauthorization.RoleAssignmentsClientListForScopeOptions{
			Filter: &filterParam,
		},
	)

	allItems, err := utils.PagerToSlice(utils.NewMappedPager[[]*armauthorization.RoleAssignment, armauthorization.RoleAssignmentsClientListForScopeResponse](
		utils.NewPagerWithContext(pager, c),
		func(resp armauthorization.RoleAssignmentsClientListForScopeResponse) []*armauthorization.RoleAssignment {
			return resp.Value
		}))
	if err != nil {
		return false, err
	}
	return slices.ContainsFunc(allItems, func(item *armauthorization.RoleAssignment) bool {
		return *item.Properties.RoleDefinitionID == roleDefinitionResourceID
	}), nil
}

func (p *RoleAssignmentProvisioner) AssignRole(c context.Context, client *armauthorization.RoleAssignmentsClient, roleDefinitionResourceID string) error {

	log.Ctx(c).Debug().Msgf("Create role assignments: ID: %s, scope: %s, definition: %s", p.AssignedTo.String(), p.Scope, roleDefinitionResourceID)
	_, err := client.Create(
		c,
		p.Scope,
		uuid.NewString(),
		armauthorization.RoleAssignmentCreateParameters{
			Properties: &armauthorization.RoleAssignmentProperties{
				PrincipalID:      to.Ptr(p.AssignedTo.String()),
				RoleDefinitionID: &roleDefinitionResourceID,
			},
		},
		nil)
	return err
}
