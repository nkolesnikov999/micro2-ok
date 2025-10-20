package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestListSuccess() {
	var (
		req   = &inventoryV1.ListPartsRequest{}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word(), gofakeit.Word()},
				Metadata:  map[string]*model.Value{"key1": {StringValue: "value1"}},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word()},
				Metadata:  map[string]*model.Value{"key2": {Int64Value: 42}},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 2)
	s.Require().Equal(parts[0].Uuid, res.Parts[0].Uuid)
	s.Require().Equal(parts[0].Name, res.Parts[0].Name)
	s.Require().Equal(parts[1].Uuid, res.Parts[1].Uuid)
	s.Require().Equal(parts[1].Name, res.Parts[1].Name)
}

func (s *APISuite) TestListWithFilter() {
	var (
		filter = &inventoryV1.PartsFilter{
			Uuids:                 []string{gofakeit.UUID(), gofakeit.UUID()},
			Names:                 []string{"test1", "test2"},
			Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
			ManufacturerCountries: []string{"USA", "Germany"},
			Tags:                  []string{"tag1", "tag2"},
		}
		req            = &inventoryV1.ListPartsRequest{Filter: filter}
		expectedFilter = model.PartsFilter{
			Uuids:                 filter.Uuids,
			Names:                 filter.Names,
			Categories:            []model.Category{model.CategoryEngine},
			ManufacturerCountries: filter.ManufacturerCountries,
			Tags:                  filter.Tags,
		}
		parts = []model.Part{
			{
				Uuid:          filter.Uuids[0],
				Name:          "test1",
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Tags:          []string{"tag1"},
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, expectedFilter).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(parts[0].Uuid, res.Parts[0].Uuid)
	s.Require().Equal(parts[0].Name, res.Parts[0].Name)
}

func (s *APISuite) TestListEmptyResult() {
	req := &inventoryV1.ListPartsRequest{}

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return([]model.Part{}, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Empty(res.Parts)
}

func (s *APISuite) TestListServiceError() {
	var (
		serviceErr = gofakeit.Error()
		req        = &inventoryV1.ListPartsRequest{}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return([]model.Part{}, serviceErr)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)
	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
	s.Require().Contains(st.Message(), "internal error")
}

func (s *APISuite) TestListNotFound() {
	req := &inventoryV1.ListPartsRequest{}

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return([]model.Part{}, model.ErrPartNotFound)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
	s.Require().Contains(st.Message(), "parts not found")
}

func (s *APISuite) TestListWithNilFilter() {
	var (
		req   = &inventoryV1.ListPartsRequest{Filter: nil}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
}

func (s *APISuite) TestListWithUUIDFilter() {
	var (
		partUUIDs = []string{gofakeit.UUID(), gofakeit.UUID()}
		filter    = &inventoryV1.PartsFilter{
			Uuids: partUUIDs,
		}
		req   = &inventoryV1.ListPartsRequest{Filter: filter}
		parts = []model.Part{
			{
				Uuid:          partUUIDs[0],
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
			{
				Uuid:          partUUIDs[1],
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryWing,
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 partUUIDs,
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 2)
	s.Require().Equal(parts[0].Uuid, res.Parts[0].Uuid)
	s.Require().Equal(parts[1].Uuid, res.Parts[1].Uuid)
}

func (s *APISuite) TestListWithNameFilter() {
	var (
		names  = []string{"part1", "part2"}
		filter = &inventoryV1.PartsFilter{
			Names: names,
		}
		req   = &inventoryV1.ListPartsRequest{Filter: filter}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          names[0],
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 names,
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(parts[0].Name, res.Parts[0].Name)
}

func (s *APISuite) TestListWithCategoryFilter() {
	var (
		categories = []inventoryV1.Category{
			inventoryV1.Category_CATEGORY_ENGINE,
			inventoryV1.Category_CATEGORY_WING,
		}
		filter = &inventoryV1.PartsFilter{
			Categories: categories,
		}
		req            = &inventoryV1.ListPartsRequest{Filter: filter}
		expectedFilter = model.PartsFilter{
			Uuids:                 []string{},
			Names:                 []string{},
			Categories:            []model.Category{model.CategoryEngine, model.CategoryWing},
			ManufacturerCountries: []string{},
			Tags:                  []string{},
		}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, expectedFilter).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(parts[0].Category, model.Category(res.Parts[0].Category))
}

func (s *APISuite) TestListWithCountryFilter() {
	var (
		countries = []string{"USA", "Germany", "Japan"}
		filter    = &inventoryV1.PartsFilter{
			ManufacturerCountries: countries,
		}
		req   = &inventoryV1.ListPartsRequest{Filter: filter}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: "USA",
					Website: gofakeit.URL(),
				},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: countries,
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
}

func (s *APISuite) TestListWithTagFilter() {
	var (
		tags   = []string{"electronics", "premium", "new"}
		filter = &inventoryV1.PartsFilter{
			Tags: tags,
		}
		req   = &inventoryV1.ListPartsRequest{Filter: filter}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Tags:          []string{"electronics", "premium"},
				CreatedAt:     gofakeit.Date(),
				UpdatedAt:     gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  tags,
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Contains(res.Parts[0].Tags, "electronics")
	s.Require().Contains(res.Parts[0].Tags, "premium")
}

func (s *APISuite) TestListWithMultipleFilters() {
	var (
		filter = &inventoryV1.PartsFilter{
			Uuids:                 []string{gofakeit.UUID()},
			Names:                 []string{"test"},
			Categories:            []inventoryV1.Category{inventoryV1.Category_CATEGORY_ENGINE},
			ManufacturerCountries: []string{"USA"},
			Tags:                  []string{"premium"},
		}
		req            = &inventoryV1.ListPartsRequest{Filter: filter}
		expectedFilter = model.PartsFilter{
			Uuids:                 filter.Uuids,
			Names:                 filter.Names,
			Categories:            []model.Category{model.CategoryEngine},
			ManufacturerCountries: filter.ManufacturerCountries,
			Tags:                  filter.Tags,
		}
		parts = []model.Part{
			{
				Uuid:          filter.Uuids[0],
				Name:          "test",
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: "USA",
					Website: gofakeit.URL(),
				},
				Tags:      []string{"premium"},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, expectedFilter).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(parts[0].Uuid, res.Parts[0].Uuid)
	s.Require().Equal(parts[0].Name, res.Parts[0].Name)
	s.Require().Equal(parts[0].Category, model.Category(res.Parts[0].Category))
}

func (s *APISuite) TestListWithNilDimensions() {
	var (
		req   = &inventoryV1.ListPartsRequest{}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions:    nil, // nil dimensions
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word()},
				Metadata:  map[string]*model.Value{"key": {StringValue: "value"}},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Nil(res.Parts[0].Dimensions)
}

func (s *APISuite) TestListWithNilManufacturer() {
	var (
		req   = &inventoryV1.ListPartsRequest{}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: nil, // nil manufacturer
				Tags:         []string{gofakeit.Word()},
				Metadata:     map[string]*model.Value{"key": {StringValue: "value"}},
				CreatedAt:    gofakeit.Date(),
				UpdatedAt:    gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Nil(res.Parts[0].Manufacturer)
}

func (s *APISuite) TestListWithEmptyMetadata() {
	var (
		req   = &inventoryV1.ListPartsRequest{}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word()},
				Metadata:  make(map[string]*model.Value), // empty metadata
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Empty(res.Parts[0].Metadata)
}

func (s *APISuite) TestListWithZeroStock() {
	var (
		req   = &inventoryV1.ListPartsRequest{}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: 0, // zero stock
				Category:      model.CategoryEngine,
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word()},
				Metadata:  map[string]*model.Value{"key": {StringValue: "value"}},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(int64(0), res.Parts[0].StockQuantity)
}

func (s *APISuite) TestListWithNegativePrice() {
	var (
		req   = &inventoryV1.ListPartsRequest{}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         -100.0, // negative price
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word()},
				Metadata:  map[string]*model.Value{"key": {StringValue: "value"}},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(-100.0, res.Parts[0].Price)
}

func (s *APISuite) TestListWithUnspecifiedCategory() {
	var (
		req   = &inventoryV1.ListPartsRequest{}
		parts = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryUnspecified, // unspecified category
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word()},
				Metadata:  map[string]*model.Value{"key": {StringValue: "value"}},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(model.CategoryUnspecified, model.Category(res.Parts[0].Category))
}

func (s *APISuite) TestListWithVeryLongName() {
	var (
		veryLongName = gofakeit.Sentence() // very long name
		req          = &inventoryV1.ListPartsRequest{}
		parts        = []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          veryLongName,
				Description:   gofakeit.Sentence(),
				Price:         gofakeit.Price(100, 1000),
				StockQuantity: int64(gofakeit.IntRange(1, 100)),
				Category:      model.CategoryEngine,
				Dimensions: &model.Dimensions{
					Length: gofakeit.Float64Range(1, 100),
					Width:  gofakeit.Float64Range(1, 100),
					Height: gofakeit.Float64Range(1, 100),
					Weight: gofakeit.Float64Range(1, 100),
				},
				Manufacturer: &model.Manufacturer{
					Name:    gofakeit.Company(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      []string{gofakeit.Word()},
				Metadata:  map[string]*model.Value{"key": {StringValue: "value"}},
				CreatedAt: gofakeit.Date(),
				UpdatedAt: gofakeit.Date(),
			},
		}
	)

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, 1)
	s.Require().Equal(veryLongName, res.Parts[0].Name)
}

func (s *APISuite) TestListWithManyParts() {
	var (
		numParts = gofakeit.IntRange(10, 20)
		req      = &inventoryV1.ListPartsRequest{}
		parts    = make([]model.Part, numParts)
	)

	for i := 0; i < numParts; i++ {
		parts[i] = model.Part{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Sentence(),
			Price:         gofakeit.Price(100, 1000),
			StockQuantity: int64(gofakeit.IntRange(1, 100)),
			Category:      model.CategoryEngine,
			Dimensions: &model.Dimensions{
				Length: gofakeit.Float64Range(1, 100),
				Width:  gofakeit.Float64Range(1, 100),
				Height: gofakeit.Float64Range(1, 100),
				Weight: gofakeit.Float64Range(1, 100),
			},
			Manufacturer: &model.Manufacturer{
				Name:    gofakeit.Company(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			},
			Tags:      []string{gofakeit.Word()},
			Metadata:  map[string]*model.Value{"key": {StringValue: "value"}},
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	}

	s.inventoryService.On("ListParts", s.ctx, model.PartsFilter{
		Uuids:                 []string{},
		Names:                 []string{},
		Categories:            []model.Category{},
		ManufacturerCountries: []string{},
		Tags:                  []string{},
	}).Return(parts, nil)

	res, err := s.api.ListParts(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.Parts, numParts)
}
