package order

import (
	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
)

func (s *ServiceSuite) TestOrderCancelByIdSuccess() {
	var (
		orderUUID = gofakeit.UUID()
		order     = model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPendingPayment,
		}
		param = model.UpdateOrderParam{OrderUUID: orderUUID, Status: lo.ToPtr(model.OrderStatusCancelled)}
	)
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.repo.On("UpdateOrder", s.ctx, param).Return(nil).Once()

	err := s.service.OrderCancelById(s.ctx, orderUUID)
	s.NoError(err)
}
func (s *ServiceSuite) TestOrderCancelByIdUpdateErr() {
	var (
		orderUUID = gofakeit.UUID()
		order     = model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPendingPayment,
		}
		param = model.UpdateOrderParam{OrderUUID: orderUUID, Status: lo.ToPtr(model.OrderStatusCancelled)}
	)
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()
	s.repo.On("UpdateOrder", s.ctx, param).Return(model.ErrOrderNotFound).Once()

	err := s.service.OrderCancelById(s.ctx, orderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
func (s *ServiceSuite) TestOrderCancelByIdConflictErr() {
	var (
		orderUUID = gofakeit.UUID()
		order     = model.Order{
			OrderUUID: orderUUID,
			Status:    model.OrderStatusPaid,
		}
	)
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(order, nil).Once()

	err := s.service.OrderCancelById(s.ctx, orderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrNotPendingStatus)
}

func (s *ServiceSuite) TestOrderCancelByIdNotFoundErr() {
	var (
		orderUUID = gofakeit.UUID()
	)
	s.repo.On("GetOrder", s.ctx, orderUUID).Return(model.Order{}, model.ErrOrderNotFound).Once()

	err := s.service.OrderCancelById(s.ctx, orderUUID)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
}
