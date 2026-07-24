package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/H1dEx/ms-rocket/order/internal/service/mocks"
	orderV1 "github.com/H1dEx/ms-rocket/shared/pkg/openapi/order/v1"
)

type APISuite struct {
	suite.Suite
	ctx     context.Context
	service *mocks.OrderService
	api     orderV1.Handler
}

func (a *APISuite) SetupTest() {
	a.ctx = context.Background()
	a.service = mocks.NewOrderService(a.T())
	a.api = NewOrderAPI(a.service)
}

func (a *APISuite) TearDownTest() {}

func TestService(t *testing.T) {
	suite.Run(t, new(APISuite))
}
