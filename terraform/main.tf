# Configure the Azure provider
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.75.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.5.1"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.43.0"
    }
  }
}

provider "azurerm" {
  skip_provider_registration = true
  features {
  }
}

data "azurerm_client_config" "current" {}

variable "resource_group_name" {
  type = string
}

variable "cosmosdb_account_name" {
  type = string
}

variable "principal_id" {
  type = string
}

variable "gha_subject_identifier" {
  type = string
}

data "azurerm_resource_group" "default" {
  name = var.resource_group_name
}

data "azurerm_cosmosdb_account" "default" {
  name                = var.cosmosdb_account_name
  resource_group_name = data.azurerm_resource_group.default.name
}

data "azurerm_cosmosdb_sql_role_definition" "contributor" {
  resource_group_name = data.azurerm_resource_group.default.name
  account_name        = data.azurerm_cosmosdb_account.default.name
  role_definition_id  = "00000000-0000-0000-0000-000000000002"
}

resource "random_uuid" "backendSqlRoleAssignmentName" {}

resource "azurerm_cosmosdb_sql_role_assignment" "backend" {
  name                = random_uuid.backendSqlRoleAssignmentName.result
  resource_group_name = data.azurerm_resource_group.default.name
  account_name        = data.azurerm_cosmosdb_account.default.name
  role_definition_id  = data.azurerm_cosmosdb_sql_role_definition.contributor.id
  principal_id        = var.principal_id
  scope               = data.azurerm_cosmosdb_account.default.id
}

resource "random_pet" "default" {}

resource "azurerm_cosmosdb_sql_database" "db" {
  name                = "smallkms-${random_pet.default.id}"
  resource_group_name = data.azurerm_cosmosdb_account.default.resource_group_name
  account_name        = data.azurerm_cosmosdb_account.default.name
}

resource "azurerm_cosmosdb_sql_container" "kmsdbContainer" {
  name                  = "Certs"
  resource_group_name   = data.azurerm_cosmosdb_account.default.resource_group_name
  account_name          = data.azurerm_cosmosdb_account.default.name
  database_name         = azurerm_cosmosdb_sql_database.db.name
  partition_key_path    = "/namespaceId"
  partition_key_version = 1

  autoscale_settings {
    max_throughput = 1000
  }
}

resource "azurerm_user_assigned_identity" "backendManagedIdentity" {
  location            = data.azurerm_resource_group.default.location
  name                = "smallkms-backend-${random_pet.default.id}"
  resource_group_name = data.azurerm_resource_group.default.name
}

resource "random_uuid" "backendIdentitySqlRoleAssignmentName" {}
resource "azurerm_cosmosdb_sql_role_assignment" "backendManagedIdentityDBAccess" {
  name                = random_uuid.backendIdentitySqlRoleAssignmentName.result
  resource_group_name = data.azurerm_resource_group.default.name
  account_name        = data.azurerm_cosmosdb_account.default.name
  role_definition_id  = data.azurerm_cosmosdb_sql_role_definition.contributor.id
  principal_id        = azurerm_user_assigned_identity.backendManagedIdentity.principal_id
  scope               = data.azurerm_cosmosdb_account.default.id
}

resource "azurerm_key_vault" "default" {
  name                       = "smallkms-${random_pet.default.id}"
  location                   = data.azurerm_resource_group.default.location
  resource_group_name        = data.azurerm_resource_group.default.name
  tenant_id                  = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days = 7
  purge_protection_enabled   = false
  enable_rbac_authorization  = true
  sku_name                   = "standard"
}

resource "azurerm_log_analytics_workspace" "default" {
  name                = "smallkms-log-${random_pet.default.id}"
  location            = data.azurerm_resource_group.default.location
  resource_group_name = data.azurerm_resource_group.default.name
  retention_in_days   = 30
}

resource "azurerm_container_app_environment" "backend" {
  name                       = "smallkms-backend-env-${random_pet.default.id}"
  location                   = data.azurerm_resource_group.default.location
  resource_group_name        = data.azurerm_resource_group.default.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.default.id
}

resource "azurerm_container_app" "backend" {
  name                         = "smallkms-${random_pet.default.id}"
  container_app_environment_id = azurerm_container_app_environment.backend.id
  resource_group_name          = data.azurerm_resource_group.default.name
  revision_mode                = "Single"

  template {
    container {
      name   = "examplecontainerapp"
      image  = "mcr.microsoft.com/azuredocs/containerapps-helloworld:latest"
      cpu    = 0.25
      memory = "0.5Gi"
    }
  }
}

resource "azurerm_storage_account" "default" {
  name                     = join("", ["smallkms", replace(random_pet.default.id, "-", "")])
  location                 = data.azurerm_resource_group.default.location
  resource_group_name      = data.azurerm_resource_group.default.name
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "certs" {
  name                  = "certs"
  storage_account_name  = azurerm_storage_account.default.name
  container_access_type = "private"
}

resource "azurerm_container_registry" "acr" {
  name                = join("", ["smallkmscr", replace(random_pet.default.id, "-", "")])
  location            = data.azurerm_resource_group.default.location
  resource_group_name = data.azurerm_resource_group.default.name
  sku                 = "Basic"
  admin_enabled       = false
}

resource "azurerm_user_assigned_identity" "deployment" {
  location            = data.azurerm_resource_group.default.location
  name                = "smallkms-deployment-${random_pet.default.id}"
  resource_group_name = data.azurerm_resource_group.default.name
}

resource "azurerm_role_assignment" "example" {
  scope                = azurerm_container_registry.acr.id
  role_definition_name = "AcrPush"
  principal_id         = azurerm_user_assigned_identity.deployment.principal_id
}

resource "azurerm_federated_identity_credential" "deploymentGHA" {
  name                = "smallkms-deployment-gha-${random_pet.default.id}"
  resource_group_name = data.azurerm_resource_group.default.name
  audience            = ["api://AzureADTokenExchange"]
  issuer              = "https://token.actions.githubusercontent.com"
  parent_id           = azurerm_user_assigned_identity.deployment.id
  subject             = var.gha_subject_identifier
}
