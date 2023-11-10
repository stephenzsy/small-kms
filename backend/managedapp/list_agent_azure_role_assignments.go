package managedapp

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/labstack/echo/v4"
	"github.com/stephenzsy/small-kms/backend/api"
	"github.com/stephenzsy/small-kms/backend/base"
	"github.com/stephenzsy/small-kms/backend/cert"
	cloudauthzaz "github.com/stephenzsy/small-kms/backend/cloud/authz/az"
	"github.com/stephenzsy/small-kms/backend/cloud/containerregistry/acr"
	"github.com/stephenzsy/small-kms/backend/cloudutils"
	"github.com/stephenzsy/small-kms/backend/internal/authz"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/secret"
	"github.com/stephenzsy/small-kms/backend/utils"
)

// ListAgentAzureRoleAssignments implements ServerInterface.
func (s *server) ListAgentAzureRoleAssignments(ec echo.Context, namespaceKind base.NamespaceKind, namespaceId base.ID, configName AgentConfigName) (err error) {
	c := ec.(ctx.RequestContext)

	if !authz.AuthorizeAdminOnly(c) {
		return s.RespondRequireAdmin(c)
	}
	c = ns.WithNSContext(c, namespaceKind, namespaceId)

	var pagers []utils.ItemsPager[*armauthorization.RoleAssignment]
	switch configName {
	case AgentConfigNameServer,
		AgentConfigNameRadius:
		pagers, err = s.getAzureRoleAssignmentPagers(c, configName)
		if err != nil {
			return nil
		}
	default:
		return fmt.Errorf("%w: invalid agent config name", base.ErrResponseStatusBadRequest)
	}

	resultPager := utils.NewSerializableItemsPager(
		utils.NewMappedItemsPager[*base.AzureRoleAssignment, *armauthorization.RoleAssignment](utils.NewChainedItemPagers(pagers...), func(ra *armauthorization.RoleAssignment) *base.AzureRoleAssignment {
			if ra == nil {
				return nil
			}
			return &base.AzureRoleAssignment{
				ID:               ra.ID,
				Name:             ra.Name,
				RoleDefinitionId: ra.Properties.RoleDefinitionID,
				PrincipalId:      ra.Properties.PrincipalID,
			}
		}))

	return api.RespondPagerList(c, resultPager)
}

func (s *server) getAzureRoleAssignmentPagers(c ctx.RequestContext, configName AgentConfigName) ([]utils.ItemsPager[*armauthorization.RoleAssignment], error) {
	nsCtx := ns.GetNSContext(c)
	if _, ok := nsCtx.ID().AsUUID(); !ok {
		return nil, fmt.Errorf("%w: invalid namespace identifier", base.ErrResponseStatusBadRequest)
	}
	assignedTo := nsCtx.ID().UUID()
	c, armRAClient, err := s.WithDelegatedARMAuthRoleAssignmentsClient(c)
	if err != nil {
		return nil, err
	}

	subscriptionIDBuilder := &cloudutils.AzureSubscriptionResourceIDBuilder{
		SubscriptionID: s.GetAzSubscriptionID(),
	}

	switch configName {
	case AgentConfigNameServer:
		doc, err := ApiReadAgentConfigDoc(c)
		if err != nil {
			return nil, err
		}
		if doc.GlobalACRImageRef == "" {
			return nil, fmt.Errorf("%w: image ref is not specified", base.ErrResponseStatusBadRequest)
		}

		pagers := make([]utils.ItemsPager[*armauthorization.RoleAssignment], 0, 2)
		// ACR Pull
		{
			acrName, err := acr.ExtractACRName(doc.GlobalACRImageRef)
			if err != nil {
				return nil, err
			}
			if acrResourceGroupName := s.EnvService().Default("AZURE_RESOURCE_GROUP_NAME", "", "ACR_"); acrResourceGroupName != "" {
				scope := subscriptionIDBuilder.WithResourceGroup(acrResourceGroupName).WithContainerRegistry(acrName).Build()
				pagers = append(pagers, cloudauthzaz.ListRoleAssignments(c, armRAClient, scope, assignedTo))
			}
		}
		// Key Vault Secrets User
		{
			scope := subscriptionIDBuilder.WithResourceGroup(s.GetResourceGroupName()).WithKeyVault(s.GetKeyVaultName(), "secrets",
				cert.GetKeyStoreName(nsCtx.Kind(), nsCtx.ID(), doc.TLSCertificatePolicyID)).Build()
			pagers = append(pagers, cloudauthzaz.ListRoleAssignments(c, armRAClient, scope, assignedTo))
		}
		return pagers, nil
	case AgentConfigNameRadius:
		doc, err := apiReadAgentConfigRadiusDoc(c)
		if err != nil {
			return nil, err
		}
		if doc.GlobalRadiusServerACRImageRef == "" {
			return nil, fmt.Errorf("%w: image ref is not specified", base.ErrResponseStatusBadRequest)
		}

		pagers := make([]utils.ItemsPager[*armauthorization.RoleAssignment], 0, 2)
		// ACR Pull
		{
			acrName, err := acr.ExtractACRName(doc.GlobalRadiusServerACRImageRef)
			if err != nil {
				return nil, err
			}
			if acrResourceGroupName := s.EnvService().Default("AZURE_RESOURCE_GROUP_NAME", "", "ACR_"); acrResourceGroupName != "" {
				scope := subscriptionIDBuilder.WithResourceGroup(acrResourceGroupName).WithContainerRegistry(acrName).Build()
				pagers = append(pagers, cloudauthzaz.ListRoleAssignments(c, armRAClient, scope, assignedTo))
			}
		}
		// Key Vault Secrets User
		{
			for _, client := range doc.Clients {
				scope := subscriptionIDBuilder.WithResourceGroup(s.GetResourceGroupName()).WithKeyVault(s.GetKeyVaultName(), "secrets",
					secret.GetKeyStoreName(nsCtx.Kind(), nsCtx.ID(), client.SecretPolicyId)).Build()
				pagers = append(pagers, cloudauthzaz.ListRoleAssignments(c, armRAClient, scope, assignedTo))
			}
		}
		return pagers, nil
	}
	return nil, nil
}
