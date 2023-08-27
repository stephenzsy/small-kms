package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type msClientPrincipalClaims struct {
	Type  string `json:"typ"`
	Value string `json:"val"`
}

type msClientPrincipal struct {
	Claims []msClientPrincipalClaims `json:"claims"`
}

type ContextKey string

const HasAdminAppRoleContextKey ContextKey = "HasAdminAppRole"

func HandleAadAuthMiddleware(ctx *gin.Context) {
	// Intercept the headers here
	var err error
	var decodedClaims []byte
	p := msClientPrincipal{}
	encodedPrincipal := ctx.Request.Header.Get("X-Ms-Client-Principal")
	if len(encodedPrincipal) == 0 {
		goto SkipClaims
	}
	decodedClaims, err = base64.StdEncoding.DecodeString(encodedPrincipal)
	if err != nil {
		goto SkipClaims
	}
	err = json.Unmarshal(decodedClaims, &p)
	if err != nil {
		goto SkipClaims
	}
	for _, c := range p.Claims {
		if c.Type == "roles" && c.Value == "App.Admin" {
			ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), HasAdminAppRoleContextKey, true))
		}
	}

SkipClaims:
	ctx.Next()
}

func HasAdminAppRole(ctx *gin.Context) bool {
	return ctx.Request.Context().Value(HasAdminAppRoleContextKey) == true
}

func GetCallerID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("X-Ms-Client-Principal-Id")
}
