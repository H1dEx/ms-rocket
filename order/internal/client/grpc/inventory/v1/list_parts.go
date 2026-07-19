package v1

import (
	"context"
	"errors"

	"github.com/H1dEx/ms-rocket/order/internal/client/converter"
	"github.com/H1dEx/ms-rocket/order/internal/model"
	inventory_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, uuids []string) ([]model.Part, error) {
	res, err := c.inventoryClient.ListParts(ctx, &inventory_v1.ListPartsRequest{Filter: &inventory_v1.PartsFilter{Uuids: uuids}})
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("empty response")
	}

	parts := make([]model.Part, 0, len(res.GetParts()))

	for _, part := range res.GetParts() {
		parts = append(parts, converter.InventoryPartToModel(part))
	}

	return parts, nil
}
