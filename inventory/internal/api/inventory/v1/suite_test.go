package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/H1dEx/ms-rocket/inventory/internal/service/mocks"
)

type ApiSuite struct {
	suite.Suite

	ctx     context.Context
	service *mocks.InventoryService
	api     *api
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()
	a.service = mocks.NewInventoryService(a.T())
	a.api = NewApi(a.service)
}

func (a *ApiSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
