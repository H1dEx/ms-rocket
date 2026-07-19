package inventory

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartSuccess() {
	part := s.GenPart()

	s.repo.On("GetPart", s.ctx, part.Uuid).Return(part, nil).Once()

	res, err := s.service.GetPart(s.ctx, part.Uuid)

	s.NoError(err)
	s.NotEmpty(res)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestGetPartNotFound() {
	uuid := gofakeit.UUID()
	s.repo.On("GetPart", s.ctx, uuid).Return(model.Part{}, model.ErrPartNotFound)

	res, err := s.service.GetPart(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, model.ErrPartNotFound)
	s.Empty(res)
}
