package v1

import (
	"github.com/H1dEx/ms-rocket/order/internal/client/grpc"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

var _ grpc.InventoryClient = (*client)(nil)

type client struct {
	inventoryClient inventoryV1.InventoryServiceClient
}

func NewInventoryClient(inventoryClient inventoryV1.InventoryServiceClient) *client {
	return &client{inventoryClient: inventoryClient}
}
