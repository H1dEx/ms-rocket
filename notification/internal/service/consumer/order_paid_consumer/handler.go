package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) HandleOrderPaid(ctx context.Context, message kafka.Message) error {
	event, err := s.decoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "failed to decode order paid event", zap.Error(err))
		return err
	}
	logger.Info(ctx, "order paid event", zap.Any("event", event))
	err = s.telegramService.SendPaidNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "failed to send paid notification", zap.Error(err))
		return err
	}
	return nil
}
