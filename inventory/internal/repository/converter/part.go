package converter

import (
	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	repoModel "github.com/H1dEx/ms-rocket/inventory/internal/repository/model"
)

func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToRepoModel(part.Category),
		Dimensions:    DimensionsToRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func MetadataValueToRepoModel(v *model.MetadataValue) *repoModel.MetadataValue {
	if v == nil {
		return nil
	}
	return &repoModel.MetadataValue{
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}
}

func MetadataToRepoModel(m map[string]*model.MetadataValue) map[string]*repoModel.MetadataValue {
	metadata := make(map[string]*repoModel.MetadataValue, len(m))
	for k, v := range m {
		metadata[k] = MetadataValueToRepoModel(v)
	}
	return metadata
}

func ManufacturerToRepoModel(m *model.Manufacturer) *repoModel.Manufacturer {
	if m == nil {
		return nil
	}
	return &repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func DimensionsToRepoModel(dim *model.Dimensions) *repoModel.Dimensions {
	if dim == nil {
		return nil
	}
	return &repoModel.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func CategoryToRepoModel(category model.Category) repoModel.Category {
	switch category {
	case model.CategoryUnknown:
		return repoModel.CategoryUnknown
	case model.CategoryEngine:
		return repoModel.CategoryEngine
	case model.CategoryFuel:
		return repoModel.CategoryFuel
	case model.CategoryPorthole:
		return repoModel.CategoryPorthole
	case model.CategoryWing:
		return repoModel.CategoryWing
	}
	return repoModel.CategoryUnknown
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func MetadataValueToModel(v *repoModel.MetadataValue) *model.MetadataValue {
	if v == nil {
		return nil
	}
	return &model.MetadataValue{
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}
}

func MetadataToModel(m map[string]*repoModel.MetadataValue) map[string]*model.MetadataValue {
	metadata := make(map[string]*model.MetadataValue, len(m))
	for k, v := range m {
		metadata[k] = MetadataValueToModel(v)
	}
	return metadata
}

func ManufacturerToModel(m *repoModel.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func DimensionsToModel(dim *repoModel.Dimensions) *model.Dimensions {
	if dim == nil {
		return nil
	}
	return &model.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func CategoryToModel(category repoModel.Category) model.Category {
	switch category {
	case repoModel.CategoryUnknown:
		return model.CategoryUnknown
	case repoModel.CategoryEngine:
		return model.CategoryEngine
	case repoModel.CategoryFuel:
		return model.CategoryFuel
	case repoModel.CategoryPorthole:
		return model.CategoryPorthole
	case repoModel.CategoryWing:
		return model.CategoryWing
	}
	return model.CategoryUnknown
}

func CategoriesToRepoModel(c []model.Category) []repoModel.Category {
	categories := make([]repoModel.Category, 0, len(c))
	for _, v := range c {
		categories = append(categories, CategoryToRepoModel(v))
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
