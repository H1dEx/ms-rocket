package inventory

import "github.com/brianvoe/gofakeit/v7"

func (r *RepositorySuite) TestGetPartNotFound() {
	var uuid = gofakeit.UUID()

	part, err := r.repo.GetPart(r.ctx, uuid)
	r.Error(err)
	r.Empty(part)
}