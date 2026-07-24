package v1

import (
	"github.com/H1dEx/ms-rocket/inventory/internal/service"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	inventoryService service.InventoryService
}

func NewAPI(svc service.InventoryService) inventoryV1.InventoryServiceServer {
	return &api{inventoryService: svc}
}
