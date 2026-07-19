package order

import (
	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestGetOrderByIDSuccess() {
	var (
		uuid = gofakeit.UUID()
		part = model.Order{
			OrderUUID: uuid,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(part, nil).Once()

	response, err := s.service.GetOrderByID(s.ctx, uuid)
	s.NoError(err)
	s.Equal(response, part)
}
func (s *ServiceSuite) TestGetOrderByIDError() {
	var (
		uuid = gofakeit.UUID()
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(model.Order{}, model.ErrOrderNotFound).Once()

	response, err := s.service.GetOrderByID(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(response)
}