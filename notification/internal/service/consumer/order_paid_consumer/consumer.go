package order_paid_consumer

import (
	"context"

	decoder "github.com/H1dEx/ms-rocket/notification/internal/converter/kafka"
	def "github.com/H1dEx/ms-rocket/notification/internal/service"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
)

var _ def.OrderPaidService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	decoder           decoder.OrderPaidDecoder
	telegramService   def.TelegramService
}

func NewService(orderPaidConsumer kafka.Consumer, decoder decoder.OrderPaidDecoder, telegramService def.TelegramService) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		decoder:           decoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	err := s.orderPaidConsumer.Consume(ctx, s.HandleOrderPaid)
	if err != nil {
		return err
	}
	return nil
}
