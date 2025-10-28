package part

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/nkolesnikov999/micro2-OK/inventory/internal/model"
	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

func (s *RepositorySuite) TestGetPartSuccess() {
	testPart := repoModel.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(),
		Price:         gofakeit.Price(100, 1000),
		StockQuantity: int64(gofakeit.IntRange(1, 100)),
		Category:      repoModel.CategoryEngine,
		Dimensions: &repoModel.Dimensions{
			Length: gofakeit.Float64Range(1.0, 300.0),
			Width:  gofakeit.Float64Range(1.0, 300.0),
			Height: gofakeit.Float64Range(0.5, 150.0),
			Weight: gofakeit.Float64Range(0.1, 500.0),
		},
		Manufacturer: &repoModel.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.Word(), gofakeit.Word()},
		Metadata: map[string]*repoModel.Value{
			"test_key": {
				StringValue: "test_value",
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Collection("parts").InsertOne(s.ctx, testPart)
	s.Require().NoError(err)

	result, err := s.repository.GetPart(s.ctx, testPart.Uuid)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testPart.Uuid, result.Uuid)
	s.Equal(testPart.Name, result.Name)
	s.Equal(testPart.Description, result.Description)
	s.Equal(testPart.Price, result.Price)
	s.Equal(testPart.StockQuantity, result.StockQuantity)
	s.Equal(model.Category(testPart.Category), result.Category)
	s.Require().NotNil(result.Dimensions)
	s.Equal(testPart.Dimensions.Length, result.Dimensions.Length)
	s.Equal(testPart.Dimensions.Width, result.Dimensions.Width)
	s.Equal(testPart.Dimensions.Height, result.Dimensions.Height)
	s.Equal(testPart.Dimensions.Weight, result.Dimensions.Weight)
	s.Require().NotNil(result.Manufacturer)
	s.Equal(testPart.Manufacturer.Name, result.Manufacturer.Name)
	s.Equal(testPart.Manufacturer.Country, result.Manufacturer.Country)
	s.Equal(testPart.Manufacturer.Website, result.Manufacturer.Website)
	s.Equal(testPart.Tags, result.Tags)
	s.Require().NotNil(result.Metadata)
	s.Equal(testPart.Metadata["test_key"].StringValue, result.Metadata["test_key"].StringValue)
}

func (s *RepositorySuite) TestGetPartNotFound() {
	nonExistentUUID := gofakeit.UUID()

	result, err := s.repository.GetPart(s.ctx, nonExistentUUID)
	s.Require().Error(err)
	s.Require().Equal(model.ErrPartNotFound, err)
	s.Require().Equal(model.Part{}, result)
}

func (s *RepositorySuite) TestGetPartWithEmptyUUID() {
	result, err := s.repository.GetPart(s.ctx, "")
	s.Require().Error(err)
	s.Require().Equal(model.ErrPartNotFound, err)
	s.Require().Equal(model.Part{}, result)
}

func (s *RepositorySuite) TestGetPartWithInvalidUUID() {
	result, err := s.repository.GetPart(s.ctx, "invalid-uuid")
	s.Require().Error(err)
	s.Require().Equal(model.ErrPartNotFound, err)
	s.Require().Equal(model.Part{}, result)
}

func (s *RepositorySuite) TestGetPartWithContextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := s.repository.GetPart(ctx, gofakeit.UUID())
	s.Require().Error(err)
	s.Require().Equal(model.Part{}, result)
	s.Require().Contains(err.Error(), "context canceled")
}
