package order

import (
	"github.com/H1dEx/ms-rocket/order/internal/client/grpc"
	"github.com/H1dEx/ms-rocket/order/internal/repository"
	def "github.com/H1dEx/ms-rocket/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	repo repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient grpc.PaymentClient
}

func NewOrderService(repo repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		repo: repo,
		inventoryClient: inventoryClient,
		paymentClient: paymentClient,
	}
}