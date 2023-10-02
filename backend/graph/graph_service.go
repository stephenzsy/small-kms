package graph

import (
	ctx "context"
	"fmt"

	"github.com/google/uuid"
	msgraphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/stephenzsy/small-kms/backend/common"
	"github.com/stephenzsy/small-kms/backend/kmsdoc"
	"github.com/stephenzsy/small-kms/backend/utils"
)

type GraphProfileDocument interface {
	kmsdoc.KmsDocument
	GetDisplayName() string
	init(tenantID uuid.UUID, obj GraphProfileable, odataType MsGraphOdataType)
	IsValid() bool
	GetOdataType() MsGraphOdataType
}

type GraphService interface {
	common.CommonConfig
	GetGraphObjectByID(c ctx.Context, objectID uuid.UUID) (msgraphmodels.DirectoryObjectable, error)

	GetDeviceByDeviceID(c ctx.Context, deviceID uuid.UUID) (msgraphmodels.Deviceable, error)

	GetGraphProfileDoc(c ctx.Context, objectID uuid.UUID, odataType MsGraphOdataType) (GraphProfileDocument, error)

	DeleteGraphProfileDoc(c ctx.Context, doc GraphProfileDocument) error

	ListGraphProfilesByType(c ctx.Context, odataType MsGraphOdataType) ([]GraphProfileDocument, error)

	NewGraphProfileDoc(tenantID uuid.UUID, obj GraphProfileable) GraphProfileDocument

	NewGraphProfileDocWithType(tenantID uuid.UUID, obj GraphProfileable, odataType MsGraphOdataType) GraphProfileDocument
}

type graphService struct {
	common.CommonConfig
}

func NewGraphService(config common.CommonConfig) GraphService {
	s := graphService{
		CommonConfig: config,
	}

	return &s
}

func (s *graphService) GetGraphProfileDoc(c ctx.Context, objectID uuid.UUID, odataType MsGraphOdataType) (GraphProfileDocument, error) {
	docID := kmsdoc.NewKmsDocID(kmsdoc.DocTypeMsGraphObject, objectID)
	doc := newDocByODataType(odataType)
	if err := kmsdoc.AzCosmosRead(c, s.AzCosmosContainerClient(), s.TenantID(), docID, doc); err != nil {
		return nil, common.WrapAzRsNotFoundErr(err, fmt.Sprintf("graphobjectprofile/%s", objectID.String()))
	}
	return doc, nil
}

func (s *graphService) GetGraphObjectByID(c ctx.Context, objectID uuid.UUID) (msgraphmodels.DirectoryObjectable, error) {
	r, err := s.MsGraphClient().DirectoryObjects().ByDirectoryObjectId(objectID.String()).Get(c, nil)
	return r, common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("graphobject/%s", objectID.String()))
}

func (s *graphService) GetDeviceByDeviceID(c ctx.Context, deviceID uuid.UUID) (device msgraphmodels.Deviceable, err error) {
	r, err := s.MsGraphClient().DevicesWithDeviceId(utils.ToPtr(deviceID.String())).Get(c, nil)
	return r, common.WrapMsGraphNotFoundErr(err, fmt.Sprintf("deviceswithdeviceid/%s", deviceID.String()))
}

func (s *graphService) DeleteGraphProfileDoc(c ctx.Context, doc GraphProfileDocument) error {
	return kmsdoc.AzCosmosDelete(c, s.AzCosmosContainerClient(), doc)
}

func (s *graphService) ListGraphProfilesByType(c ctx.Context, odataType MsGraphOdataType) ([]GraphProfileDocument, error) {
	pager := s.queryProfilesByType(c, odataType)
	l, err := utils.PagerToList[GraphDoc](c, pager)
	if err != nil {
		return nil, err
	}
	result := make([]GraphProfileDocument, len(l))
	for i, item := range l {
		result[i] = &item
	}
	return result, nil
}

func newDocByODataType(odataType MsGraphOdataType) GraphProfileDocument {
	switch odataType {
	case MsGraphOdataTypeDevice:
		return new(DeviceDoc)
	case MsGraphOdataTypeApplication:
		return new(ApplicationDoc)
	case MsGraphOdataTypeServicePrincipal:
		return new(ServicePrincipalDoc)
	case MsGraphOdataTypeUser:
		return new(UserDoc)
	case MsGraphOdataTypeGroup:
		return new(GroupDoc)
	}
	return new(GraphDoc)
}

func (s *graphService) NewGraphProfileDoc(tenantID uuid.UUID, obj GraphProfileable) GraphProfileDocument {
	return s.NewGraphProfileDocWithType(tenantID, obj, MsGraphOdataType(utils.NilToDefault(obj.GetOdataType())))
}

func (s *graphService) NewGraphProfileDocWithType(tenantID uuid.UUID, obj GraphProfileable, odataType MsGraphOdataType) GraphProfileDocument {
	doc := newDocByODataType(MsGraphOdataType(utils.NilToDefault(obj.GetOdataType())))
	doc.init(tenantID, obj, odataType)
	return doc
}
