package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMock "github.com/H1dEx/ms-rocket/order/internal/client/grpc/mocks"
	"github.com/H1dEx/ms-rocket/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx          context.Context
	repo         *mocks.OrderRepository
	paymentCli   *clientMock.PaymentClient
	inventoryCli *clientMock.InventoryClient
	service      *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewOrderRepository(s.T())
	s.paymentCli = clientMock.NewPaymentClient(s.T())
	s.inventoryCli = clientMock.NewInventoryClient(s.T())
	s.service = NewOrderService(s.repo, s.inventoryCli, s.paymentCli)
}

func (s *ServiceSuite) TearDownTest() {}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
