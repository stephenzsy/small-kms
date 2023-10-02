package auth

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
)

type ContextKey string

const (
	appAuthIdentityContextKey string = "smallkms.appAuthIdentity"
	roleKeyAppAdmin           string = "App.Admin"
)

type AuthIdentity interface {
	HasAdminRole() bool
	ClientPrincipalID() uuid.UUID
	ClientPrincipalName() string
	AppIDClaim() uuid.UUID
	GetOnBehalfOfTokenCredential(s common.CommonConfig, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error)
}

type authIdentity struct {
	msClientPrincipalID   uuid.UUID
	msClientPrincipalName string
	appRoles              map[string]bool
	appIDClaim            uuid.UUID
	bearerToken           string
}

func GetAuthIdentity(ctx context.Context) (identity AuthIdentity, ok bool) {
	ctxValue := ctx.Value(appAuthIdentityContextKey)
	identity, ok = ctxValue.(AuthIdentity)
	return
}

func (i authIdentity) HasAdminRole() bool {
	return i.appRoles[roleKeyAppAdmin]
}

func (i authIdentity) ClientPrincipalID() uuid.UUID {
	return i.msClientPrincipalID
}

func (i authIdentity) ClientPrincipalName() string {
	return i.msClientPrincipalName
}

func (i authIdentity) AppIDClaim() uuid.UUID {
	return i.appIDClaim
}

func (i authIdentity) GetOnBehalfOfTokenCredential(s common.CommonConfig, opts *azidentity.OnBehalfOfCredentialOptions) (azcore.TokenCredential, error) {
	return s.NewOnBehalfOfCredential(i.bearerToken, opts)
}
