package order_producer

import (
	"context"

	"github.com/H1dEx/ms-rocket/order/internal/converter"
	"github.com/H1dEx/ms-rocket/order/internal/model"
	def "github.com/H1dEx/ms-rocket/order/internal/service"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var _ def.ProducerService = (*service)(nil)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (s *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := converter.OrderPaidEventToModel(event)

	encoded, err := proto.Marshal(&msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal order paid event", zap.Error(err))
		return err
	}

	err = s.orderProducer.Send(ctx, []byte(event.EventUUID), encoded)
	if err != nil {
		logger.Error(ctx, "failed to send order paid event", zap.Error(err))
		return err
	}

	return nil
}
