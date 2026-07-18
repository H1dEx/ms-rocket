package order

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/samber/lo"
)

func (s *service) OrderCancelById(ctx context.Context, orderUUID string) error {
	order, err := s.repo.GetOrder(ctx, orderUUID)
	if err != nil {
		return nil
	}

	if order.Status != model.OrderStatusPendingPayment {
		return model.ErrNotPendingStatus
	}

	param := model.UpdateOrderParam{OrderUUID: orderUUID, Status: lo.ToPtr(model.OrderStatusCancelled)}
	err = s.repo.UpdateOrder(ctx, param)

	return err
}
