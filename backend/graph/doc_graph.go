package graph

import (
	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type MsGraphOdataType string

const (
	MsGraphOdataTypeDevice           MsGraphOdataType = "#microsoft.graph.device"
	MsGraphOdataTypeUser             MsGraphOdataType = "#microsoft.graph.user"
	MsGraphOdataTypeGroup            MsGraphOdataType = "#microsoft.graph.group"
	MsGraphOdataTypeApplication      MsGraphOdataType = "#microsoft.graph.application"
	MsGraphOdataTypeServicePrincipal MsGraphOdataType = "#microsoft.graph.servicePrincipal"
)

// this kind of docs represent a graph object browsable in the app, it should not be used to persist data, always query live graph service
type GraphDoc struct {
	kmsdoc.BaseDoc

	GraphType   MsGraphOdataType `json:"@odata.type"`
	DisplayName string           `json:"displayName"`

	service *graphService
}

type GraphProfileable interface {
	msgraphmodels.DirectoryObjectable
	GetDisplayName() *string
}

func (s *graphService) init(doc *GraphDoc, graphObj GraphProfileable, graphType MsGraphOdataType) {
	oid, _ := uuid.Parse(utils.NilToDefault(graphObj.GetId()))
	var docTypeExtName kmsdoc.KmsDocTypeExtName
	switch doc.GraphType {
	case MsGraphOdataTypeDevice:
		docTypeExtName = kmsdoc.DocTypeExtNameDevice
	case MsGraphOdataTypeUser:
		docTypeExtName = kmsdoc.DocTypeExtNameUser
	case MsGraphOdataTypeGroup:
		docTypeExtName = kmsdoc.DocTypeExtNameGroup
	case MsGraphOdataTypeApplication:
		docTypeExtName = kmsdoc.DocTypeExtNameApplication
	case MsGraphOdataTypeServicePrincipal:
		docTypeExtName = kmsdoc.DocTypeExtNameServicePrincipal
	}
	doc.NamespaceID = s.TenantID()
	doc.BaseDoc.ID = kmsdoc.NewKmsDocIDExt(kmsdoc.DocTypeMsGraphObject, oid, docTypeExtName)
	doc.DisplayName = utils.NilToDefault(graphObj.GetDisplayName())
}
