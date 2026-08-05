package order_consumer

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (s *service) HandleOrderAssembled(ctx context.Context, message kafka.Message) error {
	event, err := s.decoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "failed to decode order assembled event", zap.Error(err))
		return err
	}
	logger.Info(ctx, "order assembled event", zap.Any("event", event))

	err = s.repo.UpdateOrder(ctx, model.UpdateOrderParam{
		OrderUUID: event.OrderUUID,
		Status:    lo.ToPtr(model.OrderStatusShipped),
	})
	if err != nil {
		logger.Error(ctx, "failed to update order", zap.Error(err))
		return err
	}

	logger.Info(ctx, "order status updated to shipped", zap.String("order_uuid", event.OrderUUID))

	return nil
}
