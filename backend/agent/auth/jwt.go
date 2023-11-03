package agentauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
)

func NewSignedAgentAuthJWT(signingMethod jwt.SigningMethod, subject, endpoint string, key cloudkey.CloudSignatureKey) (string, error) {
	iat := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   subject,
		Audience:  jwt.ClaimStrings{endpoint},
		IssuedAt:  jwt.NewNumericDate(iat),
		ExpiresAt: jwt.NewNumericDate(iat.Add(time.Hour)),
		ID:        uuid.New().String(),
	}
	token := jwt.NewWithClaims(signingMethod, &claims)
	token.Header["kid"] = key.KeyID()
	return token.SignedString(key)
}
