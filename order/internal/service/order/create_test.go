package order

import (
	"errors"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"

	"github.com/H1dEx/ms-rocket/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	var (
		orderUUID  = gofakeit.UUID()
		userUUID   = gofakeit.UUID()
		partUUIDS  = []string{gofakeit.UUID(), gofakeit.UUID()}
		parts      = []model.Part{{Uuid: partUUIDS[0], Price: 5}, {Uuid: partUUIDS[1], Price: 5}}
		totalPrice = float32(10)
		order      = model.Order{
			OrderUUID:  orderUUID,
			UserUUID:   userUUID,
			PartUuids:  partUUIDS,
			TotalPrice: float32(totalPrice),
		}
	)

	s.inventoryCli.On("ListParts", s.ctx, partUUIDS).Return(parts, nil).Once()
	s.repo.On("CreateOrder", s.ctx, mock.Anything, userUUID, partUUIDS, totalPrice).Return(nil).Once()
	s.repo.On("GetOrder", s.ctx, mock.Anything).Return(order, nil).Once()

	responce, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDS)

	s.NoError(err)
	s.Equal(responce, order)
}

func (s *ServiceSuite) TestCreateOrderNotFoundErr() {
	var (
		userUUID   = gofakeit.UUID()
		partUUIDS  = []string{gofakeit.UUID(), gofakeit.UUID()}
		parts      = []model.Part{{Uuid: partUUIDS[0], Price: 5}, {Uuid: partUUIDS[1], Price: 5}}
		totalPrice = float32(10)
	)

	s.inventoryCli.On("ListParts", s.ctx, partUUIDS).Return(parts, nil).Once()
	s.repo.On("CreateOrder", s.ctx, mock.Anything, userUUID, partUUIDS, totalPrice).Return(nil).Once()
	s.repo.On("GetOrder", s.ctx, mock.Anything).Return(model.Order{}, model.ErrOrderNotFound).Once()

	responce, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDS)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderNotFound)
	s.Empty(responce)
}

func (s *ServiceSuite) TestCreateOrderCreateErr() {
	var (
		userUUID    = gofakeit.UUID()
		partUUIDS   = []string{gofakeit.UUID(), gofakeit.UUID()}
		parts       = []model.Part{{Uuid: partUUIDS[0], Price: 5}, {Uuid: partUUIDS[1], Price: 5}}
		totalPrice  = float32(10)
		ErrCreating = errors.New("creating error")
	)

	s.inventoryCli.On("ListParts", s.ctx, partUUIDS).Return(parts, nil).Once()
	s.repo.On("CreateOrder", s.ctx, mock.Anything, userUUID, partUUIDS, totalPrice).Return(ErrCreating).Once()

	responce, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDS)

	s.Error(err)
	s.ErrorIs(err, ErrCreating)
	s.Empty(responce)
}

func (s *ServiceSuite) TestCreateOrderGetLessPartsErr() {
	var (
		userUUID  = gofakeit.UUID()
		partUUIDS = []string{gofakeit.UUID(), gofakeit.UUID()}
		parts     = []model.Part{{Uuid: partUUIDS[0], Price: 5}, {Uuid: partUUIDS[1], Price: 5}}
	)

	s.inventoryCli.On("ListParts", s.ctx, partUUIDS).Return(parts[:1], nil).Once()

	responce, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDS)

	s.Error(err)
	s.ErrorIs(err, model.ErrPartsNotFound)
	s.Empty(responce)
}

func (s *ServiceSuite) TestCreateOrderGettingPartsErr() {
	var (
		userUUID        = gofakeit.UUID()
		partUUIDS       = []string{gofakeit.UUID(), gofakeit.UUID()}
		ErrGettingParts = errors.New("getting parts error")
	)

	s.inventoryCli.On("ListParts", s.ctx, partUUIDS).Return(nil, ErrGettingParts).Once()

	responce, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDS)

	s.Error(err)
	s.ErrorIs(err, ErrGettingParts)
	s.Empty(responce)
}
