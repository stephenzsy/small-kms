package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type msClientPrincipalClaims struct {
	Type  string `json:"typ"`
	Value string `json:"val"`
}

type msClientPrincipal struct {
	Claims []msClientPrincipalClaims `json:"claims"`
}

const msClientPrincipalHasAdminRole string = "MsClientPrincipalHasAdminRole"
const msClientPrincipalDeviceId string = "MsClientPrincipalDeviceId"

const msClientPrincipalName string = "MsClientPrincipalName"

const msClientPrincipalClaimType_DeviceID string = "http://schemas.microsoft.com/2012/01/devicecontext/claims/identifier"

func HandleAadAuthMiddleware(ctx *gin.Context) {
	a := authIdentity{
		appRoles: make(map[string]bool),
	}
	// Intercept the headers here
	var err error
	var decodedClaims []byte
	p := msClientPrincipal{}
	a.msClientPrincipalIDstr = ctx.Request.Header.Get("X-Ms-Client-Principal-Id")
	if parsedCallerId, err := uuid.Parse(a.msClientPrincipalIDstr); err == nil {
		a.msClientPrincipalID = parsedCallerId
	}

	a.msClientPrincipalName = ctx.Request.Header.Get("X-Ms-Client-Principal-Name")

	encodedPrincipal := ctx.Request.Header.Get("X-Ms-Client-Principal")
	if len(encodedPrincipal) == 0 {
		log.Warn().Msg("No X-Ms-Client-Principal header found")
		goto afterParsePrincipalClaims
	}
	decodedClaims, err = base64.StdEncoding.DecodeString(encodedPrincipal)
	if err != nil {
		log.Warn().Msg("Error decoding X-Ms-Client-Principal header")
		goto afterParsePrincipalClaims
	}
	if err = json.Unmarshal(decodedClaims, &p); err != nil {
		log.Warn().Msgf("Error unmarshal X-Ms-Client-Principal header: %s", encodedPrincipal)
		goto afterParsePrincipalClaims
	} else {
		for _, c := range p.Claims {
			if c.Type == "roles" {
				a.appRoles[c.Value] = true
			}
		}
	}

afterParsePrincipalClaims:
	SetAuthContext(ctx, context.WithValue(context.Background(), authIdentityContextKey, a))
	ctx.Next()
}
