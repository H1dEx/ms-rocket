package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/H1dEx/ms-rocket/inventory/internal/model"
	inventoryV1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
)

func PartFilterToModel(filter *inventoryV1.PartsFilter) *model.PartFilter {
	if filter == nil {
		return &model.PartFilter{}
	}
	categories := make([]model.Category, 0, len(filter.Categories))

	for _, v := range filter.Categories {
		categories = append(categories, CategoryToModel(v))
	}

	return &model.PartFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
		Categories:            categories,
	}
}

func PartToProto(part model.Part) *inventoryV1.Part {
	var updatedAt *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAt = timestamppb.New(*part.UpdatedAt)
	}

	var createdAt *timestamppb.Timestamp

	if part.CreatedAt != nil {
		updatedAt = timestamppb.New(*part.CreatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToProto(part.Category),
		Dimensions:    DimensionsToProto(part.Dimensions),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}

func CategoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryUnknown:
		return inventoryV1.Category_CATEGORY_UNKNOWN
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	}
	return inventoryV1.Category_CATEGORY_UNKNOWN
}

func CategoryToModel(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_UNKNOWN:
		return model.CategoryUnknown
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	}
	return model.CategoryUnknown
}

func DimensionsToProto(dim *model.Dimensions) *inventoryV1.Dimensions {
	if dim == nil {
		return nil
	}

	return &inventoryV1.Dimensions{
		Length: dim.Length,
		Width:  dim.Width,
		Height: dim.Height,
		Weight: dim.Weight,
	}
}

func ManufacturerToProto(m *model.Manufacturer) *inventoryV1.Manufacturer {
	if m == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func MetadataValueToProto(v *model.MetadataValue) *inventoryV1.Value {
	if v == nil {
		return nil
	}
	switch {
	case v.BoolValue != nil:
		return &inventoryV1.Value{ValueType: &inventoryV1.Value_BoolValue{BoolValue: *v.BoolValue}}
	case v.DoubleValue != nil:
		return &inventoryV1.Value{ValueType: &inventoryV1.Value_DoubleValue{DoubleValue: *v.DoubleValue}}
	case v.Int64Value != nil:
		return &inventoryV1.Value{ValueType: &inventoryV1.Value_Int64Value{Int64Value: *v.Int64Value}}
	case v.StringValue != nil:
		return &inventoryV1.Value{ValueType: &inventoryV1.Value_StringValue{StringValue: *v.StringValue}}
	}
	return nil
}

func MetadataToRepoModel(m map[string]*model.MetadataValue) map[string]*inventoryV1.Value {
	metadata := make(map[string]*inventoryV1.Value, len(m))
	for k, v := range m {
		metadata[k] = MetadataValueToProto(v)
	}
	return metadata
}
