package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ContextKey string

const (
	appAuthContextKey string = "smallkms.appAuthContext"
	roleKeyAppAdmin   string = "App.Admin"
)
const (
	authIdentityContextKey ContextKey = "appIdentity"
)

type AuthIdentity interface {
	HasAdminRole() bool
	ClientPrincipalID() uuid.UUID
	ClientPrincipalName() string
}

type authIdentity struct {
	msClientPrincipalIDstr string
	msClientPrincipalID    uuid.UUID
	msClientPrincipalName  string
	appRoles               map[string]bool
}

func SetAuthContext(c *gin.Context, ctx context.Context) {
	c.Set(appAuthContextKey, ctx)
}

func GetAuthContext(c *gin.Context) (context.Context, bool) {
	if ctx, ok := c.Value(appAuthContextKey).(context.Context); ok {
		return ctx, ok
	}
	return nil, false
}

func GetAuthIdentity(c *gin.Context) (AuthIdentity, bool) {
	if ctx, ok := GetAuthContext(c); ok {
		identity, ok := ctx.Value(authIdentityContextKey).(AuthIdentity)
		return identity, ok
	}
	return nil, false
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
