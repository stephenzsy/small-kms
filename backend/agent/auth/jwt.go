package agentauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
)

func NewSignedAgentAuthJWT(signingMethod jwt.SigningMethod, subject string, key cloudkey.CloudSignatureKey) (string, error) {
	nbf := time.Now()
	token := jwt.NewWithClaims(signingMethod, jwt.RegisteredClaims{
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(nbf),
		ExpiresAt: jwt.NewNumericDate(nbf.Add(time.Hour)),
		ID:        uuid.New().String(),
	})
	token.Header["kid"] = key.KeyID()
	return token.SignedString(key)
}
