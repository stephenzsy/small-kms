package admin

import (
	"github.com/google/uuid"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/internal/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/models"
	ns "github.com/stephenzsy/small-kms/backend/namespace"
)

type RequestContext = common.RequestContext
type ServiceConfig = models.ServiceConfigComposed

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

var serviceConfigDocLocator = models.NewResourceLocator(
	models.NewNamespaceID(models.NamespaceKindProfile, common.StringIdentifier(ns.ProfileNamespaceIDNameBuiltin)),
	models.NewResourceID(models.ResourceKindReserved, common.StringIdentifier("service-config")))

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
