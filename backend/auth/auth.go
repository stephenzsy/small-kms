package auth

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type msClientPrincipalClaims struct {
	Type  string `json:"typ"`
	Value string `json:"val"`
}

type msClientPrincipal struct {
	Claims []msClientPrincipalClaims `json:"claims"`
}

type ContextKey string

const msClientPrincipalHasAdminRole string = "MsClientPrincipalHasAdminRole"
const msClientPrincipalId string = "MsClientPrincipalId"

func HandleAadAuthMiddleware(ctx *gin.Context) {
	// Intercept the headers here
	var err error
	var decodedClaims []byte
	p := msClientPrincipal{}
	callerIdStr := ctx.Request.Header.Get("X-Ms-Client-Principal-Id")
	if parsedCallerId, err := uuid.Parse(callerIdStr); err != nil {
		ctx.Set(msClientPrincipalId, parsedCallerId)
	}
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
			ctx.Set(msClientPrincipalHasAdminRole, true)
		}
	}

SkipClaims:
	ctx.Next()
}

func CallerPrincipalHasAdminRole(ctx *gin.Context) bool {
	return ctx.Value(msClientPrincipalHasAdminRole) == true
}

func CallerPrincipalId(c *gin.Context) uuid.UUID {
	return c.Value(msClientPrincipalId).(uuid.UUID)
}
