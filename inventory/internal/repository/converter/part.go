package converter

import (
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

func ToRepoPart(part model.Part) repoModel.Part {
	return repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      ToRepoCategory(part.Category),
		Dimensions:    ToRepoDimensions(part.Dimensions),
		Manufacturer:  ToRepoManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ToRepoValueMap(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func ToModelPart(part repoModel.Part) model.Part {
	return model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      ToModelCategory(part.Category),
		Dimensions:    ToModelDimensions(part.Dimensions),
		Manufacturer:  ToModelManufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      ToModelValueMap(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func ToRepoCategory(category model.Category) repoModel.Category {
	return repoModel.Category(category)
}

func ToModelCategory(category repoModel.Category) model.Category {
	return model.Category(category)
}

func ToRepoDimensions(dimensions *model.Dimensions) *repoModel.Dimensions {
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

func ToModelDimensions(dimensions *repoModel.Dimensions) *model.Dimensions {
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

func ToRepoManufacturer(manufacturer *model.Manufacturer) *repoModel.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToModelManufacturer(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func ToRepoValue(value *model.Value) *repoModel.Value {
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

func ToModelValue(value *repoModel.Value) *model.Value {
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

func ToRepoValueMap(metadata map[string]*model.Value) map[string]*repoModel.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*repoModel.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ToRepoValue(value)
	}
	return result
}

func ToModelValueMap(metadata map[string]*repoModel.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}
	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		result[key] = ToModelValue(value)
	}
	return result
}
