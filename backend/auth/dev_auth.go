package auth

import (
	"context"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DevJwtClaims struct {
	jwt.Claims
	AppID      string   `json:"appid"`
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
		claims := DevJwtClaims{}
		parser := jwt.NewParser()
		_, _, err := parser.ParseUnverified(a.bearerToken, &claims)
		if err != nil {
			c.Logger().Warn("invalid Authorization header", err)
			return c.NoContent(401)
		}
		a.appIDClaim, _ = uuid.Parse(claims.AppID)
		a.msClientPrincipalID, _ = uuid.Parse(claims.ObjectID)
		a.msClientPrincipalName = claims.UniqueName
		for _, r := range claims.Roles {
			a.appRoles[r] = true
		}
		ctx := c.Request().Context()
		ctx = context.WithValue(ctx, appAuthIdentityContextKey, a)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
