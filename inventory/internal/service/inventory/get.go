package inventory

import (
	"context"

	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	part, err := s.repo.GetPart(ctx, uuid)
	if err != nil {
		logger.Debug(ctx, "failed to get part", zap.Error(err))
		return model.Part{}, err
	}

	return part, nil
}
