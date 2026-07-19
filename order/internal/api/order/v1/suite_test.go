package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/H1dEx/ms-rocket/order/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite
	ctx     context.Context
	service *mocks.OrderService
	api     *api
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()
	a.service = mocks.NewOrderService(a.T())
	a.api = NewOrderApi(a.service)
}

func (s *ApiSuite) TearDownTest() {}

func TestService(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
