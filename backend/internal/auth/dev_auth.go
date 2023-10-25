package auth

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

type devJwtClaims struct {
	jwt.Claims
	AppID      string   `json:"azp"`
	ObjectID   string   `json:"oid"`
	UniqueName string   `json:"unique_name"`
	Roles      []string `json:"roles,omitempty"`
}

func UnverifiedAADJwtAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if len(authHeader) == 0 {
			c.Logger().Warn("No Authorization header found")
			return c.NoContent(401)
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Logger().Warn("Invalid Authorization header")
			return c.NoContent(401)
		}
		a := authIdentity{
			appRoles: make(map[string]bool),
		}
		a.bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
		claims := devJwtClaims{}
		parser := jwt.NewParser()
		_, _, err := parser.ParseUnverified(a.bearerToken, &claims)
		if err != nil {
			c.Logger().Warn("invalid Authorization header", err)
			return c.NoContent(401)
		}
		a.msClientPrincipalID, _ = uuid.Parse(claims.ObjectID)
		a.msClientPrincipalName = claims.UniqueName
		for _, r := range claims.Roles {
			a.appRoles[r] = true
		}
		a.appID = claims.AppID
		return next(ctx.EchoContextWithValue(c, authIdentityContextKey, &a, false))
	}
}
