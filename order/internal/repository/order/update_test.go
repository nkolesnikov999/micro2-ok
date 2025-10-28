package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *RepositorySuite) TestUpdateOrderSuccess() {
	// Создаем тестовый заказ
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.50,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Создаем заказ в базе данных
	parts := make([]model.Part, len(originalOrder.PartUuids))
	for i, id := range originalOrder.PartUuids {
		parts[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
	s.Require().NoError(err)

	// Обновляем заказ
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}, // Добавляем еще одну часть
		TotalPrice:      250.75,                                          // Увеличиваем цену
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	err = s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем, что заказ обновлен
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	// Проверяем обновленные данные
	s.Equal(updatedOrder.OrderUUID, result.OrderUUID)
	s.Equal(updatedOrder.UserUUID, result.UserUUID)
	s.Equal(updatedOrder.PartUuids, result.PartUuids)
	s.Equal(updatedOrder.TotalPrice, result.TotalPrice)
	s.Equal(updatedOrder.TransactionUUID, result.TransactionUUID)
	s.Equal(updatedOrder.PaymentMethod, result.PaymentMethod)
	s.Equal(updatedOrder.Status, result.Status)
}

func (s *RepositorySuite) TestUpdateOrderNotFound() {
	// Пытаемся обновить несуществующий заказ
	nonExistentUUID := uuid.New()

	updatedOrder := model.Order{
		OrderUUID:       nonExistentUUID,
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	err := s.repository.UpdateOrder(s.ctx, nonExistentUUID, updatedOrder)
	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderNotFound, err)
}

func (s *RepositorySuite) TestUpdateOrderWithContextCancellation() {
	// Создаем заказ для обновления
	orderUUID := uuid.New()
	userUUID := uuid.New()

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	parts := make([]model.Part, len(originalOrder.PartUuids))
	for i, id := range originalOrder.PartUuids {
		parts[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
	s.Require().NoError(err)

	// Создаем отмененный контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	// Пытаемся обновить заказ с отмененным контекстом
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      200.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PAID",
	}

	err = s.repository.UpdateOrder(ctx, orderUUID, updatedOrder)
	s.Require().Error(err)
	// Проверяем, что ошибка связана с отменой контекста
	s.Require().Contains(err.Error(), "context canceled")
}

func (s *RepositorySuite) TestUpdateOrderStatusToPaid() {
	// Создаем заказ со статусом PENDING_PAYMENT
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем статус на PAID
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: uuid.New().String(),
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленный статус
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal("PAID", result.Status)
	s.Equal(updatedOrder.TransactionUUID, result.TransactionUUID)
	s.Equal(updatedOrder.PaymentMethod, result.PaymentMethod)
}

func (s *RepositorySuite) TestUpdateOrderStatusToCancelled() {
	// Создаем заказ со статусом PENDING_PAYMENT
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем статус на CANCELLED
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленный статус
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal("CANCELLED", result.Status)
}

func (s *RepositorySuite) TestUpdateOrderWithEmptyPartUUIDs() {
	// Создаем заказ с частями
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем заказ с пустым списком частей
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []uuid.UUID{}, // Пустой список
		TotalPrice:      0.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленные данные
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal([]uuid.UUID{}, result.PartUuids)
	s.Equal(0.0, result.TotalPrice)
	s.Equal("CANCELLED", result.Status)
}

func (s *RepositorySuite) TestUpdateOrderWithManyPartUUIDs() {
	// Создаем заказ с одной частью
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем заказ с большим количеством частей
	updatedPartUUIDs := make([]uuid.UUID, 10)
	for i := 0; i < 10; i++ {
		updatedPartUUIDs[i] = uuid.New()
	}

	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       updatedPartUUIDs,
		TotalPrice:      1000.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленные данные
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(updatedPartUUIDs, result.PartUuids)
	s.Equal(1000.0, result.TotalPrice)
}

func (s *RepositorySuite) TestUpdateOrderWithNegativeTotalPrice() {
	// Создаем заказ с положительной ценой
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем заказ с отрицательной ценой
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      -50.0, // Отрицательная цена
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленные данные
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(-50.0, result.TotalPrice)
	s.Equal("CANCELLED", result.Status)
}

func (s *RepositorySuite) TestUpdateOrderWithZeroTotalPrice() {
	// Создаем заказ с положительной ценой
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем заказ с нулевой ценой
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      0.0, // Нулевая цена
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленные данные
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(0.0, result.TotalPrice)
	s.Equal("CANCELLED", result.Status)
}

func (s *RepositorySuite) TestUpdateOrderWithLargeTotalPrice() {
	// Создаем заказ с обычной ценой
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем заказ с большой ценой
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      99999999.99, // Большая цена (в пределах DECIMAL(10,2))
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленные данные
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(99999999.99, result.TotalPrice)
}

func (s *RepositorySuite) TestUpdateOrderWithDifferentUserUUID() {
	// Создаем заказ с одним пользователем
	orderUUID := uuid.New()
	originalUserUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        originalUserUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем заказ с другим пользователем
	newUserUUID := uuid.New()
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        newUserUUID, // Новый пользователь
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленные данные
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(newUserUUID, result.UserUUID)
}

func (s *RepositorySuite) TestUpdateOrderWithTransactionUUID() {
	// Создаем заказ без transaction UUID
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	originalOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	{
		parts := make([]model.Part, len(originalOrder.PartUuids))
		for i, id := range originalOrder.PartUuids {
			parts[i] = model.Part{Uuid: id}
		}
		err := s.repository.CreateOrder(s.ctx, originalOrder, model.PartsFilter{Uuids: originalOrder.PartUuids}, parts)
		s.Require().NoError(err)
	}

	// Обновляем заказ с transaction UUID
	transactionUUID := uuid.New().String()
	updatedOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: transactionUUID,
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	err := s.repository.UpdateOrder(s.ctx, orderUUID, updatedOrder)
	s.Require().NoError(err)

	// Проверяем обновленные данные
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Equal(transactionUUID, result.TransactionUUID)
	s.Equal("CARD", result.PaymentMethod)
	s.Equal("PAID", result.Status)
}
