package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

func (r *repository) InitParts(ctx context.Context, parts []model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(parts) == 0 {
		parts = CreateTestParts(100)
	}

	for _, part := range parts {
		r.parts[part.Uuid] = part
	}
	return nil
}

func CreateTestParts(count int) []model.Part {
	parts := make([]model.Part, 0, count)
	for range count {
		parts = append(parts, model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      randomCategory(),
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          fakeTags(),
			CreatedAt:     gofakeit.Date(),
			UpdatedAt:     gofakeit.Date(),
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
