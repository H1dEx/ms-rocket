package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

type RepositorySuite struct {
	suite.Suite
	ctx context.Context

	repo *repository
}

func (r *RepositorySuite) SetupTest() {
	logger.SetNopLogger()
	r.ctx = context.Background()
	// r.repo = NewRepository()
}

func (r *RepositorySuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	t.Skip("inventory repository tests require MongoDB (testcontainers); skipped until wired")
	suite.Run(t, new(RepositorySuite))
}
