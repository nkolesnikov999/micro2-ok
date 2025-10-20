package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	inventoryV1 "github.com/nkolesnikov999/micro2-OK/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetNotFound() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(model.Part{}, model.ErrPartNotFound)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetServiceError() {
	var (
		serviceErr = gofakeit.Error()
		uuid       = gofakeit.UUID()

		req = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(model.Part{}, serviceErr)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
	s.Require().Contains(st.Message(), "internal error")
}

func (s *APISuite) TestGetSuccess() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
			Metadata:  map[string]*model.Value{"key": {StringValue: "value"}},
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Equal(part.Uuid, res.Part.Uuid)
	s.Require().Equal(part.Name, res.Part.Name)
	s.Require().Equal(part.Description, res.Part.Description)
	s.Require().Equal(part.Price, res.Part.Price)
	s.Require().Equal(part.StockQuantity, res.Part.StockQuantity)
	s.Require().Equal(part.Category, model.Category(res.Part.Category))
	s.Require().Equal(part.Tags, res.Part.Tags)
}

func (s *APISuite) TestGetInvalidUUID() {
	invalidUUIDs := []string{
		"invalid-uuid",
		"not-a-uuid",
		"123",
		"",
	}

	for _, invalidUUID := range invalidUUIDs {
		req := &inventoryV1.GetPartRequest{
			Uuid: invalidUUID,
		}

		res, err := s.api.GetPart(s.ctx, req)
		s.Require().Error(err)
		s.Require().Nil(res)

		st, ok := status.FromError(err)
		s.Require().True(ok)
		s.Require().Equal(codes.InvalidArgument, st.Code())
		s.Require().Contains(st.Message(), "invalid uuid format")
	}
}

func (s *APISuite) TestGetWithNilDimensions() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Nil(res.Part.Dimensions)
}

func (s *APISuite) TestGetWithNilManufacturer() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Nil(res.Part.Manufacturer)
}

func (s *APISuite) TestGetWithEmptyMetadata() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Empty(res.Part.Metadata)
}

func (s *APISuite) TestGetWithZeroStock() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Equal(int64(0), res.Part.StockQuantity)
}

func (s *APISuite) TestGetWithNegativePrice() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Equal(-100.0, res.Part.Price)
}

func (s *APISuite) TestGetWithUnspecifiedCategory() {
	var (
		uuid = gofakeit.UUID()
		req  = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Equal(model.CategoryUnspecified, model.Category(res.Part.Category))
}

func (s *APISuite) TestGetWithVeryLongName() {
	var (
		uuid         = gofakeit.UUID()
		veryLongName = gofakeit.Sentence() // very long name
		req          = &inventoryV1.GetPartRequest{
			Uuid: uuid,
		}
		part = model.Part{
			Uuid:          uuid,
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
		}
	)

	s.inventoryService.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.api.GetPart(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().NotNil(res.Part)
	s.Require().Equal(veryLongName, res.Part.Name)
}
