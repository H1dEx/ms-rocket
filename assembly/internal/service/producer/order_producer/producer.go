package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/H1dEx/ms-rocket/assembly/internal/converter"
	"github.com/H1dEx/ms-rocket/assembly/internal/model"
	def "github.com/H1dEx/ms-rocket/assembly/internal/service"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

var _ def.ProducerService = (*service)(nil)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{orderProducer}
}

func (s *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := converter.ShipAssembledEventToEvent(event)

	encoded, err := proto.Marshal(&msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ship assembled event", zap.Error(err))
		return err
	}

	err = s.orderProducer.Send(ctx, []byte(event.EventUUID), encoded)
	if err != nil {
		logger.Error(ctx, "failed to send ship assembled event", zap.Error(err))
		return err
	}
	return nil
}
