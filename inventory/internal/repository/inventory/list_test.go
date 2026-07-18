package inventory

import (
	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (r *RepositorySuite) TestListPartsEmptyResponse() {
	var uuid = gofakeit.UUID()

	parts, err := r.repo.ListParts(r.ctx, &model.PartFilter{Uuids: []string{uuid}})

	r.NoError(err)
	r.Empty(parts)
}
