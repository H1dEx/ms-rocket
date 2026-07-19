package inventory

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/suite"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/H1dEx/ms-rocket/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context

	service *service
	repo    *mocks.InventoryRepository
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewInventoryRepository(s.T())
	s.service = NewService(s.repo)
}

func (s *ServiceSuite) TearDownTest() {}

func (s *ServiceSuite) GenPart() model.Part {
	return model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.Comment(),
		Price:         gofakeit.Uint64(),
		StockQuantity: gofakeit.Int64(),
	}
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
