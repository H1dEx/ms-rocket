//go:build integration

package integration

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	uuid := gofakeit.UUID()
	partDoc := bson.M{
		"_id":         uuid,
		"name":        gofakeit.Word(),
		"description": gofakeit.Word(),
		"price":       uint64(gofakeit.Uint32()),
		"quantity":    int64(gofakeit.Number(1, 1000)),
		"category":    "ENGINE",
		"dimensions": bson.M{
			"length": gofakeit.Float64(),
			"width":  gofakeit.Float64(),
			"height": gofakeit.Float64(),
			"weight": gofakeit.Float64(),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Word(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags": []string{gofakeit.Word()},
		"metadata": bson.M{
			gofakeit.Word(): bson.M{
				"stringvalue": gofakeit.Word(),
			},
		},
		"created_at": time.Now().UTC(),
	}

	_, err := env.Mongo.Client().
		Database(env.Mongo.Config().Database).
		Collection(partsCollectionName).
		InsertOne(ctx, partDoc)
	if err != nil {
		return "", fmt.Errorf("failed to insert test part: %w", err)
	}
	return uuid, nil
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	_, err := env.Mongo.Client().
		Database(env.Mongo.Config().Database).
		Collection(partsCollectionName).
		DeleteMany(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to clear parts collection: %w", err)
	}
	return nil
}
