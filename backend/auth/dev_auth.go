package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type DevJwtClaims struct {
	jwt.Claims
	AppID      string   `json:"appid"`
	ObjectID   string   `json:"oid"`
	UniqueName string   `json:"unique_name"`
	Roles      []string `json:"roles,omitempty"`
}

func HandleDevJWTMiddleware(ctx *gin.Context) {
	authHeader := ctx.Request.Header.Get("Authorization")
	if len(authHeader) == 0 {
		log.Warn().Msg("No Authorization header found")
		ctx.JSON(401, gin.H{"error": "No Authorization header"})
		return
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		log.Warn().Msg("Invalid Authorization header")
		ctx.JSON(401, gin.H{"error": "Invalid Authorization header"})
		return
	}
	encodedJwt := strings.TrimPrefix(authHeader, "Bearer ")
	claims := DevJwtClaims{}
	parser := jwt.NewParser()
	_, _, err := parser.ParseUnverified(encodedJwt, &claims)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid Authorization header")
		ctx.JSON(401, gin.H{"error": "Invalid Authorization header"})
		return
	}
	a := authIdentity{
		appRoles: make(map[string]bool),
	}
	a.appIDClaim, _ = uuid.Parse(claims.AppID)
	a.msClientPrincipalID, _ = uuid.Parse(claims.ObjectID)
	a.msClientPrincipalName = claims.UniqueName
	for _, r := range claims.Roles {
		a.appRoles[r] = true
	}
	ctx.Set(appAuthIdentityContextKey, a)
	ctx.Next()
}
