package inventory

import (
	"sync"

	def "github.com/H1dEx/ms-rocket/inventory/internal/repository"
	repModel "github.com/H1dEx/ms-rocket/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu sync.RWMutex

	parts map[string]repModel.Part
}

func NewRepository() *repository {
	rep := &repository{
		parts: make(map[string]repModel.Part),
	}

	firstMock := repModel.Part{
		Uuid:  "111",
		Name:  "First detail",
		Price: 100,
	}
	secondMock := repModel.Part{
		Uuid:  "222",
		Name:  "Second detail",
		Price: 200,
	}

	rep.parts[firstMock.Uuid] = firstMock
	rep.parts[secondMock.Uuid] = secondMock
	return rep
}
