package inventory

import (
	"context"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) GetList(ctx context.Context, filter *model.PartFilter) ([]model.Part, error) {
	parts, err := s.repo.ListParts(ctx, filter)
	if err != nil {
		logger.Error(ctx, "failed to list parts", zap.Error(err))
		return nil, err
	}

	return parts, err
}
