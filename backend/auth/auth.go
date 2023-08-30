package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"

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
const HasScepAppRoleContextKey ContextKey = "HasScepAppRole"

func HandleAadAuthMiddleware(ctx *gin.Context) {
	// Intercept the headers here
	var err error
	var decodedClaims []byte
	p := msClientPrincipal{}
	encodedPrincipal := ctx.Request.Header.Get("X-Ms-Client-Principal")
	if len(encodedPrincipal) == 0 {
		log.Println("No X-Ms-Client-Principal header found")
		goto SkipClaims
	}
	decodedClaims, err = base64.StdEncoding.DecodeString(encodedPrincipal)
	if err != nil {
		log.Println("Error decoding X-Ms-Client-Principal header")
		goto SkipClaims
	}
	err = json.Unmarshal(decodedClaims, &p)
	if err != nil {
		log.Printf("Error unmarshal X-Ms-Client-Principal header: %s", encodedPrincipal)
		goto SkipClaims
	}
	for _, c := range p.Claims {
		if c.Type == "roles" && c.Value == "App.Admin" {
			ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), HasAdminAppRoleContextKey, true))
		} else if c.Type == "roles" && c.Value == "App.Scep" {
			ctx.Request = ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), HasScepAppRoleContextKey, true))
		}
	}

SkipClaims:
	ctx.Next()
}

func CallerHasAdminAppRole(ctx *gin.Context) bool {
	return ctx.Request.Context().Value(HasAdminAppRoleContextKey) == true
}

func CallerHasScepAppRole(ctx *gin.Context) bool {
	return ctx.Request.Context().Value(HasScepAppRoleContextKey) == true
}

func GetCallerID(ctx *gin.Context) string {
	return ctx.Request.Header.Get("X-Ms-Client-Principal-Id")
}
