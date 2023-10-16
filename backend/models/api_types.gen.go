// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/oapi-codegen/runtime"
	openapi_types "github.com/oapi-codegen/runtime/types"
	externalRef0 "github.com/stephenzsy/small-kms/backend/shared"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for CreateProfileRequestType.
const (
	ProfileTypeManagedApplication CreateProfileRequestType = "managed-application"
)

// Defines values for LinkedCertificateTemplateUsage.
const (
	LinkedCertificateTemplateUsageClientAuthorization       LinkedCertificateTemplateUsage = "cliant-authorization"
	LinkedCertificateTemplateUsageMemberDelegatedEnrollment LinkedCertificateTemplateUsage = "member-delegated-enrollment"
)

// Defines values for PatchServiceConfigParamsConfigPath.
const (
	ServiceConfigPathAppRoleIds             PatchServiceConfigParamsConfigPath = "appRoleIds"
	ServiceConfigPathAzureContainerRegistry PatchServiceConfigParamsConfigPath = "azureContainerRegistry"
	ServiceConfigPathAzureSubscriptionId    PatchServiceConfigParamsConfigPath = "azureSubscriptionId"
	ServiceConfigPathKeyvaultArmResourceId  PatchServiceConfigParamsConfigPath = "keyvaultArmResourceId"
)

// AzureRoleAssignment defines model for AzureRoleAssignment.
type AzureRoleAssignment struct {
	Id               *string `json:"id,omitempty"`
	Name             *string `json:"name,omitempty"`
	PrincipalId      *string `json:"principalId,omitempty"`
	RoleDefinitionId *string `json:"roleDefinitionId,omitempty"`
}

// CertificateLifetimeTrigger defines model for CertificateLifetimeTrigger.
type CertificateLifetimeTrigger struct {
	DaysBeforeExpiry   *int32 `json:"days_before_expiry,omitempty"`
	LifetimePercentage *int32 `json:"lifetime_percentage,omitempty"`
}

// CertificateTemplate defines model for CertificateTemplate.
type CertificateTemplate struct {
	// Deleted Time when the deleted was deleted
	Deleted        *time.Time                   `json:"deleted,omitempty"`
	Id             Identifier                   `json:"id"`
	IssuerTemplate externalRef0.ResourceLocator `json:"issuerTemplate"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties   externalRef0.JwkProperties    `json:"keyProperties"`
	KeyStorePath    *string                       `json:"keyStorePath,omitempty"`
	LifetimeTrigger CertificateLifetimeTrigger    `json:"lifetimeTrigger"`
	LinkTo          *externalRef0.ResourceLocator `json:"linkTo,omitempty"`
	Locator         ResourceLocator               `json:"locator"`
	Metadata        map[string]interface{}        `json:"metadata,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string `json:"subjectCommonName"`

	// Updated Time when the resoruce was last updated
	Updated          *time.Time                      `json:"updated,omitempty"`
	UpdatedBy        *string                         `json:"updatedBy,omitempty"`
	Usages           []externalRef0.CertificateUsage `json:"usages"`
	ValidityInMonths int32                           `json:"validity_months"`
}

// CertificateTemplateFields Certificate fields, may accept template substitutions
type CertificateTemplateFields struct {
	IssuerTemplate externalRef0.ResourceLocator `json:"issuerTemplate"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties    externalRef0.JwkProperties      `json:"keyProperties"`
	KeyStorePath     *string                         `json:"keyStorePath,omitempty"`
	LifetimeTrigger  CertificateLifetimeTrigger      `json:"lifetimeTrigger"`
	Usages           []externalRef0.CertificateUsage `json:"usages"`
	ValidityInMonths int32                           `json:"validity_months"`
}

// CertificateTemplateParameters Certificate fields, may accept template substitutions
type CertificateTemplateParameters struct {
	IssuerTemplate *externalRef0.ResourceLocator `json:"issuerTemplate,omitempty"`
	KeyExportable  *bool                         `json:"keyExportable,omitempty"`

	// KeyProperties Property bag of JSON Web Key (RFC 7517) with additional fields, all bytes are base64url encoded
	KeyProperties   *externalRef0.JwkProperties `json:"keyProperties,omitempty"`
	KeyStorePath    *string                     `json:"keyStorePath,omitempty"`
	LifetimeTrigger *CertificateLifetimeTrigger `json:"lifetimeTrigger,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string                          `json:"subjectCommonName"`
	Usages            []externalRef0.CertificateUsage `json:"usages"`
	ValidityInMonths  *int32                          `json:"validity_months,omitempty"`
}

// CertificateTemplateRef defines model for CertificateTemplateRef.
type CertificateTemplateRef struct {
	// Deleted Time when the deleted was deleted
	Deleted  *time.Time                    `json:"deleted,omitempty"`
	Id       externalRef0.Identifier       `json:"id"`
	LinkTo   *externalRef0.ResourceLocator `json:"linkTo,omitempty"`
	Locator  externalRef0.ResourceLocator  `json:"locator"`
	Metadata map[string]interface{}        `json:"metadata,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string `json:"subjectCommonName"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// CertificateTemplateRefFields defines model for CertificateTemplateRefFields.
type CertificateTemplateRefFields struct {
	LinkTo *externalRef0.ResourceLocator `json:"linkTo,omitempty"`

	// SubjectCommonName Common name
	SubjectCommonName string `json:"subjectCommonName"`
}

// CreateLinkedCertificateTemplateParameters defines model for CreateLinkedCertificateTemplateParameters.
type CreateLinkedCertificateTemplateParameters struct {
	TargetTemplate externalRef0.ResourceLocator   `json:"targetTemplate"`
	Usage          LinkedCertificateTemplateUsage `json:"usage"`
}

// CreateManagedApplicationProfileRequest defines model for CreateManagedApplicationProfileRequest.
type CreateManagedApplicationProfileRequest struct {
	Name string                   `json:"name"`
	Type CreateProfileRequestType `json:"type"`
}

// CreateProfileRequest defines model for CreateProfileRequest.
type CreateProfileRequest struct {
	union json.RawMessage
}

// CreateProfileRequestType defines model for CreateProfileRequestType.
type CreateProfileRequestType string

// LinkedCertificateTemplateUsage defines model for LinkedCertificateTemplateUsage.
type LinkedCertificateTemplateUsage string

// Profile defines model for Profile.
type Profile = ProfileRef

// ProfileRef defines model for ProfileRef.
type ProfileRef struct {
	// Deleted Time when the deleted was deleted
	Deleted *time.Time `json:"deleted,omitempty"`

	// DisplayName Display name of the resource
	DisplayName string                  `json:"displayName"`
	Id          externalRef0.Identifier `json:"id"`

	// IsAppManaged Whether the resource is managed by the application
	IsAppManaged *bool                        `json:"isAppManaged,omitempty"`
	Locator      externalRef0.ResourceLocator `json:"locator"`
	Metadata     map[string]interface{}       `json:"metadata,omitempty"`
	Type         externalRef0.NamespaceKind   `json:"type"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// ProfileRefFields defines model for ProfileRefFields.
type ProfileRefFields struct {
	// DisplayName Display name of the resource
	DisplayName string `json:"displayName"`

	// IsAppManaged Whether the resource is managed by the application
	IsAppManaged *bool                      `json:"isAppManaged,omitempty"`
	Type         externalRef0.NamespaceKind `json:"type"`
}

// ServiceConfig defines model for ServiceConfig.
type ServiceConfig struct {
	AppRoleIds struct {
		AgentActiveHost openapi_types.UUID `json:"Agent.ActiveHost"`
		AppAdmin        openapi_types.UUID `json:"App.Admin"`
	} `json:"appRoleIds"`
	AzureContainerRegistry struct {
		ArmResourceId string `json:"armResourceId"`
		LoginServer   string `json:"loginServer"`
		Name          string `json:"name"`
	} `json:"azureContainerRegistry"`
	AzureSubscriptionId string `json:"azureSubscriptionId"`

	// Deleted Time when the deleted was deleted
	Deleted               *time.Time                   `json:"deleted,omitempty"`
	Id                    externalRef0.Identifier      `json:"id"`
	KeyvaultArmResourceId string                       `json:"keyvaultArmResourceId"`
	Locator               externalRef0.ResourceLocator `json:"locator"`
	Metadata              map[string]interface{}       `json:"metadata,omitempty"`

	// Updated Time when the resoruce was last updated
	Updated   *time.Time `json:"updated,omitempty"`
	UpdatedBy *string    `json:"updatedBy,omitempty"`
}

// ServiceConfigFields defines model for ServiceConfigFields.
type ServiceConfigFields struct {
	AppRoleIds struct {
		AgentActiveHost openapi_types.UUID `json:"Agent.ActiveHost"`
		AppAdmin        openapi_types.UUID `json:"App.Admin"`
	} `json:"appRoleIds"`
	AzureContainerRegistry struct {
		ArmResourceId string `json:"armResourceId"`
		LoginServer   string `json:"loginServer"`
		Name          string `json:"name"`
	} `json:"azureContainerRegistry"`
	AzureSubscriptionId   string `json:"azureSubscriptionId"`
	KeyvaultArmResourceId string `json:"keyvaultArmResourceId"`
}

// AgentConfigNameParameter defines model for AgentConfigNameParameter.
type AgentConfigNameParameter = externalRef0.AgentConfigName

// CertificateIdPathParameter defines model for CertificateIdPathParameter.
type CertificateIdPathParameter = externalRef0.Identifier

// CertificateTemplateIdentifierParameter defines model for CertificateTemplateIdentifierParameter.
type CertificateTemplateIdentifierParameter = externalRef0.Identifier

// IncludeCertificateParameter defines model for IncludeCertificateParameter.
type IncludeCertificateParameter = bool

// NamespaceIdParameter defines model for NamespaceIdParameter.
type NamespaceIdParameter = externalRef0.Identifier

// NamespaceKindParameter defines model for NamespaceKindParameter.
type NamespaceKindParameter = externalRef0.NamespaceKind

// AgentConfigurationResponse defines model for AgentConfigurationResponse.
type AgentConfigurationResponse = externalRef0.AgentConfiguration

// CertificateResponse defines model for CertificateResponse.
type CertificateResponse = externalRef0.CertificateInfo

// PatchServiceConfigJSONBody defines parameters for PatchServiceConfig.
type PatchServiceConfigJSONBody = interface{}

// PatchServiceConfigParamsConfigPath defines parameters for PatchServiceConfig.
type PatchServiceConfigParamsConfigPath string

// GetAgentConfigurationParams defines parameters for GetAgentConfiguration.
type GetAgentConfigurationParams struct {
	RefreshToken               *string `form:"refreshToken,omitempty" json:"refreshToken,omitempty"`
	XSmallkmsIfVersionNotMatch *string `json:"X-Smallkms-If-Version-Not-Match,omitempty"`
}

// IssueCertificateFromTemplateParams defines parameters for IssueCertificateFromTemplate.
type IssueCertificateFromTemplateParams struct {
	IncludeCertificate *IncludeCertificateParameter            `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
	Force              *bool                                   `form:"force,omitempty" json:"force,omitempty"`
	Enroll             *bool                                   `form:"enroll,omitempty" json:"enroll,omitempty"`
	Tags               *[]externalRef0.TemplatedCertificateTag `form:"tags,omitempty" json:"tags,omitempty"`
}

// AddKeyVaultRoleAssignmentParams defines parameters for AddKeyVaultRoleAssignment.
type AddKeyVaultRoleAssignmentParams struct {
	RoleDefinitionId string `form:"roleDefinitionId" json:"roleDefinitionId"`
}

// GetCertificateParams defines parameters for GetCertificate.
type GetCertificateParams struct {
	IncludeCertificate *IncludeCertificateParameter `form:"includeCertificate,omitempty" json:"includeCertificate,omitempty"`
	TemplateId         *externalRef0.Identifier     `form:"templateId,omitempty" json:"templateId,omitempty"`
}

// CreateProfileJSONRequestBody defines body for CreateProfile for application/json ContentType.
type CreateProfileJSONRequestBody = CreateProfileRequest

// PatchServiceConfigJSONRequestBody defines body for PatchServiceConfig for application/json ContentType.
type PatchServiceConfigJSONRequestBody = PatchServiceConfigJSONBody

// AgentCallbackJSONRequestBody defines body for AgentCallback for application/json ContentType.
type AgentCallbackJSONRequestBody = externalRef0.AgentCallbackRequest

// PutAgentConfigurationJSONRequestBody defines body for PutAgentConfiguration for application/json ContentType.
type PutAgentConfigurationJSONRequestBody = externalRef0.AgentConfigurationParameters

// PutCertificateTemplateJSONRequestBody defines body for PutCertificateTemplate for application/json ContentType.
type PutCertificateTemplateJSONRequestBody = CertificateTemplateParameters

// CreateLinkedCertificateTemplateJSONRequestBody defines body for CreateLinkedCertificateTemplate for application/json ContentType.
type CreateLinkedCertificateTemplateJSONRequestBody = CreateLinkedCertificateTemplateParameters

// AsCreateManagedApplicationProfileRequest returns the union data inside the CreateProfileRequest as a CreateManagedApplicationProfileRequest
func (t CreateProfileRequest) AsCreateManagedApplicationProfileRequest() (CreateManagedApplicationProfileRequest, error) {
	var body CreateManagedApplicationProfileRequest
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromCreateManagedApplicationProfileRequest overwrites any union data inside the CreateProfileRequest as the provided CreateManagedApplicationProfileRequest
func (t *CreateProfileRequest) FromCreateManagedApplicationProfileRequest(v CreateManagedApplicationProfileRequest) error {
	v.Type = "managed-application"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeCreateManagedApplicationProfileRequest performs a merge with any union data inside the CreateProfileRequest, using the provided CreateManagedApplicationProfileRequest
func (t *CreateProfileRequest) MergeCreateManagedApplicationProfileRequest(v CreateManagedApplicationProfileRequest) error {
	v.Type = "managed-application"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t CreateProfileRequest) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"type"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t CreateProfileRequest) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "managed-application":
		return t.AsCreateManagedApplicationProfileRequest()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t CreateProfileRequest) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *CreateProfileRequest) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}
