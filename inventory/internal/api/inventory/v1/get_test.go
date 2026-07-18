package v1

import (
	"errors"

	"github.com/H1dEx/ms-rocket/inventory/internal/converter"
	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *ApiSuite) TestGetPartSuccess() {
	uuid := gofakeit.UUID()
	part := model.Part{Uuid: uuid}
	a.service.On("GetPart", a.ctx, uuid).Return(part, nil)

	res, err := a.api.GetPart(a.ctx, &inventoryV1.GetPartRequest{Uuid: uuid})
	a.NoError(err)
	expect := &inventoryV1.GetPartResponse{Part: converter.PartToProto(part)}
	a.Equal(expect, res)
}

func (a *ApiSuite) TestGetPartNotFound() {
	uuid := gofakeit.UUID()

	expectedErr := status.Errorf(codes.NotFound, "part with UUID %s not found", uuid)
	a.service.On("GetPart", a.ctx, uuid).Return(model.Part{}, model.ErrPartNotFound)

	res, err := a.api.GetPart(a.ctx, &inventoryV1.GetPartRequest{Uuid: uuid})
	a.Error(err)
	a.ErrorIs(err, expectedErr)
	a.Empty(res)
}

func (a *ApiSuite) TestGetPartUnknownError() {
	uuid := gofakeit.UUID()

	expectedErr := errors.New("Unknown err")
	a.service.On("GetPart", a.ctx, uuid).Return(model.Part{}, expectedErr)

	res, err := a.api.GetPart(a.ctx, &inventoryV1.GetPartRequest{Uuid: uuid})
	a.Error(err)
	a.ErrorIs(err, expectedErr)
	a.Empty(res)
}
