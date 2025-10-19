package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      randomCategory(),
		Dimensions:    fakeDimensions(),
		Manufacturer:  fakeManufacturer(),
		Tags:          fakeTags(),
		Metadata:      fakeMetadata(),
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestGetPartError() {
	var (
		repoErr = gofakeit.Error()
		uuid    = gofakeit.UUID()
	)

	s.partRepository.On("GetPart", s.ctx, uuid).Return(model.Part{}, repoErr)

	res, err := s.service.GetPart(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}

func (s *ServiceSuite) TestGetPartWithEmptyUUID() {
	var (
		emptyUUID = ""
		repoErr   = gofakeit.Error()
	)

	s.partRepository.On("GetPart", s.ctx, emptyUUID).Return(model.Part{}, repoErr)

	res, err := s.service.GetPart(s.ctx, emptyUUID)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}

func (s *ServiceSuite) TestGetPartWithInvalidUUID() {
	var (
		invalidUUID = "invalid-uuid-format"
		repoErr     = gofakeit.Error()
	)

	s.partRepository.On("GetPart", s.ctx, invalidUUID).Return(model.Part{}, repoErr)

	res, err := s.service.GetPart(s.ctx, invalidUUID)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}

func (s *ServiceSuite) TestGetPartWithNilManufacturer() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      randomCategory(),
		Dimensions:    fakeDimensions(),
		Manufacturer:  nil, // nil manufacturer
		Tags:          fakeTags(),
		Metadata:      fakeMetadata(),
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Nil(res.Manufacturer)
}

func (s *ServiceSuite) TestGetPartWithNilDimensions() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      randomCategory(),
		Dimensions:    nil, // nil dimensions
		Manufacturer:  fakeManufacturer(),
		Tags:          fakeTags(),
		Metadata:      fakeMetadata(),
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Nil(res.Dimensions)
}

func (s *ServiceSuite) TestGetPartWithEmptyTags() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      randomCategory(),
		Dimensions:    fakeDimensions(),
		Manufacturer:  fakeManufacturer(),
		Tags:          []string{}, // empty tags
		Metadata:      fakeMetadata(),
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Empty(res.Tags)
}

func (s *ServiceSuite) TestGetPartWithEmptyMetadata() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      randomCategory(),
		Dimensions:    fakeDimensions(),
		Manufacturer:  fakeManufacturer(),
		Tags:          fakeTags(),
		Metadata:      make(map[string]*model.Value), // empty metadata
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Empty(res.Metadata)
}

func (s *ServiceSuite) TestGetPartWithZeroStock() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: 0, // zero stock
		Category:      randomCategory(),
		Dimensions:    fakeDimensions(),
		Manufacturer:  fakeManufacturer(),
		Tags:          fakeTags(),
		Metadata:      fakeMetadata(),
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Equal(int64(0), res.StockQuantity)
}

func (s *ServiceSuite) TestGetPartWithNegativePrice() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         -100.0, // negative price
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      randomCategory(),
		Dimensions:    fakeDimensions(),
		Manufacturer:  fakeManufacturer(),
		Tags:          fakeTags(),
		Metadata:      fakeMetadata(),
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Equal(-100.0, res.Price)
}

func (s *ServiceSuite) TestGetPartWithUnspecifiedCategory() {
	part := model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      model.CategoryUnspecified, // unspecified category
		Dimensions:    fakeDimensions(),
		Manufacturer:  fakeManufacturer(),
		Tags:          fakeTags(),
		Metadata:      fakeMetadata(),
		CreatedAt:     gofakeit.Date(),
		UpdatedAt:     gofakeit.Date(),
	}

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Equal(model.CategoryUnspecified, res.Category)
}

func (s *ServiceSuite) TestGetPartWithVeryLongName() {
	var (
		veryLongName = gofakeit.Sentence(1000) // very long name
		part         = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          veryLongName,
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      randomCategory(),
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          fakeTags(),
			Metadata:      fakeMetadata(),
			CreatedAt:     gofakeit.Date(),
			UpdatedAt:     gofakeit.Date(),
		}
	)

	s.partRepository.On("GetPart", s.ctx, part.Uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.Uuid)
	s.NoError(err)
	s.Equal(part, res)
	s.Equal(veryLongName, res.Name)
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
