package inventory

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/H1dEx/ms-rocket/inventory/internal/repository/converter"
	repoModel "github.com/H1dEx/ms-rocket/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, partUUID string) (model.Part, error) {
	var part repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"_id": partUUID}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}

	return converter.PartToModel(part), nil
}
