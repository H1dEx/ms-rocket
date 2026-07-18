package inventory

import (
	"context"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
)

func (s *service) GetList(ctx context.Context, filter *model.PartFilter) ([]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		return nil, err
	}

	return parts, err
}