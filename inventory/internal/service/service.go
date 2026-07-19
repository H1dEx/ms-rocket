package service

import (
	"context"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
)

type InventoryService interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	GetList(ctx context.Context, filter *model.PartFilter) ([]model.Part, error)
}
