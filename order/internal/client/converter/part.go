package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func ToModelPartList(parts []*inventoryV1.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, ToModelPart(part))
	}
	return result
}

func ToModelPart(part *inventoryV1.Part) model.Part {
	partUUID, err := uuid.Parse(part.GetUuid())
	if err != nil {
		// If UUID is invalid, keep zero value to indicate missing/invalid ID.
		partUUID = uuid.UUID{}
	}

	var createdAt, updatedAt time.Time
	if part.GetCreatedAt() != nil {
		createdAt = part.GetCreatedAt().AsTime()
	}
	if part.GetUpdatedAt() != nil {
		updatedAt = part.GetUpdatedAt().AsTime()
	}

	return model.Part{
		Uuid:          partUUID,
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      ToModelCategory(part.GetCategory()),
		Dimensions:    ToModelDimensions(part.GetDimensions()),
		Manufacturer:  ToModelManufacturer(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      ToModelPartMetadata(part.GetMetadata()),
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
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

func ToModelPartMetadata(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ToModelValue(value)
	}
	return result
}

func ToModelValue(value *inventoryV1.Value) *model.Value {
	if value == nil {
		return nil
	}

	modelValue := &model.Value{}

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

func ToModelCategory(category inventoryV1.Category) model.Category {
	return model.Category(category)
}

func ToProtoPartsFilter(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, cat := range filter.Categories {
		categories = append(categories, inventoryV1.Category(cat))
	}

	uuids := make([]string, 0, len(filter.Uuids))
	for _, uuid := range filter.Uuids {
		uuids = append(uuids, uuid.String())
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
