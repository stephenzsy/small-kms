// Package agentmodels provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package agentmodels

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/oapi-codegen/runtime"
	externalRef0 "github.com/stephenzsy/small-kms/backend/models"
)

// Defines values for AgentConfigName.
const (
	AgentConfigNameIdentity AgentConfigName = "identity"
)

// Agent defines model for Agent.
type Agent = externalRef0.Profile

// AgentConfig defines model for AgentConfig.
type AgentConfig struct {
	union json.RawMessage
}

// AgentConfigBundle defines model for AgentConfigBundle.
type AgentConfigBundle struct {
	EnvGuards []string          `json:"envGuards"`
	Expires   *time.Time        `json:"expires,omitempty"`
	Id        string            `json:"id"`
	Items     *[]AgentConfigRef `json:"items,omitempty"`
}

// AgentConfigIdentity defines model for AgentConfigIdentity.
type AgentConfigIdentity = agentConfigIdentityComposed

// AgentConfigIdentityFields defines model for AgentConfigIdentityFields.
type AgentConfigIdentityFields struct {
	KeyCredentialCertificatePolicyId string `json:"keyCredentialCertificatePolicyId"`
}

// AgentConfigName defines model for AgentConfigName.
type AgentConfigName string

// AgentConfigRef defines model for AgentConfigRef.
type AgentConfigRef struct {
	Name    AgentConfigName `json:"name"`
	Updated time.Time       `json:"updated"`
	Version string          `json:"version"`
}

// CreateAgentConfigRequest defines model for CreateAgentConfigRequest.
type CreateAgentConfigRequest struct {
	union json.RawMessage
}

// CreateAgentRequest defines model for CreateAgentRequest.
type CreateAgentRequest struct {
	// AppId The Application ID (Client ID) of the agent
	AppId string `json:"appId,omitempty"`

	// DisplayName The display name of the agent application
	DisplayName string `json:"displayName,omitempty"`
}

// AgentConfigResponse defines model for AgentConfigResponse.
type AgentConfigResponse = AgentConfig

// AgentResposne defines model for AgentResposne.
type AgentResposne = Agent

// AsAgentConfigIdentity returns the union data inside the AgentConfig as a AgentConfigIdentity
func (t AgentConfig) AsAgentConfigIdentity() (AgentConfigIdentity, error) {
	var body AgentConfigIdentity
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAgentConfigIdentity overwrites any union data inside the AgentConfig as the provided AgentConfigIdentity
func (t *AgentConfig) FromAgentConfigIdentity(v AgentConfigIdentity) error {
	v.Name = "identity"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAgentConfigIdentity performs a merge with any union data inside the AgentConfig, using the provided AgentConfigIdentity
func (t *AgentConfig) MergeAgentConfigIdentity(v AgentConfigIdentity) error {
	v.Name = "identity"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t AgentConfig) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"name"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t AgentConfig) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "identity":
		return t.AsAgentConfigIdentity()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t AgentConfig) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *AgentConfig) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// AsAgentConfigIdentityFields returns the union data inside the CreateAgentConfigRequest as a AgentConfigIdentityFields
func (t CreateAgentConfigRequest) AsAgentConfigIdentityFields() (AgentConfigIdentityFields, error) {
	var body AgentConfigIdentityFields
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromAgentConfigIdentityFields overwrites any union data inside the CreateAgentConfigRequest as the provided AgentConfigIdentityFields
func (t *CreateAgentConfigRequest) FromAgentConfigIdentityFields(v AgentConfigIdentityFields) error {
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeAgentConfigIdentityFields performs a merge with any union data inside the CreateAgentConfigRequest, using the provided AgentConfigIdentityFields
func (t *CreateAgentConfigRequest) MergeAgentConfigIdentityFields(v AgentConfigIdentityFields) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t CreateAgentConfigRequest) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *CreateAgentConfigRequest) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}