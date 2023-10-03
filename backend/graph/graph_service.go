package graph

import (
	ctx "context"
	"fmt"

	"github.com/google/uuid"
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

	GetGraphProfileDoc(c ctx.Context, objectID uuid.UUID, odataType MsGraphOdataType) (GraphProfileDocument, error)

	DeleteGraphProfileDoc(c ctx.Context, doc GraphProfileDocument) error

	ListGraphProfilesByType(c ctx.Context, odataType MsGraphOdataType) ([]GraphDoc, error)

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
	if doc.GetDeleted() != nil && !doc.GetDeleted().IsZero() {
		return nil, fmt.Errorf("%w:graphobjectprofile/%s:deleted", common.ErrStatusNotFound, objectID)
	}
	return doc, nil
}

func (s *graphService) DeleteGraphProfileDoc(c ctx.Context, doc GraphProfileDocument) error {
	return kmsdoc.AzCosmosDelete(c, s.AzCosmosContainerClient(), doc)
}

func (s *graphService) ListGraphProfilesByType(c ctx.Context, odataType MsGraphOdataType) ([]GraphDoc, error) {
	pager := s.queryProfilesByType(c, odataType)
	return utils.PagerToList[GraphDoc](c, pager)
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
