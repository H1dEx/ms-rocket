package inventory

import (
	"context"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.repo.GetPart(ctx, uuid)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
