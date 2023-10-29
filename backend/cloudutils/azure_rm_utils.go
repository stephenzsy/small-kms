package cloudutils

import (
	"fmt"

	"github.com/google/uuid"
)

type AzureRMIDBuilder interface {
	Build() string
}

// subscription id builder
type AzureRMSubscriptionIDBuilder interface {
	AzureRMIDBuilder
	WithResourceGroup(resourceGroup string) AzureRMResourceGroupIDBuilder
	WithRoleDefinitionID(roleDefinitionID uuid.UUID) AzureRMSubscriptionIDBuilder
}

type AzureSubscriptionResourceIDBuilder struct {
	SubscriptionID string
}

func (b *AzureSubscriptionResourceIDBuilder) Build() string {
	return fmt.Sprintf("/subscriptions/%s", b.SubscriptionID)
}

func (b *AzureSubscriptionResourceIDBuilder) WithRoleDefinitionID(roleDefinitionID uuid.UUID) AzureRMSubscriptionIDBuilder {
	return &AzureRoleDefinitionResourceIDBuilder{
		AzureSubscriptionResourceIDBuilder: *b,
		RoleDefinitionID:                   roleDefinitionID,
	}
}

func (b *AzureSubscriptionResourceIDBuilder) WithResourceGroup(resourceGroup string) AzureRMResourceGroupIDBuilder {
	return &AzureResourceGroupResourceIDBuilder{
		AzureSubscriptionResourceIDBuilder: *b,
		ResourceGroup:                      resourceGroup,
	}
}

var _ AzureRMSubscriptionIDBuilder = (*AzureSubscriptionResourceIDBuilder)(nil)

// role assignment id builder
type AzureRoleDefinitionResourceIDBuilder struct {
	AzureSubscriptionResourceIDBuilder
	RoleDefinitionID uuid.UUID
}

func (b *AzureRoleDefinitionResourceIDBuilder) Build() string {
	return fmt.Sprintf("%s/providers/Microsoft.Authorization/roleDefinitions/%s", b.AzureSubscriptionResourceIDBuilder.Build(), b.RoleDefinitionID.String())
}

var _ AzureRMSubscriptionIDBuilder = (*AzureRoleDefinitionResourceIDBuilder)(nil)

// resource group id builder

type AzureRMResourceGroupIDBuilder interface {
	AzureRMIDBuilder
	WithKeyVault(vaultName, category, itemName string) AzureRMIDBuilder
	WithContainerRegistry(registryName string) AzureRMIDBuilder
}

type AzureResourceGroupResourceIDBuilder struct {
	AzureSubscriptionResourceIDBuilder
	ResourceGroup string
}

// WithContainerRegistry implements AzureRMResourceGroupIDBuilder.
func (b *AzureResourceGroupResourceIDBuilder) WithContainerRegistry(registryName string) AzureRMIDBuilder {
	return &AzureContainerRegistryResourceIDBuilder{
		AzureResourceGroupResourceIDBuilder: *b,
		RegistryName:                        registryName,
	}
}

func (b *AzureResourceGroupResourceIDBuilder) WithKeyVault(vaultName, catetory, itemName string) AzureRMIDBuilder {
	return &AzureKeyVaultResourceIDBuilder{
		AzureResourceGroupResourceIDBuilder: *b,
		VaultName:                           vaultName,
		Category:                            catetory,
		ItemName:                            itemName,
	}
}

func (b *AzureResourceGroupResourceIDBuilder) Build() string {
	return fmt.Sprintf("%s/resourceGroups/%s", b.AzureSubscriptionResourceIDBuilder.Build(), b.ResourceGroup)
}

var _ AzureRMResourceGroupIDBuilder = (*AzureResourceGroupResourceIDBuilder)(nil)

// key vault id builder
type AzureKeyVaultResourceIDBuilder struct {
	AzureResourceGroupResourceIDBuilder
	VaultName string
	Category  string
	ItemName  string
}

func (b *AzureKeyVaultResourceIDBuilder) Build() string {
	return fmt.Sprintf("%s/providers/Microsoft.KeyVault/vaults/%s/%s/%s", b.AzureResourceGroupResourceIDBuilder.Build(), b.VaultName, b.Category, b.ItemName)
}

type AzureContainerRegistryResourceIDBuilder struct {
	AzureResourceGroupResourceIDBuilder
	RegistryName string
}

func (b *AzureContainerRegistryResourceIDBuilder) Build() string {
	return fmt.Sprintf("%s/providers/Microsoft.ContainerRegistry/registries/%s", b.AzureResourceGroupResourceIDBuilder.Build(), b.RegistryName)
}

var _ AzureRMIDBuilder = (*AzureContainerRegistryResourceIDBuilder)(nil)
