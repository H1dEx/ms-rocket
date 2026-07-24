package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/H1dEx/ms-rocket/payment/internal/service/mocks"
	paymentV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/payment/v1"
)

type APISuite struct {
	suite.Suite
	ctx     context.Context
	service *mocks.PaymentService

	api paymentV1.PaymentServiceServer
}

func (a *APISuite) SetupTest() {
	a.ctx = context.Background()
	a.service = mocks.NewPaymentService(a.T())
	a.api = NewAPI(a.service)
}

func (a *APISuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
