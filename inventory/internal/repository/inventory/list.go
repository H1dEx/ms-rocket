package inventory

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	"github.com/H1dEx/ms-rocket/inventory/internal/repository/converter"
	repoModel "github.com/H1dEx/ms-rocket/inventory/internal/repository/model"
	"github.com/H1dEx/ms-rocket/platform/pkg/logger"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartFilter) ([]model.Part, error) {
	f := converter.PartFilterToRepoModel(filter)
	filterSet := BuildFilterSets(f)

	cursor, err := r.collection.Find(ctx, prepareFilterSets(filterSet))
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			logger.Error(ctx, "failed to close cursor", zap.Error(err))
		}
	}()

	var parts []repoModel.Part

	if err := cursor.All(ctx, &parts); err != nil {
		return nil, err
	}

	results := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		results = append(results, converter.PartToModel(part))
	}
	return results, nil
}

type FilterSets struct {
	Uuids                 map[string]struct{}
	Names                 map[string]struct{}
	Categories            map[string]struct{}
	ManufacturerCountries map[string]struct{}
	Tags                  map[string]struct{}
}

func BuildFilterSets(f *repoModel.PartFilter) FilterSets {
	if f == nil {
		return FilterSets{}
	}

	return FilterSets{
		Uuids:                 filterToMap(f.Uuids),
		Names:                 filterToMap(f.Names),
		Categories:            filterToMap(f.Categories),
		ManufacturerCountries: filterToMap(f.ManufacturerCountries),
		Tags:                  filterToMap(f.Tags),
	}
}

func filterToMap[T comparable](f []T) map[T]struct{} {
	if f == nil {
		return nil
	}
	acc := make(map[T]struct{}, len(f))
	for _, v := range f {
		acc[v] = struct{}{}
	}
	return acc
}

func prepareFilterSets(f FilterSets) bson.M {
	filter := bson.M{}

	if len(f.Uuids) > 0 {
		ids := make([]string, 0, len(f.Uuids))
		for k := range f.Uuids {
			ids = append(ids, k)
		}
		filter["_id"] = bson.M{"$in": ids}
	}

	if len(f.Names) > 0 {
		names := make([]string, 0, len(f.Names))
		for k := range f.Names {
			names = append(names, k)
		}
		filter["name"] = bson.M{"$in": names}
	}

	if len(f.Categories) > 0 {
		cats := make([]string, 0, len(f.Categories))
		for k := range f.Categories {
			cats = append(cats, k)
		}
		filter["category"] = bson.M{"$in": cats}
	}

	if len(f.ManufacturerCountries) > 0 {
		countries := make([]string, 0, len(f.ManufacturerCountries))
		for k := range f.ManufacturerCountries {
			countries = append(countries, k)
		}
		filter["manufacturer.country"] = bson.M{"$in": countries}
	}

	if len(f.Tags) > 0 {
		tags := make([]string, 0, len(f.Tags))
		for k := range f.Tags {
			tags = append(tags, k)
		}
		filter["tags"] = bson.M{"$in": tags}
	}

	return filter
}
