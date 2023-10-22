package admin

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	ctx "github.com/stephenzsy/small-kms/backend/internal/context"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
	"github.com/stephenzsy/small-kms/backend/shared"
)

type RequestContext = ctx.RequestContext
type ServiceConfig = models.ServiceConfigComposed
type PatchServiceConfigParamsConfigPath = models.PatchServiceConfigParamsConfigPath

type ServiceConfigDoc struct {
	kmsdoc.BaseDoc
	AppRoleIds struct {
		AgentActiveHost uuid.UUID `json:"Agent.ActiveHost"`
		AppAdmin        uuid.UUID `json:"App.Admin"`
	} `json:"appRoleIds"`
	AzureContainerRegistry struct {
		ArmResourceId string `json:"armResourceId"`
		LoginServer   string `json:"loginServer"`
		Name          string `json:"name"`
	} `json:"azureContainerRegistry"`
	AzureSubscriptionId   string `json:"azureSubscriptionId"`
	KeyvaultArmResourceId string `json:"keyvaultArmResourceId"`
}

func (d *ServiceConfigDoc) toModel() *ServiceConfig {
	if d == nil {
		return nil
	}
	m := ServiceConfig{}
	d.BaseDoc.PopulateResourceRef(&m.ResourceRef)
	m.AppRoleIds.AgentActiveHost = d.AppRoleIds.AgentActiveHost
	m.AppRoleIds.AppAdmin = d.AppRoleIds.AppAdmin
	m.AzureContainerRegistry.ArmResourceId = d.AzureContainerRegistry.ArmResourceId
	m.AzureContainerRegistry.LoginServer = d.AzureContainerRegistry.LoginServer
	m.AzureContainerRegistry.Name = d.AzureContainerRegistry.Name
	m.AzureSubscriptionId = d.AzureSubscriptionId
	m.KeyvaultArmResourceId = d.KeyvaultArmResourceId
	return &m
}

var serviceConfigDocLocator = shared.NewResourceLocator(
	shared.NewNamespaceIdentifier(shared.NamespaceKindProfile, shared.StringIdentifier(ns.ProfileNamespaceIDNameBuiltin)),
	shared.NewResourceIdentifier(shared.ResourceKindReserved, shared.StringIdentifier("service-config")))

func GetServiceConfig(c RequestContext) (*ServiceConfig, error) {
	doc := ServiceConfigDoc{}
	err := kmsdoc.Read(c, serviceConfigDocLocator, &doc)
	if err != nil {
		if err, isNotFound := common.IsAzCosmosNotFound(err); isNotFound {
			// store a new doc there
			doc.ID = serviceConfigDocLocator.GetID()
			doc.NamespaceID = serviceConfigDocLocator.GetNamespaceID()
			if err := kmsdoc.Create(c, &doc); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return doc.toModel(), nil
}

func PatchServiceConfig(c RequestContext, configPath PatchServiceConfigParamsConfigPath) (*ServiceConfig, error) {
	var err error
	doc := ServiceConfigDoc{}
	var patchData any
	switch configPath {
	case models.ServiceConfigPathKeyvaultArmResourceId:
		err = c.Bind(&doc.KeyvaultArmResourceId)
		patchData = &doc.KeyvaultArmResourceId
	case models.ServiceConfigPathAzureSubscriptionId:
		err = c.Bind(&doc.AzureSubscriptionId)
		patchData = &doc.AzureSubscriptionId
	case models.ServiceConfigPathAzureContainerRegistry:
		err = c.Bind(&doc.AzureContainerRegistry)
		patchData = &doc.AzureContainerRegistry
	case models.ServiceConfigPathAppRoleIds:
		err = c.Bind(&doc.AppRoleIds)
		patchData = &doc.AppRoleIds
	default:
		return nil, fmt.Errorf("%w:invalid config path", common.ErrStatusBadRequest)
	}
	if err != nil {
		return nil, fmt.Errorf("%w:invalid input body", common.ErrStatusBadRequest)
	}
	ops := azcosmos.PatchOperations{}
	ops.AppendSet("/"+string(configPath), patchData)
	err = kmsdoc.PatchWithWriteBack(c, serviceConfigDocLocator, &doc, ops)
	return doc.toModel(), err
}
