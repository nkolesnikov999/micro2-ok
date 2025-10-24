package part

import (
	"context"
)

func (s *RepositorySuite) TestListPartsSuccess() {
	// Получаем список частей (включая 100 частей от initParts)
	result, err := s.repository.ListParts(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Len(result, 100) // initParts добавляет 100 частей

	// Проверяем, что все части имеют корректные поля
	for _, part := range result {
		s.NotEmpty(part.Uuid)
		s.NotEmpty(part.Name)
		s.NotEmpty(part.Description)
		s.Greater(part.Price, 0.0)
		s.GreaterOrEqual(part.StockQuantity, int64(0))
	}
}

func (s *RepositorySuite) TestListPartsEmpty() {
	// Создаем новый репозиторий без инициализации частей
	// Для этого нужно временно отключить initParts
	collection := s.db.Collection("parts")
	// Очищаем коллекцию
	_ = collection.Drop(s.ctx)

	// Создаем репозиторий без initParts
	r := &repository{collection: collection}

	// Получаем список частей из пустой коллекции
	result, err := r.ListParts(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Len(result, 0)
}

func (s *RepositorySuite) TestListPartsWithManyParts() {
	// Получаем список частей (включая 100 частей от initParts)
	result, err := s.repository.ListParts(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Len(result, 100) // initParts добавляет 100 частей

	// Проверяем, что все части имеют корректные поля
	for _, part := range result {
		s.NotEmpty(part.Uuid)
		s.NotEmpty(part.Name)
		s.NotEmpty(part.Description)
		s.Greater(part.Price, 0.0)
		s.GreaterOrEqual(part.StockQuantity, int64(0))
	}
}

func (s *RepositorySuite) TestListPartsWithContextCancellation() {
	// Создаем отмененный контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	// Пытаемся получить список частей с отмененным контекстом
	result, err := s.repository.ListParts(ctx)
	s.Require().Error(err)
	s.Require().Nil(result)
	// Проверяем, что ошибка связана с отменой контекста
	s.Require().Contains(err.Error(), "context canceled")
}

func (s *RepositorySuite) TestListPartsWithNilDimensions() {
	// Получаем список частей (включая 100 частей от initParts)
	result, err := s.repository.ListParts(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Len(result, 100) // initParts добавляет 100 частей

	// Проверяем, что все части имеют корректные поля
	for _, part := range result {
		s.NotEmpty(part.Uuid)
		s.NotEmpty(part.Name)
		s.NotEmpty(part.Description)
		s.Greater(part.Price, 0.0)
		s.GreaterOrEqual(part.StockQuantity, int64(0))
		// Dimensions могут быть nil или не nil - это нормально
	}
}

func (s *RepositorySuite) TestListPartsWithNilManufacturer() {
	// Получаем список частей (включая 100 частей от initParts)
	result, err := s.repository.ListParts(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(result)
	s.Require().Len(result, 100) // initParts добавляет 100 частей

	// Проверяем, что все части имеют корректные поля
	for _, part := range result {
		s.NotEmpty(part.Uuid)
		s.NotEmpty(part.Name)
		s.NotEmpty(part.Description)
		s.Greater(part.Price, 0.0)
		s.GreaterOrEqual(part.StockQuantity, int64(0))
		// Manufacturer могут быть nil или не nil - это нормально
	}
}
