package order

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

func (s *ServiceSuite) TestPayOrderByIdSuccess() {
	var (
		uuid          = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		transactionID = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
		order         = model.Order{
			OrderUUID: uuid,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPendingPayment,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(order, nil).Once()
	s.paymentCli.On("PayOrder", s.ctx, uuid, userUUID, paymentMethod).Return(transactionID, nil)
	s.repo.On("UpdateOrder", s.ctx, model.UpdateOrderParam{OrderUUID: uuid, PaymentMethod: &paymentMethod, Status: lo.ToPtr(model.OrderStatusPaid), TransactionUUID: &transactionID}).Return(nil).Once()
	s.orderProducer.On("ProduceOrderPaid", s.ctx, model.OrderPaidEvent{EventUUID: mock.Anything, OrderUUID: uuid, UserUUID: userUUID, PaymentMethod: paymentMethod, TransactionUUID: transactionID}).Return(nil).Once()
	response, err := s.service.PayOrderByID(s.ctx, uuid, paymentMethod)
	s.NoError(err)
	s.Equal(response, transactionID)
}

func (s *ServiceSuite) TestPayOrderByIdUpdateError() {
	var (
		uuid          = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		transactionID = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
		order         = model.Order{
			OrderUUID: uuid,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPendingPayment,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(order, nil).Once()
	s.paymentCli.On("PayOrder", s.ctx, uuid, userUUID, paymentMethod).Return(transactionID, nil)
	s.repo.On("UpdateOrder", s.ctx, model.UpdateOrderParam{OrderUUID: uuid, PaymentMethod: &paymentMethod, Status: lo.ToPtr(model.OrderStatusPaid), TransactionUUID: &transactionID}).Return(model.ErrOrderNotFound).Once()

	response, err := s.service.PayOrderByID(s.ctx, uuid, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(response)
}

func (s *ServiceSuite) TestPayOrderByIdProduceError() {
	var (
		uuid          = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		transactionID = gofakeit.UUID()
		paymentMethod = model.PaymentMethodCard
		order         = model.Order{
			OrderUUID: uuid,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPendingPayment,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(order, nil).Once()
	s.paymentCli.On("PayOrder", s.ctx, uuid, userUUID, paymentMethod).Return(transactionID, nil)
	s.repo.On("UpdateOrder", s.ctx, model.UpdateOrderParam{OrderUUID: uuid, PaymentMethod: &paymentMethod, Status: lo.ToPtr(model.OrderStatusPaid), TransactionUUID: &transactionID}).Return(nil).Once()
	s.orderProducer.On("ProduceOrderPaid", s.ctx, model.OrderPaidEvent{EventUUID: mock.Anything, OrderUUID: uuid, UserUUID: userUUID, PaymentMethod: paymentMethod, TransactionUUID: transactionID}).Return(errors.New("produce error")).Once()
	response, err := s.service.PayOrderByID(s.ctx, uuid, paymentMethod)
	s.Error(err)
	s.Empty(response)
}

func (s *ServiceSuite) TestPayOrderByIdPayError() {
	var (
		uuid          = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		ErrPay        = errors.New("pay error")
		paymentMethod = model.PaymentMethodCard
		order         = model.Order{
			OrderUUID: uuid,
			UserUUID:  userUUID,
			Status:    model.OrderStatusPendingPayment,
		}
	)
	s.repo.On("GetOrder", s.ctx, uuid).Return(order, nil).Once()
	s.paymentCli.On("PayOrder", s.ctx, uuid, userUUID, paymentMethod).Return("", ErrPay)

	response, err := s.service.PayOrderByID(s.ctx, uuid, paymentMethod)
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

	response, err := s.service.PayOrderByID(s.ctx, uuid, paymentMethod)
	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(response)
}
