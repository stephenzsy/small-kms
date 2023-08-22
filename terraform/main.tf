# Configure the Azure provider
terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.70.0"
    }
  }
}

provider "azurerm" {
  skip_provider_registration = true
  features {
  }
}

variable "resource_group_name" {
  type = string
}

variable "cosmosdb_account_name" {
  type = string
}

variable "principal_id" {
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

resource "azurerm_cosmosdb_sql_role_assignment" "backend" {
  name                = "9c2b265e-6012-4644-af81-a8420580dd2e"
  resource_group_name = data.azurerm_resource_group.default.name
  account_name        = data.azurerm_cosmosdb_account.default.name
  role_definition_id  = data.azurerm_cosmosdb_sql_role_definition.contributor.id
  principal_id        = var.principal_id
  scope               = data.azurerm_cosmosdb_account.default.id
}
