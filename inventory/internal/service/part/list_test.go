package part

import (
	"github.com/brianvoe/gofakeit/v7"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsSuccess() {
	parts := []model.Part{
		{
			Uuid:          gofakeit.UUID(),
			Name:          "Engine Part 1",
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      model.CategoryEngine,
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          []string{"engine", "high-performance"},
			Metadata:      fakeMetadata(),
			CreatedAt:     gofakeit.Date(),
			UpdatedAt:     gofakeit.Date(),
		},
		{
			Uuid:          gofakeit.UUID(),
			Name:          "Wing Part 1",
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      model.CategoryWing,
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          []string{"wing", "aerodynamic"},
			Metadata:      fakeMetadata(),
			CreatedAt:     gofakeit.Date(),
			UpdatedAt:     gofakeit.Date(),
		},
	}

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, model.PartsFilter{})
	s.NoError(err)
	s.Equal(parts, res)
}

func (s *ServiceSuite) TestListPartsWithEmptyFilter() {
	parts := []model.Part{
		{
			Uuid:          gofakeit.UUID(),
			Name:          "Test Part",
			Description:   gofakeit.Sentence(10),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      model.CategoryEngine,
			Dimensions:    fakeDimensions(),
			Manufacturer:  fakeManufacturer(),
			Tags:          []string{"test"},
			Metadata:      fakeMetadata(),
			CreatedAt:     gofakeit.Date(),
			UpdatedAt:     gofakeit.Date(),
		},
	}

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	// Пустой фильтр должен возвращать все детали
	res, err := s.service.ListParts(s.ctx, model.PartsFilter{})
	s.NoError(err)
	s.Equal(parts, res)
}

func (s *ServiceSuite) TestListPartsWithUUIDFilter() {
	var (
		uuid1 = gofakeit.UUID()
		uuid2 = gofakeit.UUID()
		parts = []model.Part{
			{
				Uuid:          uuid1,
				Name:          "Part 1",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"test"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
			{
				Uuid:          uuid2,
				Name:          "Part 2",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"test"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
		filter = model.PartsFilter{
			Uuids: []string{uuid1},
		}
		expectedParts = []model.Part{parts[0]}
	)

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsWithNameFilter() {
	var (
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Engine Component",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"engine"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Wing Component",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"wing"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
		filter = model.PartsFilter{
			Names: []string{"Engine Component"},
		}
		expectedParts = []model.Part{parts[0]}
	)

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsWithCategoryFilter() {
	var (
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Engine Part",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"engine"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Wing Part",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"wing"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
		filter = model.PartsFilter{
			Categories: []model.Category{model.CategoryEngine},
		}
		expectedParts = []model.Part{parts[0]}
	)

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsWithCountryFilter() {
	var (
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "US Part",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    fakeDimensions(),
				Manufacturer: &model.Manufacturer{
					Name:    "US Company",
					Country: "USA",
					Website: "https://uscompany.com",
				},
				Tags:      []string{"us"},
				Metadata:  fakeMetadata(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          "German Part",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				Dimensions:    fakeDimensions(),
				Manufacturer: &model.Manufacturer{
					Name:    "German Company",
					Country: "Germany",
					Website: "https://germancompany.com",
				},
				Tags:      []string{"german"},
				Metadata:  fakeMetadata(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
		filter = model.PartsFilter{
			ManufacturerCountries: []string{"USA"},
		}
		expectedParts = []model.Part{parts[0]}
	)

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsWithTagFilter() {
	var (
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "High Performance Part",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"high-performance", "premium"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Standard Part",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"standard", "basic"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
		filter = model.PartsFilter{
			Tags: []string{"high-performance"},
		}
		expectedParts = []model.Part{parts[0]}
	)

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsWithMultipleFilters() {
	var (
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Engine Component",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    fakeDimensions(),
				Manufacturer: &model.Manufacturer{
					Name:    "US Company",
					Country: "USA",
					Website: "https://uscompany.com",
				},
				Tags:      []string{"high-performance", "premium"},
				Metadata:  fakeMetadata(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Wing Component",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				Dimensions:    fakeDimensions(),
				Manufacturer: &model.Manufacturer{
					Name:    "German Company",
					Country: "Germany",
					Website: "https://germancompany.com",
				},
				Tags:      []string{"standard", "basic"},
				Metadata:  fakeMetadata(),
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
		filter = model.PartsFilter{
			Categories:            []model.Category{model.CategoryEngine},
			ManufacturerCountries: []string{"USA"},
			Tags:                  []string{"high-performance"},
		}
		expectedParts = []model.Part{parts[0]}
	)

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsNoMatches() {
	var (
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Engine Part",
				Description:   gofakeit.Sentence(10),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    fakeDimensions(),
				Manufacturer:  fakeManufacturer(),
				Tags:          []string{"engine"},
				Metadata:      fakeMetadata(),
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
		filter = model.PartsFilter{
			Categories: []model.Category{model.CategoryWing}, // Ищем крылья, но есть только двигатели
		}
		expectedParts = []model.Part{} // Пустой результат
	)

	s.partRepository.On("ListParts", s.ctx).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestListPartsRepositoryError() {
	var (
		repoErr = gofakeit.Error()
		filter  = model.PartsFilter{}
	)

	s.partRepository.On("ListParts", s.ctx).Return(nil, repoErr)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Nil(res)
}
