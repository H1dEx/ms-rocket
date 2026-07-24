package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/H1dEx/ms-rocket/inventory/internal/service/mocks"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

type APISuite struct {
	suite.Suite

	ctx     context.Context
	service *mocks.InventoryService
	api     inventoryV1.InventoryServiceServer
}

func (a *APISuite) SetupTest() {
	a.ctx = context.Background()
	a.service = mocks.NewInventoryService(a.T())
	a.api = NewAPI(a.service)
}

func (a *APISuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
