package auth

import (
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

func HandleAadAuthMiddleware(ctx *gin.Context) {
	a := authIdentity{
		appRoles: make(map[string]bool),
	}
	// Intercept the headers here
	var err error
	var decodedClaims []byte
	a.msClientPrincipalID, _ = uuid.Parse(ctx.Request.Header.Get("X-Ms-Client-Principal-Id"))
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
	{
		p := msClientPrincipal{}
		if err = json.Unmarshal(decodedClaims, &p); err != nil {
			log.Warn().Msgf("Error unmarshal X-Ms-Client-Principal header: %s", encodedPrincipal)
			goto afterParsePrincipalClaims
		} else {
			for _, c := range p.Claims {
				switch c.Type {
				case "appid":
					a.appIDClaim, _ = uuid.Parse(c.Value)
				case "roles":
					a.appRoles[c.Value] = true
				}
			}
		}
	}
afterParsePrincipalClaims:
	ctx.Set(appAuthIdentityContextKey, a)
	ctx.Next()
}
