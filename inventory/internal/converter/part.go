package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func ToProtoPart(part model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      ToProtoCategory(part.Category),
		Dimensions:    ToProtoDimensions(part.Dimensions),
		Manufacturer:  ToProtoManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ToProtoValueMap(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}
}

func ToProtoParts(parts []model.Part) []*inventoryV1.Part {
	protoParts := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, ToProtoPart(part))
	}
	return protoParts
}

func ToModelPart(part *inventoryV1.Part) model.Part {
	return model.Part{
		Uuid:          part.GetUuid(),
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      ToModelCategory(part.GetCategory()),
		Dimensions:    ToModelDimensions(part.GetDimensions()),
		Manufacturer:  ToModelManufacturer(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      ToModelValueMap(part.GetMetadata()),
		CreatedAt:     part.GetCreatedAt().AsTime(),
		UpdatedAt:     part.GetUpdatedAt().AsTime(),
	}
}

func ToProtoCategory(category model.Category) inventoryV1.Category {
	return inventoryV1.Category(category)
}

func ToModelCategory(category inventoryV1.Category) model.Category {
	return model.Category(category)
}

func ToProtoDimensions(dimensions *model.Dimensions) *inventoryV1.Dimensions {
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

func ToModelDimensions(dimensions *inventoryV1.Dimensions) *model.Dimensions {
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

func ToProtoManufacturer(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToModelManufacturer(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    manufacturer.GetName(),
		Country: manufacturer.GetCountry(),
		Website: manufacturer.GetWebsite(),
	}
}

func ToProtoValue(value *model.Value) *inventoryV1.Value {
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

func ToModelValue(value *inventoryV1.Value) *model.Value {
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

func ToProtoValueMap(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ToProtoValue(value)
	}
	return result
}

func ToModelValueMap(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ToModelValue(value)
	}
	return result
}

func ToModelPartsFilter(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{
			Uuids:                 []string{},
			Names:                 []string{},
			Categories:            []model.Category{},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}
	}

	categories := make([]model.Category, 0, len(filter.GetCategories()))
	for _, cat := range filter.GetCategories() {
		categories = append(categories, ToModelCategory(cat))
	}

	// Ensure all fields are non-nil slices
	uuids := filter.GetUuids()
	if uuids == nil {
		uuids = []string{}
	}

	names := filter.GetNames()
	if names == nil {
		names = []string{}
	}

	countries := filter.GetManufacturerCountries()
	if countries == nil {
		countries = []string{}
	}

	tags := filter.GetTags()
	if tags == nil {
		tags = []string{}
	}

	return model.PartsFilter{
		Uuids:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: countries,
		Tags:                  tags,
	}
}
