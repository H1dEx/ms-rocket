package inventory

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/H1dEx/ms-rocket/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) def.InventoryRepository {
	collection := db.Collection("parts")
	return &repository{
		collection: collection,
	}
}
