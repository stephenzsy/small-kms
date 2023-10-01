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
}

type GraphService interface {
	common.CommonConfig
	GetGraphObjectByID(c ctx.Context, objectID uuid.UUID) (msgraphmodels.DirectoryObjectable, error)

	GetDeviceByDeviceID(c ctx.Context, deviceID uuid.UUID) (msgraphmodels.Deviceable, error)

	GetGraphProfileDoc(c ctx.Context, objectID uuid.UUID, docExtension kmsdoc.KmsDocTypeExtName) (GraphProfileDocument, error)
	//	StoreDevice(device msgraphmodels.Deviceable) error

	DeleteGraphProfileDoc(c ctx.Context, doc GraphProfileDocument) error

	NewDeviceDocFromGraph(msgraphmodels.Deviceable) *DeviceDoc

	ListGraphProfilesByType(c ctx.Context, docExtension kmsdoc.KmsDocTypeExtName) ([]GraphProfileDocument, error)
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

func (s *graphService) GetGraphProfileDoc(c ctx.Context, objectID uuid.UUID, docExtension kmsdoc.KmsDocTypeExtName) (GraphProfileDocument, error) {
	docID := kmsdoc.NewKmsDocIDExt(kmsdoc.DocTypeMsGraphObject, objectID, docExtension)
	switch docExtension {
	case kmsdoc.DocTypeExtNameDevice:
		doc := DeviceDoc{}
		if err := kmsdoc.AzCosmosRead(c, s.AzCosmosContainerClient(), s.TenantID(), docID, &doc); err != nil {
			return nil, common.WrapAzRsNotFoundErr(err, fmt.Sprintf("graphobjectprofile/%s", objectID.String()))
		}
		return &doc, nil
	}
	return nil, fmt.Errorf("%w: unsupported graph profile type %s", common.ErrStatusBadRequest, docExtension)
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

func (s *graphService) ListGraphProfilesByType(c ctx.Context, docExtension kmsdoc.KmsDocTypeExtName) ([]GraphProfileDocument, error) {
	pager := s.queryProfilesByType(c, docExtension)
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
