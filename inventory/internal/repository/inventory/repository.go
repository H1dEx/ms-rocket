package inventory

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/H1dEx/ms-rocket/inventory/internal/repository"
	"github.com/H1dEx/ms-rocket/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) def.InventoryRepository {
	collection := db.Collection("parts")

	test := model.Part{
		UUID:          uuid.New().String(),
		Name:          "Init part",
		Description:   "Init part",
		Price:         100,
		StockQuantity: 100,
		Category:      "Test",
		Dimensions:    model.Dimensions{Width: 100, Height: 100},
		Manufacturer:  model.Manufacturer{Name: "Test", Country: "Test", Website: "Test"},
		Tags:          []string{"Test"},
		Metadata: map[string]model.MetadataValue{
			"Test": {
				StringValue: lo.ToPtr("Test"),
			},
		},
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), test)
	if err != nil {
		log.Printf("failed to insert test: %v\n", err)
		panic(err)
	}
	return &repository{
		collection: collection,
	}
}
