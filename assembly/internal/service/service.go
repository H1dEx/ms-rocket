package service

import (
	"context"

	"github.com/H1dEx/ms-rocket/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
