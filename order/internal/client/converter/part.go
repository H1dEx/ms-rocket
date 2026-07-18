package converter

import (
	"time"

	"github.com/H1dEx/ms-rocket/order/internal/model"
	inventory_v1 "github.com/H1dEx/ms-rocket/shared/pkg/proto/inventory/v1"
	"github.com/samber/lo"
)

func CategoryToModel(category inventory_v1.Category) model.Category {
	switch category {
	case inventory_v1.Category_CATEGORY_UNKNOWN:
		return model.CategoryUnknown
	case inventory_v1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventory_v1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventory_v1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventory_v1.Category_CATEGORY_WING:
		return model.CategoryWing
	}
	return model.CategoryUnknown
}


func DimensionsToModel(dim *inventory_v1.Dimensions) *model.Dimensions {
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

func ManufacturerToModel(m *inventory_v1.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func MetadataValueToModel(v *inventory_v1.Value) *model.MetadataValue {
	if v == nil {
		return nil
	}
	switch v.ValueType.(type){
	case *inventory_v1.Value_StringValue:
		return &model.MetadataValue{StringValue: lo.ToPtr(v.GetStringValue())}
	case *inventory_v1.Value_DoubleValue:
		return &model.MetadataValue{DoubleValue: lo.ToPtr(v.GetDoubleValue())}
	case *inventory_v1.Value_BoolValue:
		return &model.MetadataValue{BoolValue: lo.ToPtr(v.GetBoolValue())}
	case *inventory_v1.Value_Int64Value:
		return &model.MetadataValue{Int64Value: lo.ToPtr(v.GetInt64Value())}
	}
	return nil
}

func MetadataToRepoModel(m map[string]*inventory_v1.Value) map[string]*model.MetadataValue {
	metadata := make(map[string]*model.MetadataValue, len(m))
	for k, v := range m {
		metadata[k] = MetadataValueToModel(v)
	}
	return metadata
}

func InventoryPartToModel(part *inventory_v1.Part) model.Part {
	var createdAt *time.Time  
	if part.CreatedAt != nil {
		createdAt = lo.ToPtr(part.CreatedAt.AsTime())
	}
	var updatedAt *time.Time  
	if part.UpdatedAt != nil {
		createdAt = lo.ToPtr(part.UpdatedAt.AsTime())
	}
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
		Metadata:      MetadataToRepoModel(part.Metadata),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}
}
