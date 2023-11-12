// Package frconfig provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package frconfig

import (
	externalRef0 "github.com/stephenzsy/small-kms/backend/base"
)

// Defines values for RadiusServerListenerType.
const (
	RadiusServerListenerTypeAcct RadiusServerListenerType = "acct"
	RadiusServerListenerTypeAuth RadiusServerListenerType = "auth"
)

// RadiusClientConfig defines model for RadiusClientConfig.
type RadiusClientConfig struct {
	Ipaddr         string          `json:"ipaddr,omitempty"`
	Name           string          `json:"name"`
	Secret         string          `json:"-"`
	SecretId       externalRef0.Id `json:"secretId,omitempty"`
	SecretPolicyId externalRef0.Id `json:"secretPolicyId,omitempty"`
}

// RadiusEapTls defines model for RadiusEapTls.
type RadiusEapTls struct {
	CertId       externalRef0.Id `json:"certId,omitempty"`
	CertPolicyId externalRef0.Id `json:"certPolicyId"`
}

// RadiusServerConfig defines model for RadiusServerConfig.
type RadiusServerConfig struct {
	Listeners []RadiusServerListenConfig `json:"listeners"`
	Name      string                     `json:"name"`
}

// RadiusServerListenConfig defines model for RadiusServerListenConfig.
type RadiusServerListenConfig struct {
	Ipaddr string                   `json:"ipaddr"`
	Port   int                      `json:"port"`
	Type   RadiusServerListenerType `json:"type"`
}

// RadiusServerListenerType defines model for RadiusServerListenerType.
type RadiusServerListenerType string
