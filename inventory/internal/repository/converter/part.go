package converter

import (
	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
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
		Metadata:      ValueMapToRepoModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
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
		Metadata:      ValueMapToModel(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func CategoryToRepoModel(category model.Category) repoModel.Category {
	return repoModel.Category(category)
}

func CategoryToModel(category repoModel.Category) model.Category {
	return model.Category(category)
}

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
