package v1

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/H1dEx/ms-rocket/inventory/internal/converter"
	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	inventory_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

func (a *ApiSuite) TestListPartsSuccess() {
	partUUID := gofakeit.UUID()
	parts := []model.Part{{Uuid: partUUID}}
	filter := &inventory_v1.ListPartsRequest{}

	a.service.On("GetList", a.ctx, &model.PartFilter{}).Return(parts, nil)

	res, err := a.api.ListParts(a.ctx, filter)
	a.NoError(err)
	expect := &inventory_v1.ListPartsResponse{Parts: []*inventory_v1.Part{converter.PartToProto(parts[0])}}
	a.Equal(res, expect)
}

func (a *ApiSuite) TestListPartsError() {
	filter := &inventory_v1.ListPartsRequest{}

	ErrTest := errors.New("test error")

	a.service.On("GetList", a.ctx, &model.PartFilter{}).Return(nil, ErrTest)

	res, err := a.api.ListParts(a.ctx, filter)
	a.Error(err)
	a.ErrorIs(err, ErrTest)
	a.Empty(res)
}
