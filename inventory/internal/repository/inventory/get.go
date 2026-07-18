package inventory

import (
	"context"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/H1dEx/ms-rocket/inventory/internal/repository/converter"
)

func (r *repository) GetPart(ctx context.Context, partUuid string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	part, ok := r.parts[partUuid]

	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return converter.PartToModel(part), nil
}
