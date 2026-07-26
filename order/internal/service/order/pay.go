package order

import (
	"context"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) PayOrderByID(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (transactionUUID string, err error) {
	order, err := s.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return "", err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return "", model.ErrNotPendingStatus
	}

	transactionID, err := s.paymentClient.PayOrder(ctx, orderUUID, order.UserUUID, paymentMethod)
	if err != nil {
		logger.Error(ctx, "failed to pay order", zap.Error(err))
		return "", err
	}

	param := model.UpdateOrderParam{OrderUUID: orderUUID, PaymentMethod: &paymentMethod, Status: lo.ToPtr(model.OrderStatusPaid), TransactionUUID: &transactionID}

	err = s.repo.UpdateOrder(ctx, param)
	if err != nil {
		logger.Error(ctx, "failed to update order", zap.Error(err))
		return "", err
	}
	return transactionID, nil
}
