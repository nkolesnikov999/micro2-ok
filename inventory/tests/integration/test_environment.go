//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"

	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/app"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/mongo"
	"github.com/nkolesnikov999/micro2-OK/platform/pkg/testcontainers/network"
)

const collectionParts = "parts"

// TestEnvironment — структура для хранения ресурсов тестового окружения
type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

// ... existing code ...
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	testPart := repoModel.Part{
		Uuid:          partUUID,
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

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionParts).InsertOne(ctx, testPart)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// ... existing code ...

// GetTestPart — создает тестовую запчасть, сохраняет в БД и возвращает её
func (env *TestEnvironment) GetTestPart(ctx context.Context) (repoModel.Part, error) {
	uuid := gofakeit.UUID()
	now := time.Now()

	part := repoModel.Part{
		Uuid:          uuid,
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(),
		Price:         gofakeit.Price(50, 5000),
		StockQuantity: int64(gofakeit.Number(0, 1000)),
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
			"test_key": {StringValue: "test_value"},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts"
	}

	if _, err := env.Mongo.Client().Database(databaseName).Collection(collectionParts).InsertOne(ctx, part); err != nil {
		return repoModel.Part{}, err
	}

	return part, nil
}

// ... existing code ...
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "parts" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(collectionParts).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

// ... existing code ...
