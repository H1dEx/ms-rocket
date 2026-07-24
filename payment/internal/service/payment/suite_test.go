package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	def "github.com/H1dEx/ms-rocket/payment/internal/service"
)

type ServiceSuite struct {
	suite.Suite
	ctx context.Context

	service def.PaymentService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.service = NewService()
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
