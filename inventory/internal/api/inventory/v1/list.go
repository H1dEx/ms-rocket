package v1

import (
	"context"

	"github.com/H1dEx/ms-rocket/inventory/internal/converter"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)
func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.inventoryService.GetList(ctx, converter.PartFilterToModel(req.GetFilter()))

	if err != nil {
		return nil, err
	}

	response := make([]*inventoryV1.Part, 0, len(parts))

	for _, part := range parts {
		response = append(response, converter.PartToProto(part))
	}

	return &inventoryV1.ListPartsResponse{Parts: response}, nil
}