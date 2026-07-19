package inventory

import (
	"github.com/H1dEx/ms-rocket/inventory/internal/repository"
	def "github.com/H1dEx/ms-rocket/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	repo repository.InventoryRepository
}

func NewService(repo repository.InventoryRepository) *service {
	return &service{repo}
}
