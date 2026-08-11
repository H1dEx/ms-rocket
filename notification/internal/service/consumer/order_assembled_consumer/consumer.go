package order_assembled_consumer

import (
	"context"

	decoder "github.com/H1dEx/ms-rocket/notification/internal/converter/kafka"
	def "github.com/H1dEx/ms-rocket/notification/internal/service"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
)

var _ def.OrderAssembledService = (*service)(nil)

type service struct {
	telegramService        def.TelegramService
	orderAssembledConsumer kafka.Consumer
	decoder                decoder.OrderAssembledDecoder
}

func NewService(orderAssembledConsumer kafka.Consumer, decoder decoder.OrderAssembledDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		decoder:                decoder,
		telegramService:        telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	err := s.orderAssembledConsumer.Consume(ctx, s.HandleOrderAssembled)
	if err != nil {
		return err
	}
	return nil
}
