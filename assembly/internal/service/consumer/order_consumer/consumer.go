package order_consumer

import (
	"context"

	converter "github.com/H1dEx/ms-rocket/assembly/internal/converter/kafka"
	def "github.com/H1dEx/ms-rocket/assembly/internal/service"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderConsumer kafka.Consumer
	decoder converter.OrderPaidDecoder
	producer def.ProducerService
}

func NewService(orderConsumer kafka.Consumer, conv converter.OrderPaidDecoder, producer def.ProducerService) *service {
	return &service{orderConsumer, conv, producer}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer")
	err := s.orderConsumer.Consume(ctx, s.HandleOrderPaid)
	if err != nil {
		logger.Error(ctx, "Error consuming order paid", zap.Error(err))
		return err
	}
	return nil
}
