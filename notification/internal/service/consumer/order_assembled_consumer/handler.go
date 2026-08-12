package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) HandleOrderAssembled(ctx context.Context, message kafka.Message) error {
	event, err := s.decoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "failed to decode order assembled event", zap.Error(err))
		return err
	}
	logger.Info(ctx, "order assembled event", zap.Any("event", event))
	err = s.telegramService.SendAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "failed to send assembled notification", zap.Error(err))
		return err
	}
	return nil
}
