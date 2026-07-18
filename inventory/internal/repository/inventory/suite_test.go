package inventory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepositorySuite struct {
	suite.Suite
	ctx context.Context

	repo  *repository
}

func (r *RepositorySuite) SetupTest() {
	r.ctx = context.Background()
	r.repo = NewRepository()
}


func (r *RepositorySuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}
