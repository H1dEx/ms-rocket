package v1

import (
	"github.com/H1dEx/ms-rocket/inventory/internal/service"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	inventoryService service.InventoryService
}

func NewApi(service service.InventoryService) *api {
	return &api{inventoryService: service}
}
