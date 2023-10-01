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

type GraphProfileDoc interface {
	kmsdoc.KmsDocument
}

type GraphService interface {
	common.CommonConfig
	GetGraphObjectByID(c ctx.Context, objectID uuid.UUID) (msgraphmodels.DirectoryObjectable, error)

	GetDeviceByDeviceID(c ctx.Context, deviceID uuid.UUID) (msgraphmodels.Deviceable, error)

	GetGraphProfileDoc(c ctx.Context, objectID uuid.UUID, docExtension kmsdoc.KmsDocTypeExtName) (GraphProfileDoc, error)
	//	StoreDevice(device msgraphmodels.Deviceable) error

	DeleteGraphProfileDoc(c ctx.Context, doc GraphProfileDoc) error

	NewDeviceDocFromGraph(msgraphmodels.Deviceable) *DeviceDoc
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

func (s *graphService) GetGraphProfileDoc(c ctx.Context, objectID uuid.UUID, docExtension kmsdoc.KmsDocTypeExtName) (GraphProfileDoc, error) {
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
	return s.MsGraphClient().DirectoryObjects().ByDirectoryObjectId(objectID.String()).Get(c, nil)
}

func (s *graphService) GetDeviceByDeviceID(c ctx.Context, deviceID uuid.UUID) (device msgraphmodels.Deviceable, err error) {
	return s.MsGraphClient().DevicesWithDeviceId(utils.ToPtr(deviceID.String())).Get(c, nil)
}

func (s *graphService) DeleteGraphProfileDoc(c ctx.Context, doc GraphProfileDoc) error {
	return kmsdoc.AzCosmosDelete(c, s.AzCosmosContainerClient(), doc)
}
