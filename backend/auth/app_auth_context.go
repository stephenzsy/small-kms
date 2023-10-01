package auth

import (
	"context"

	"github.com/google/uuid"
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
}

type authIdentity struct {
	msClientPrincipalID   uuid.UUID
	msClientPrincipalName string
	appRoles              map[string]bool
	appIDClaim            uuid.UUID
}

func GetAuthIdentity(ctx context.Context) (identity AuthIdentity, ok bool) {
	identity, ok = ctx.Value(appAuthIdentityContextKey).(AuthIdentity)
	return
}

func (i *authIdentity) HasAdminRole() bool {
	if i == nil {
		return false
	}
	return i.appRoles[roleKeyAppAdmin]
}

// use our own copy in case some code path modified it by accident
var uuidNil = uuid.UUID{}

func (i *authIdentity) ClientPrincipalID() uuid.UUID {

	if i == nil {
		return uuidNil
	}

	return i.msClientPrincipalID
}

func (i *authIdentity) ClientPrincipalName() string {

	if i == nil {
		return ""
	}

	return i.msClientPrincipalName
}

func (i *authIdentity) AppIDClaim() uuid.UUID {
	if i == nil {
		return uuidNil
	}
	return i.appIDClaim
}
