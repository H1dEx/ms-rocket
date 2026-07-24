package repository

import (
	"context"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
)

type InventoryRepository interface {
	GetPart(ctx context.Context, partUUID string) (model.Part, error)
	ListParts(ctx context.Context, filter *model.PartFilter) ([]model.Part, error)
}
