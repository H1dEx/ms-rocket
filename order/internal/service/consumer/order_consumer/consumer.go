package order_consumer

import (
	"context"

	decoder "github.com/H1dEx/ms-rocket/order/internal/converter/kafka"
	"github.com/H1dEx/ms-rocket/order/internal/repository"
	"github.com/H1dEx/ms-rocket/platform/pkg/kafka"
)

type service struct {
	orderConsumer kafka.Consumer
	decoder       decoder.OrderPaidDecoder
	repo          repository.OrderRepository
}

func NewService(repo repository.OrderRepository, orderConsumer kafka.Consumer, decoder decoder.OrderPaidDecoder) *service {
	return &service{
		orderConsumer: orderConsumer,
		decoder:       decoder,
		repo:          repo,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	err := s.orderConsumer.Consume(ctx, s.HandleOrderAssembled)
	if err != nil {
		return err
	}
	return nil
}
