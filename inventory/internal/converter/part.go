package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

// PartToProto конвертирует model.Part в inventoryV1.Part
func PartToProto(part model.Part) *inventoryV1.Part {
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
		Metadata:      ValueMapToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}
}

// PartToModel конвертирует inventoryV1.Part в model.Part
func PartToModel(part *inventoryV1.Part) model.Part {
	return model.Part{
		Uuid:          part.GetUuid(),
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      CategoryToModel(part.GetCategory()),
		Dimensions:    DimensionsToModel(part.GetDimensions()),
		Manufacturer:  ManufacturerToModel(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      ValueMapToModel(part.GetMetadata()),
		CreatedAt:     part.GetCreatedAt().AsTime(),
		UpdatedAt:     part.GetUpdatedAt().AsTime(),
	}
}

// CategoryToProto конвертирует model.Category в inventoryV1.Category
func CategoryToProto(category model.Category) inventoryV1.Category {
	return inventoryV1.Category(category)
}

// CategoryToModel конвертирует inventoryV1.Category в model.Category
func CategoryToModel(category inventoryV1.Category) model.Category {
	return model.Category(category)
}

// DimensionsToProto конвертирует model.Dimensions в inventoryV1.Dimensions
func DimensionsToProto(dimensions *model.Dimensions) *inventoryV1.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// DimensionsToModel конвертирует inventoryV1.Dimensions в model.Dimensions
func DimensionsToModel(dimensions *inventoryV1.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &model.Dimensions{
		Length: dimensions.GetLength(),
		Width:  dimensions.GetWidth(),
		Height: dimensions.GetHeight(),
		Weight: dimensions.GetWeight(),
	}
}

// ManufacturerToProto конвертирует model.Manufacturer в inventoryV1.Manufacturer
func ManufacturerToProto(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// ManufacturerToModel конвертирует inventoryV1.Manufacturer в model.Manufacturer
func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    manufacturer.GetName(),
		Country: manufacturer.GetCountry(),
		Website: manufacturer.GetWebsite(),
	}
}

// ValueToProto конвертирует model.Value в inventoryV1.Value
func ValueToProto(value *model.Value) *inventoryV1.Value {
	if value == nil {
		return nil
	}

	protoValue := &inventoryV1.Value{}

	// Устанавливаем значение в зависимости от того, какое поле заполнено
	switch {
	case value.StringValue != "":
		protoValue.Value = &inventoryV1.Value_StringValue{StringValue: value.StringValue}
	case value.Int64Value != 0:
		protoValue.Value = &inventoryV1.Value_Int64Value{Int64Value: value.Int64Value}
	case value.DoubleValue != 0:
		protoValue.Value = &inventoryV1.Value_DoubleValue{DoubleValue: value.DoubleValue}
	default:
		protoValue.Value = &inventoryV1.Value_BoolValue{BoolValue: value.BoolValue}
	}

	return protoValue
}

// ValueToModel конвертирует inventoryV1.Value в model.Value
func ValueToModel(value *inventoryV1.Value) *model.Value {
	if value == nil {
		return nil
	}

	modelValue := &model.Value{}

	// Извлекаем значение в зависимости от типа oneof поля
	switch v := value.GetValue().(type) {
	case *inventoryV1.Value_StringValue:
		modelValue.StringValue = v.StringValue
	case *inventoryV1.Value_Int64Value:
		modelValue.Int64Value = v.Int64Value
	case *inventoryV1.Value_DoubleValue:
		modelValue.DoubleValue = v.DoubleValue
	case *inventoryV1.Value_BoolValue:
		modelValue.BoolValue = v.BoolValue
	}

	return modelValue
}

// ValueMapToProto конвертирует map[string]*model.Value в map[string]*inventoryV1.Value
func ValueMapToProto(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ValueToProto(value)
	}
	return result
}

// ValueMapToModel конвертирует map[string]*inventoryV1.Value в map[string]*model.Value
func ValueMapToModel(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ValueToModel(value)
	}
	return result
}

// PartsFilterToModel конвертирует inventoryV1.PartsFilter в model.PartsFilter
func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	categories := make([]model.Category, 0, len(filter.GetCategories()))
	for _, cat := range filter.GetCategories() {
		categories = append(categories, CategoryToModel(cat))
	}

	return model.PartsFilter{
		Uuids:                 filter.GetUuids(),
		Names:                 filter.GetNames(),
		Categories:            categories,
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}
