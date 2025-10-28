package part

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

func (r *repository) initParts(ctx context.Context, count int) error {
	// Проверяем, есть ли уже данные в коллекции
	countDocs, err := r.collection.CountDocuments(ctx, primitive.M{})
	if err != nil {
		return err
	}

	// Если данные уже есть, не инициализируем повторно
	if countDocs > 0 {
		return nil
	}

	parts := createTestParts(count)

	// Если нет частей для вставки, возвращаем nil
	if len(parts) == 0 {
		return nil
	}

	_, err = r.collection.InsertMany(ctx, parts)
	return err
}

func createTestParts(count int) []interface{} {
	// Если количество отрицательное или нулевое, возвращаем пустой слайс
	if count <= 0 {
		return []interface{}{}
	}

	parts := make([]interface{}, 0, count)
	now := time.Now()

	for range count {
		parts = append(parts, model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      randomCategory(),
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          fakeTags(),
			Metadata:      fakeMetadata(),
			CreatedAt:     now,
			UpdatedAt:     now,
		})
	}
	return parts
}

func fakeDimensions() *model.Dimensions {
	return &model.Dimensions{
		Length: gofakeit.Float64Range(1.0, 300.0),
		Width:  gofakeit.Float64Range(1.0, 300.0),
		Height: gofakeit.Float64Range(0.5, 150.0),
		Weight: gofakeit.Float64Range(0.1, 500.0),
	}
}

func fakeManufacturer() *model.Manufacturer {
	return &model.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
	}
}

func fakeTags() []string {
	tags := make([]string, 0, 5) // максимальная емкость
	for range gofakeit.IntRange(1, 5) {
		tags = append(tags, gofakeit.Word())
	}
	return tags
}

func randomCategory() model.Category {
	// Генерируем случайную категорию, исключая UNSPECIFIED (значение 0)
	vals := []model.Category{
		model.CategoryEngine,
		model.CategoryFuel,
		model.CategoryPorthole,
		model.CategoryWing,
	}
	return vals[gofakeit.IntRange(0, len(vals)-1)]
}

func fakeMetadata() map[string]*model.Value {
	metadata := make(map[string]*model.Value)

	for range gofakeit.IntRange(1, 10) {
		metadata[gofakeit.Word()] = fakeMetadataValue()
	}

	return metadata
}

func fakeMetadataValue() *model.Value {
	switch gofakeit.Number(0, 3) {
	case 0:
		return &model.Value{
			StringValue: gofakeit.Word(),
		}

	case 1:
		return &model.Value{
			Int64Value: int64(gofakeit.Number(1, 100)),
		}

	case 2:
		return &model.Value{
			DoubleValue: gofakeit.Float64Range(1, 100),
		}

	case 3:
		return &model.Value{
			BoolValue: gofakeit.Bool(),
		}

	default:
		return nil
	}
}
