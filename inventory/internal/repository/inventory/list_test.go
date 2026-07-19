package inventory

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
)

func (r *RepositorySuite) TestListPartsEmptyResponse() {
	uuid := gofakeit.UUID()

	parts, err := r.repo.ListParts(r.ctx, &model.PartFilter{Uuids: []string{uuid}})

	r.NoError(err)
	r.Empty(parts)
}
