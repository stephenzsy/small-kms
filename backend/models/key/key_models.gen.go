// Package keymodels provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package keymodels

import (
	cloudkey "github.com/stephenzsy/small-kms/backend/cloud/key"
	externalRef0 "github.com/stephenzsy/small-kms/backend/models"
)

// Defines values for KeyStatus.
const (
	KeyStatusActive   KeyStatus = "active"
	KeyStatusInactive KeyStatus = "inactive"
)

// CreateKeyPolicyRequest defines model for CreateKeyPolicyRequest.
type CreateKeyPolicyRequest struct {
	DisplayName string `json:"displayName,omitempty"`
	ExpiryTime  string `json:"expiryTime,omitempty"`
	Exportable  *bool  `json:"exportable,omitempty"`

	// KeySpec these attributes should mostly confirm to JWK (RFC7517)
	KeySpec *JsonWebKeySpec `json:"keySpec,omitempty"`
}

// JsonWebKey defines model for JsonWebKey.
type JsonWebKey = cloudkey.JsonWebKey

// JsonWebKeyCurveName defines model for JsonWebKeyCurveName.
type JsonWebKeyCurveName = cloudkey.JsonWebKeyCurveName

// JsonWebKeyOperation defines model for JsonWebKeyOperation.
type JsonWebKeyOperation = cloudkey.JsonWebKeyOperation

// JsonWebKeySpec these attributes should mostly confirm to JWK (RFC7517)
type JsonWebKeySpec struct {
	Alg           string                `json:"alg,omitempty"`
	Crv           JsonWebKeyCurveName   `json:"crv,omitempty"`
	Extractable   *bool                 `json:"ext,omitempty"`
	KeyOperations []JsonWebKeyOperation `json:"key_ops,omitempty"`
	KeySize       *int                  `json:"key_size,omitempty"`
	Kty           JsonWebKeyType        `json:"kty,omitempty"`
}

// JsonWebKeyType defines model for JsonWebKeyType.
type JsonWebKeyType = cloudkey.JsonWebKeyType

// JsonWebSignatureAlgorithm defines model for JsonWebSignatureAlgorithm.
type JsonWebSignatureAlgorithm = cloudkey.JsonWebSignatureAlgorithm

// Key defines model for Key.
type Key = keyComposed

// KeyFields defines model for KeyFields.
type KeyFields struct {
	Identififier string                    `json:"identififier"`
	Jwk          JsonWebKey                `json:"jwk"`
	Nbf          *externalRef0.NumericDate `json:"nbf,omitempty"`

	// Sid Key Vault Secret ID
	KeyVaultSecretID string `json:"sid,omitempty"`
}

// KeyPolicy defines model for KeyPolicy.
type KeyPolicy = keyPolicyComposed

// KeyPolicyFields defines model for KeyPolicyFields.
type KeyPolicyFields struct {
	ExpiryTime string `json:"expiryTime,omitempty"`

	// KeySpec these attributes should mostly confirm to JWK (RFC7517)
	KeySpec JsonWebKeySpec `json:"keySpec"`
}

// KeyRef defines model for KeyRef.
type KeyRef = keyRefComposed

// KeyRefFields defines model for KeyRefFields.
type KeyRefFields struct {
	Exp              *externalRef0.NumericDate `json:"exp,omitempty"`
	Iat              externalRef0.NumericDate  `json:"iat"`
	PolicyIdentifier string                    `json:"policyIdentifier"`
	Status           KeyStatus                 `json:"status"`
}

// KeyStatus defines model for KeyStatus.
type KeyStatus string

// OneTimeKey OneTimeKey
type OneTimeKey struct {
	Exp externalRef0.NumericDate `json:"exp"`
	Iat externalRef0.NumericDate `json:"iat"`
	Jwk JsonWebKey               `json:"jwk"`
}

// KeyPolicyResponse defines model for KeyPolicyResponse.
type KeyPolicyResponse = KeyPolicy

// KeyRefsResponse defines model for KeyRefsResponse.
type KeyRefsResponse = []KeyRef

// KeyResponse defines model for KeyResponse.
type KeyResponse = Key
