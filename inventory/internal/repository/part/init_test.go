package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson/primitive"

	repoModel "github.com/nkolesnikov999/micro2-OK/inventory/internal/repository/model"
)

func (s *RepositorySuite) TestInitPartsSuccess() {
	// Очищаем коллекцию перед тестом
	_ = s.db.Collection("parts").Drop(s.ctx)

	// Создаем репозиторий без initParts
	collection := s.db.Collection("parts")
	r := &repository{collection: collection}

	// Инициализируем части
	err := r.initParts(s.ctx, 10)
	s.Require().NoError(err)

	// Проверяем, что части были добавлены
	count, err := collection.CountDocuments(s.ctx, primitive.M{})
	s.Require().NoError(err)
	s.Equal(int64(10), count)
}

func (s *RepositorySuite) TestInitPartsWithExistingData() {
	// Очищаем коллекцию перед тестом
	_ = s.db.Collection("parts").Drop(s.ctx)

	// Создаем репозиторий без initParts
	collection := s.db.Collection("parts")
	r := &repository{collection: collection}

	// Добавляем одну часть вручную
	_, err := collection.InsertOne(s.ctx, repoModel.Part{
		Uuid:          gofakeit.UUID(),
		Name:          "Test Part",
		Description:   "Test Description",
		Price:         100.0,
		StockQuantity: 10,
		Category:      repoModel.CategoryEngine,
	})
	s.Require().NoError(err)

	// Пытаемся инициализировать части (должно быть проигнорировано)
	err = r.initParts(s.ctx, 10)
	s.Require().NoError(err)

	// Проверяем, что количество частей не изменилось
	count, err := collection.CountDocuments(s.ctx, primitive.M{})
	s.Require().NoError(err)
	s.Equal(int64(1), count)
}

func (s *RepositorySuite) TestInitPartsWithZeroCount() {
	// Очищаем коллекцию перед тестом
	_ = s.db.Collection("parts").Drop(s.ctx)

	// Создаем репозиторий без initParts
	collection := s.db.Collection("parts")
	r := &repository{collection: collection}

	// Инициализируем с нулевым количеством
	err := r.initParts(s.ctx, 0)
	s.Require().NoError(err)

	// Проверяем, что части не были добавлены
	count, err := collection.CountDocuments(s.ctx, primitive.M{})
	s.Require().NoError(err)
	s.Equal(int64(0), count)
}

func (s *RepositorySuite) TestInitPartsWithNegativeCount() {
	// Очищаем коллекцию перед тестом
	_ = s.db.Collection("parts").Drop(s.ctx)

	// Создаем репозиторий без initParts
	collection := s.db.Collection("parts")
	r := &repository{collection: collection}

	// Инициализируем с отрицательным количеством
	// Теперь это должно работать без ошибки (возвращает пустой слайс)
	err := r.initParts(s.ctx, -5)
	s.Require().NoError(err)

	// Проверяем, что части не были добавлены
	count, err := collection.CountDocuments(s.ctx, primitive.M{})
	s.Require().NoError(err)
	s.Equal(int64(0), count)
}

func (s *RepositorySuite) TestInitPartsWithLargeCount() {
	// Очищаем коллекцию перед тестом
	_ = s.db.Collection("parts").Drop(s.ctx)

	// Создаем репозиторий без initParts
	collection := s.db.Collection("parts")
	r := &repository{collection: collection}

	// Инициализируем с большим количеством
	err := r.initParts(s.ctx, 1000)
	s.Require().NoError(err)

	// Проверяем, что все части были добавлены
	count, err := collection.CountDocuments(s.ctx, primitive.M{})
	s.Require().NoError(err)
	s.Equal(int64(1000), count)
}

func (s *RepositorySuite) TestCreateTestParts() {
	// Тестируем функцию createTestParts
	parts := createTestParts(5)

	s.Require().Len(parts, 5)

	for _, partInterface := range parts {
		part, ok := partInterface.(repoModel.Part)
		s.True(ok, "Part should be of type repoModel.Part")

		// Проверяем обязательные поля
		s.NotEmpty(part.Uuid)
		s.NotEmpty(part.Name)
		s.NotEmpty(part.Description)
		s.Greater(part.Price, 0.0)
		s.GreaterOrEqual(part.StockQuantity, int64(0))
		s.NotEqual(repoModel.CategoryUnspecified, part.Category)

		// Проверяем, что Dimensions не nil
		s.NotNil(part.Dimensions)
		s.Greater(part.Dimensions.Length, 0.0)
		s.Greater(part.Dimensions.Width, 0.0)
		s.Greater(part.Dimensions.Height, 0.0)
		s.Greater(part.Dimensions.Weight, 0.0)

		// Проверяем, что Manufacturer не nil
		s.NotNil(part.Manufacturer)
		s.NotEmpty(part.Manufacturer.Name)
		s.NotEmpty(part.Manufacturer.Country)
		s.NotEmpty(part.Manufacturer.Website)

		// Проверяем Tags
		s.NotNil(part.Tags)
		s.GreaterOrEqual(len(part.Tags), 1)
		s.LessOrEqual(len(part.Tags), 5)

		// Проверяем Metadata
		s.NotNil(part.Metadata)
		s.GreaterOrEqual(len(part.Metadata), 1)
		s.LessOrEqual(len(part.Metadata), 10)

		// Проверяем временные метки
		s.False(part.CreatedAt.IsZero())
		s.False(part.UpdatedAt.IsZero())
	}
}

func (s *RepositorySuite) TestFakeDimensions() {
	dimensions := fakeDimensions()

	s.NotNil(dimensions)
	s.GreaterOrEqual(dimensions.Length, 1.0)
	s.LessOrEqual(dimensions.Length, 300.0)
	s.GreaterOrEqual(dimensions.Width, 1.0)
	s.LessOrEqual(dimensions.Width, 300.0)
	s.GreaterOrEqual(dimensions.Height, 0.5)
	s.LessOrEqual(dimensions.Height, 150.0)
	s.GreaterOrEqual(dimensions.Weight, 0.1)
	s.LessOrEqual(dimensions.Weight, 500.0)
}

func (s *RepositorySuite) TestFakeManufacturer() {
	manufacturer := fakeManufacturer()

	s.NotNil(manufacturer)
	s.NotEmpty(manufacturer.Name)
	s.NotEmpty(manufacturer.Country)
	s.NotEmpty(manufacturer.Website)
}

func (s *RepositorySuite) TestFakeTags() {
	tags := fakeTags()

	s.NotNil(tags)
	s.GreaterOrEqual(len(tags), 1)
	s.LessOrEqual(len(tags), 5)

	for _, tag := range tags {
		s.NotEmpty(tag)
	}
}

func (s *RepositorySuite) TestRandomCategory() {
	// Тестируем несколько раз, чтобы убедиться в случайности
	categories := make(map[repoModel.Category]bool)

	for i := 0; i < 100; i++ {
		category := randomCategory()
		categories[category] = true

		// Проверяем, что категория не UNSPECIFIED
		s.NotEqual(repoModel.CategoryUnspecified, category)
	}

	// Проверяем, что получили разные категории
	s.Greater(len(categories), 1, "Should generate different categories")
}

func (s *RepositorySuite) TestFakeMetadata() {
	metadata := fakeMetadata()

	s.NotNil(metadata)
	s.GreaterOrEqual(len(metadata), 1)
	s.LessOrEqual(len(metadata), 10)

	for key, value := range metadata {
		s.NotEmpty(key)
		s.NotNil(value)

		// Проверяем, что только одно поле заполнено
		filledFields := 0
		if value.StringValue != "" {
			filledFields++
		}
		if value.Int64Value != 0 {
			filledFields++
		}
		if value.DoubleValue != 0 {
			filledFields++
		}
		// Для BoolValue проверяем, что оно было установлено
		// Если все остальные поля пустые, то должно быть BoolValue
		if value.StringValue == "" && value.Int64Value == 0 && value.DoubleValue == 0 {
			filledFields++ // BoolValue должно быть установлено
		}

		s.Equal(1, filledFields, "Exactly one field should be filled")
	}
}

func (s *RepositorySuite) TestFakeMetadataValue() {
	// Тестируем несколько раз, чтобы проверить разные типы
	stringCount := 0
	intCount := 0
	doubleCount := 0
	boolCount := 0

	for i := 0; i < 100; i++ {
		value := fakeMetadataValue()
		s.NotNil(value)

		// Проверяем каждый тип отдельно
		if value.StringValue != "" {
			stringCount++
		}
		if value.Int64Value != 0 {
			intCount++
		}
		if value.DoubleValue != 0 {
			doubleCount++
		}
		// Для BoolValue проверяем, что оно было установлено
		// Если все остальные поля пустые, то должно быть BoolValue
		if value.StringValue == "" && value.Int64Value == 0 && value.DoubleValue == 0 {
			boolCount++
		}
	}

	// Проверяем, что все типы генерируются
	s.Greater(stringCount, 0, "Should generate string values")
	s.Greater(intCount, 0, "Should generate int64 values")
	s.Greater(doubleCount, 0, "Should generate double values")
	s.Greater(boolCount, 0, "Should generate bool values")
}

func (s *RepositorySuite) TestInitPartsWithContextCancellation() {
	// Очищаем коллекцию перед тестом
	_ = s.db.Collection("parts").Drop(s.ctx)

	// Создаем отмененный контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	// Создаем репозиторий без initParts
	collection := s.db.Collection("parts")
	r := &repository{collection: collection}

	// Пытаемся инициализировать части с отмененным контекстом
	err := r.initParts(ctx, 10)
	s.Require().Error(err)
	s.Contains(err.Error(), "context canceled")
}

func (s *RepositorySuite) TestInitPartsDataIntegrity() {
	// Очищаем коллекцию перед тестом
	_ = s.db.Collection("parts").Drop(s.ctx)

	// Создаем репозиторий без initParts
	collection := s.db.Collection("parts")
	r := &repository{collection: collection}

	// Инициализируем части
	err := r.initParts(s.ctx, 50)
	s.Require().NoError(err)

	// Получаем все части из базы
	cursor, err := collection.Find(s.ctx, primitive.M{})
	s.Require().NoError(err)
	defer func() {
		if closeErr := cursor.Close(s.ctx); closeErr != nil {
			s.T().Logf("Failed to close cursor: %v", closeErr)
		}
	}()

	var parts []repoModel.Part
	err = cursor.All(s.ctx, &parts)
	s.Require().NoError(err)
	s.Len(parts, 50)

	// Проверяем уникальность UUID
	uuids := make(map[string]bool)
	for _, part := range parts {
		s.False(uuids[part.Uuid], "UUID should be unique: %s", part.Uuid)
		uuids[part.Uuid] = true

		// Проверяем корректность данных
		s.NotEmpty(part.Name)
		s.NotEmpty(part.Description)
		s.Greater(part.Price, 0.0)
		s.GreaterOrEqual(part.StockQuantity, int64(0))
		s.NotEqual(repoModel.CategoryUnspecified, part.Category)
	}
}
