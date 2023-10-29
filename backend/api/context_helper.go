package api

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	"github.com/stephenzsy/small-kms/backend/internal/auth"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type contextKey int

const (
	delegatedARMAuthRoleAssignmentsClient contextKey = iota
)

func (s *apiServer) WithDelegatedARMAuthRoleAssignmentsClient(c ctx.RequestContext) (ctx.RequestContext, *armauthorization.RoleAssignmentsClient, error) {
	return auth.WithDelegatedClient[armauthorization.RoleAssignmentsClient, contextKey](c, delegatedARMAuthRoleAssignmentsClient, func(creds azcore.TokenCredential) (*armauthorization.RoleAssignmentsClient, error) {
		return armauthorization.NewRoleAssignmentsClient(s.azSubscriptionID, creds, nil)
	})
}

func GetDelegatedARMAuthRoleAssignmentsClient(c context.Context) *armauthorization.RoleAssignmentsClient {
	return auth.GetDelegateClient[armauthorization.RoleAssignmentsClient, contextKey](c, delegatedARMAuthRoleAssignmentsClient)
}
