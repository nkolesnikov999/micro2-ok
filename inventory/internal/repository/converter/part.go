package converter

import (
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

// PartToRepoModel конвертирует model.Part в repository/model.Part
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
		Metadata:      ValueMapToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// PartToModel конвертирует repository/model.Part в model.Part
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
		Metadata:      ValueMapToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

// CategoryToRepoModel конвертирует model.Category в repository/model.Category
func CategoryToRepoModel(category model.Category) repoModel.Category {
	return repoModel.Category(category)
}

// CategoryToModel конвертирует repository/model.Category в model.Category
func CategoryToModel(category repoModel.Category) model.Category {
	return model.Category(category)
}

// DimensionsToRepoModel конвертирует model.Dimensions в repository/model.Dimensions
func DimensionsToRepoModel(dimensions *model.Dimensions) *repoModel.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// DimensionsToModel конвертирует repository/model.Dimensions в model.Dimensions
func DimensionsToModel(dimensions *repoModel.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}
	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// ManufacturerToRepoModel конвертирует model.Manufacturer в repository/model.Manufacturer
func ManufacturerToRepoModel(manufacturer *model.Manufacturer) *repoModel.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// ManufacturerToModel конвертирует repository/model.Manufacturer в model.Manufacturer
func ManufacturerToModel(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// ValueToRepoModel конвертирует model.Value в repository/model.Value
func ValueToRepoModel(value *model.Value) *repoModel.Value {
	if value == nil {
		return nil
	}
	return &repoModel.Value{
		StringValue: value.StringValue,
		Int64Value:  value.Int64Value,
		DoubleValue: value.DoubleValue,
		BoolValue:   value.BoolValue,
	}
}

// ValueToModel конвертирует repository/model.Value в model.Value
func ValueToModel(value *repoModel.Value) *model.Value {
	if value == nil {
		return nil
	}
	return &model.Value{
		StringValue: value.StringValue,
		Int64Value:  value.Int64Value,
		DoubleValue: value.DoubleValue,
		BoolValue:   value.BoolValue,
	}
}

// ValueMapToRepoModel конвертирует map[string]*model.Value в map[string]*repository/model.Value
func ValueMapToRepoModel(metadata map[string]*model.Value) map[string]*repoModel.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*repoModel.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ValueToRepoModel(value)
	}
	return result
}

// ValueMapToModel конвертирует map[string]*repository/model.Value в map[string]*model.Value
func ValueMapToModel(metadata map[string]*repoModel.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ValueToModel(value)
	}
	return result
}
