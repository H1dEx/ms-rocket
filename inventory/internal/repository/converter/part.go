package converter

import (
	"fmt"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	repoModel "github.com/H1dEx/ms-rocket/inventory/internal/repository/model"
)

func PartToRepoModel(part model.Part) (repoModel.Part, error) {
	if part.CreatedAt == nil {
		return repoModel.Part{}, fmt.Errorf("part created at is required")
	}
	return repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      string(part.Category),
		Dimensions:    DimensionsToRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     *part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}, nil
}

func MetadataValueToRepoModel(v model.MetadataValue) repoModel.MetadataValue {
	return repoModel.MetadataValue{
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}
}

func MetadataToRepoModel(m map[string]model.MetadataValue) map[string]repoModel.MetadataValue {
	metadata := make(map[string]repoModel.MetadataValue, len(m))
	for k, v := range m {
		metadata[k] = MetadataValueToRepoModel(v)
	}
	return metadata
}

func ManufacturerToRepoModel(m model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func DimensionsToRepoModel(dim model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     &part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func MetadataValueToModel(v repoModel.MetadataValue) model.MetadataValue {
	return model.MetadataValue{
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}
}

func MetadataToModel(m map[string]repoModel.MetadataValue) map[string]model.MetadataValue {
	metadata := make(map[string]model.MetadataValue, len(m))
	for k, v := range m {
		metadata[k] = MetadataValueToModel(v)
	}
	return metadata
}

func ManufacturerToModel(m repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func DimensionsToModel(dim repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func CategoryToModel(category string) model.Category {
	switch category {
	case string(model.CategoryEngine):
		return model.CategoryEngine
	case string(model.CategoryFuel):
		return model.CategoryFuel
	case string(model.CategoryPorthole):
		return model.CategoryPorthole
	case string(model.CategoryWing):
		return model.CategoryWing
	}
	return model.CategoryUnknown
}

func CategoriesToRepoModel(c []model.Category) []string {
	categories := make([]string, 0, len(c))
	for _, v := range c {
		categories = append(categories, string(v))
	}
	return categories
}

func PartFilterToRepoModel(f *model.PartFilter) *repoModel.PartFilter {
	if f == nil {
		return nil
	}
	return &repoModel.PartFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Categories:            CategoriesToRepoModel(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}
