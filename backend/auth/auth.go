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
const msClientPrincipalDeviceId string = "MsClientPrincipalDeviceId"

const msClientPrincipalId string = "MsClientPrincipalId"
const msClientPrincipalName string = "MsClientPrincipalName"

const msClientPrincipalClaimType_DeviceID string = "http://schemas.microsoft.com/2012/01/devicecontext/claims/identifier"

func HandleAadAuthMiddleware(ctx *gin.Context) {
	// Intercept the headers here
	var err error
	var decodedClaims []byte
	p := msClientPrincipal{}
	ctx.Set(msClientPrincipalName, ctx.Request.Header.Get("X-Ms-Client-Principal-Name"))
	callerIdStr := ctx.Request.Header.Get("X-Ms-Client-Principal-Id")
	if parsedCallerId, err := uuid.Parse(callerIdStr); err == nil {
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
		if c.Type == msClientPrincipalClaimType_DeviceID {
			if deviceID, err := uuid.Parse(c.Value); err == nil {
				ctx.Set(msClientPrincipalDeviceId, deviceID)
			}
		}
	}

SkipClaims:
	ctx.Next()
}

func CallerPrincipalId(c *gin.Context) uuid.UUID {
	if value, ok := c.Value(msClientPrincipalId).(uuid.UUID); ok {
		return value
	}
	return uuid.Nil
}

func CallerPrincipalName(c *gin.Context) string {
	if value, ok := c.Value(msClientPrincipalName).(string); ok {
		return value
	}
	return ""
}

func CallerPrincipalHasAdminRole(ctx *gin.Context) bool {
	return ctx.Value(msClientPrincipalHasAdminRole) == true
}

func CallerPrincipalDeviceID(c *gin.Context) uuid.UUID {
	if value, ok := c.Value(msClientPrincipalDeviceId).(uuid.UUID); ok {
		return value
	}
	return uuid.Nil
}
