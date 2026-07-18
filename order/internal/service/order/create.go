package order

import (
	"context"
	"fmt"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	"github.com/google/uuid"
)

func findMissingId(parts []model.Part, ids []string) []string {
	foundMap := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		foundMap[part.Uuid] = struct{}{}
	}

	notFound := []string{}
	for _, id := range ids {
		if _, ok := foundMap[id]; !ok {
			notFound = append(notFound, id)
		}
	}

	return notFound
}

func (s *service) CreateOrder(ctx context.Context, userUUID string, partUUIDs []string) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, partUUIDs)
	if err != nil {
		return model.Order{}, err
	}
	if len(parts) < len(partUUIDs) {
		ids := findMissingId(parts, partUUIDs)
		return model.Order{}, fmt.Errorf("Not found details with uuids %v : %w", ids, model.ErrPartsNotFound)
	}

	var sum float32
	for _, p := range parts {
		sum += float32(p.Price)
	}
	uuid := uuid.New().String()

	err = s.repo.CreateOrder(ctx, uuid, userUUID, partUUIDs, sum)

	if err != nil {
		return model.Order{}, err
	}
	order, err := s.repo.GetOrder(ctx, uuid)

	if err != nil {
		return model.Order{}, err
	}
	return order, nil
}
