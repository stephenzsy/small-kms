package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type msClientPrincipalClaims struct {
	Type  string `json:"typ"`
	Value string `json:"val"`
}

type msClientPrincipal struct {
	Claims []msClientPrincipalClaims `json:"claims"`
}

func ProxiedAADAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		a := authIdentity{
			appRoles: make(map[string]bool),
		}

		// Intercept the headers here
		var err error
		var decodedClaims []byte
		headers := c.Request().Header
		a.msClientPrincipalID, err = uuid.Parse(headers.Get("X-Ms-Client-Principal-Id"))
		if err != nil {
			c.Logger().Errorf("Error parsing X-Ms-Client-Principal-Id header: %s", err.Error())
			return c.NoContent(http.StatusUnauthorized)
		}
		a.msClientPrincipalName = headers.Get("X-Ms-Client-Principal-Name")
		a.bearerToken = headers.Get("Authorization")[7:]
		encodedPrincipal := headers.Get("X-Ms-Client-Principal")
		if len(encodedPrincipal) == 0 {
			c.Logger().Warn("No X-Ms-Client-Principal header found")
			goto afterParsePrincipalClaims
		}
		decodedClaims, err = base64.StdEncoding.DecodeString(encodedPrincipal)
		if err != nil {
			c.Logger().Warn("Error decoding X-Ms-Client-Principal header")
			goto afterParsePrincipalClaims
		}
		{
			p := msClientPrincipal{}
			if err = json.Unmarshal(decodedClaims, &p); err != nil {
				c.Logger().Warnf("Error unmarshal X-Ms-Client-Principal header: %s", encodedPrincipal)
				goto afterParsePrincipalClaims
			} else {
				for _, c := range p.Claims {
					switch c.Type {
					case "roles":
						a.appRoles[c.Value] = true
					}
				}
			}
		}
	afterParsePrincipalClaims:
		return next(ctx.EchoContextWithValue(c, authIdentityContextKey, &a, false))
	}
}

var _ echo.MiddlewareFunc = ProxiedAADAuth
