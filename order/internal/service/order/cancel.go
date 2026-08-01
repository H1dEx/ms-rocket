package order

import (
	"context"

	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) OrderCancelByID(ctx context.Context, orderUUID string) error {
	order, err := s.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		logger.Debug(ctx, "failed to get order", zap.Error(err))
		return err
	}

	if order.Status != model.OrderStatusPendingPayment {
		return model.ErrNotPendingStatus
	}

	param := model.UpdateOrderParam{OrderUUID: orderUUID, Status: lo.ToPtr(model.OrderStatusCancelled)}
	err = s.repo.UpdateOrder(ctx, param)

	return err
}
