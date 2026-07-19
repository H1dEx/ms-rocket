package order

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderByIdSuccess() {
	var (
		uuid          = gofakeit.UUID()
		userUuid      = gofakeit.UUID()
		transactionID = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
		order         = model.Order{
			OrderUUID: uuid,
			UserUUID:  userUuid,
			Status:    model.OrderStatusPendingPayment,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(order, nil).Once()
	s.paymentCli.On("PayOrder", s.ctx, uuid, userUuid, paymentMethod).Return(transactionID, nil)
	s.repo.On("UpdateOrder", s.ctx, model.UpdateOrderParam{OrderUUID: uuid, PaymentMethod: &paymentMethod, Status: lo.ToPtr(model.OrderStatusPaid), TransactionUUID: &transactionID}).Return(nil).Once()

	response, err := s.service.PayOrderById(s.ctx, uuid, paymentMethod)
	s.NoError(err)
	s.Equal(response, transactionID)
}

func (s *ServiceSuite) TestPayOrderByIdUpdateError() {
	var (
		uuid          = gofakeit.UUID()
		userUuid      = gofakeit.UUID()
		transactionID = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
		order         = model.Order{
			OrderUUID: uuid,
			UserUUID:  userUuid,
			Status:    model.OrderStatusPendingPayment,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(order, nil).Once()
	s.paymentCli.On("PayOrder", s.ctx, uuid, userUuid, paymentMethod).Return(transactionID, nil)
	s.repo.On("UpdateOrder", s.ctx, model.UpdateOrderParam{OrderUUID: uuid, PaymentMethod: &paymentMethod, Status: lo.ToPtr(model.OrderStatusPaid), TransactionUUID: &transactionID}).Return(model.ErrOrderNotFound).Once()

	response, err := s.service.PayOrderById(s.ctx, uuid, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(response)
}

func (s *ServiceSuite) TestPayOrderByIdPayError() {
	var (
		uuid          = gofakeit.UUID()
		userUuid      = gofakeit.UUID()
		ErrPay        = errors.New("pay error")
		paymentMethod = model.PaymentMethodCard
		order         = model.Order{
			OrderUUID: uuid,
			UserUUID:  userUuid,
			Status:    model.OrderStatusPendingPayment,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(order, nil).Once()
	s.paymentCli.On("PayOrder", s.ctx, uuid, userUuid, paymentMethod).Return("", ErrPay)

	response, err := s.service.PayOrderById(s.ctx, uuid, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, ErrPay)
	s.Empty(response)
}

func (s *ServiceSuite) TestPayOrderByIdGetError() {
	var (
		uuid          = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(model.Order{}, model.ErrOrderNotFound).Once()

	response, err := s.service.PayOrderById(s.ctx, uuid, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(response)
}
