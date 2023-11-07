package auth

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
)

func PreconfiguredKeysJWTAuthorization(keys []cloudkey.JsonWebSignatureKey, aud string) echo.MiddlewareFunc {

	keyMapping := make(map[string]*cloudkey.JsonWebSignatureKey, len(keys))
	for _, key := range keys {
		keyMapping[key.KeyID] = &key
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if k, ok := keyMapping[token.Header["kid"].(string)]; ok {
			return k.PublicKey(), nil
		}
		return nil, jwt.ErrInvalidKey
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.String(http.StatusUnauthorized, "missing authorization header")
			}
			authToken := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.ParseWithClaims(authToken, &jwt.RegisteredClaims{}, keyFunc, jwt.WithAudience(aud))
			if err != nil || !token.Valid {
				return c.String(http.StatusUnauthorized, "invalid authorization token")
			}

			return next(ctx.EchoContextWithValue(c, jwtClaimsContextKey, token.Claims, false))
		}
	}
}
