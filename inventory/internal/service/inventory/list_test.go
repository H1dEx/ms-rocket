package inventory

import (
	"errors"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
)

func (s *ServiceSuite) TestGetListSuccess() {
	partOne := s.GenPart()
	partTwo := s.GenPart()
	partThree := s.GenPart()

	uuids := []string{partOne.Uuid, partTwo.Uuid, partThree.Uuid}
	result := []model.Part{partOne, partTwo, partThree}
	filter := &model.PartFilter{Uuids: uuids}

	s.repo.On("ListParts", s.ctx, filter).Return(result, nil).Once()

	res, err := s.service.GetList(s.ctx, filter)

	s.NoError(err)
	s.Equal(res, result)
}

func (s *ServiceSuite) TestGetListRepoError() {
	filter := &model.PartFilter{}
	ErrTest := errors.New("test error")

	s.repo.On("ListParts", s.ctx, filter).Return(nil, ErrTest)

	res, err := s.service.GetList(s.ctx, filter)
	s.Error(err)
	s.ErrorIs(err, ErrTest)
	s.Empty(res)
}
