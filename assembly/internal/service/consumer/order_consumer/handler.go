package order_consumer

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/assembly/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) HandleOrderPaid(ctx context.Context, message kafka.Message) error {
	event, err := s.decoder.Decode(message.Value)
	if err != nil {
		logger.Error(ctx, "Error decoding order paid", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order paid event", zap.Any("event", event))
	buildTimeSeconds := rand.Intn(10) + 1

	select {
	case <-time.After(time.Duration(buildTimeSeconds) * time.Second):
	case <-ctx.Done():
		logger.Info(ctx, "Context cancelled, skipping order assembly")
		return nil
	}
	return s.producer.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		EventUUID:        event.EventUUID,
		OrderUUID:        event.OrderUUID,
		UserUUID:         event.UserUUID,
		BuildTimeSeconds: int64(buildTimeSeconds),
	})
}
