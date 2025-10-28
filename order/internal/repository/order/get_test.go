package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *RepositorySuite) TestGetOrderSuccess() {
	// Создаем тестовый заказ
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.50,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Готовим список деталей
	parts := make([]model.Part, len(partUUIDs))
	for i, id := range partUUIDs {
		parts[i] = model.Part{Uuid: id}
	}
	// Создаем заказ в базе данных
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: testOrder.PartUuids}, parts)
	s.Require().NoError(err)

	// Получаем заказ
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)

	// Проверяем данные
	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestGetOrderNotFound() {
	// Пытаемся получить несуществующий заказ
	nonExistentUUID := uuid.New()

	result, err := s.repository.GetOrder(s.ctx, nonExistentUUID)
	s.Require().Error(err)
	s.Require().Equal(model.Order{}, result)
	s.Require().Equal(model.ErrOrderNotFound, err)
}

func (s *RepositorySuite) TestGetOrderWithContextCancellation() {
	// Создаем отмененный контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	// Пытаемся получить заказ с отмененным контекстом
	result, err := s.repository.GetOrder(ctx, uuid.New())
	s.Require().Error(err)
	s.Require().Equal(model.Order{}, result)
	// Проверяем, что ошибка связана с отменой контекста
	s.Require().Contains(err.Error(), "context canceled")
}

func (s *RepositorySuite) TestGetOrderWithPaidStatus() {
	// Создаем заказ со статусом PAID
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}
	transactionUUID := uuid.New()

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      250.75,
		TransactionUUID: transactionUUID.String(),
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	// Создаем заказ в базе данных
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: testOrder.PartUuids}, []model.Part{{Uuid: partUUIDs[0]}})
	s.Require().NoError(err)

	// Получаем заказ
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)

	// Проверяем данные
	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.TransactionUUID, result.TransactionUUID)
	s.Equal(testOrder.PaymentMethod, result.PaymentMethod)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestGetOrderWithCancelledStatus() {
	// Создаем заказ со статусом CANCELLED
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      500.00,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "CANCELLED",
	}

	// Создаем заказ в базе данных
	parts3 := make([]model.Part, len(partUUIDs))
	for i, id := range partUUIDs {
		parts3[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: testOrder.PartUuids}, parts3)
	s.Require().NoError(err)

	// Получаем заказ
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)

	// Проверяем данные
	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestGetOrderWithEmptyPartUUIDs() {
	// Создаем заказ с пустым списком частей
	orderUUID := uuid.New()
	userUUID := uuid.New()

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []uuid.UUID{}, // Пустой список
		TotalPrice:      0.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Создаем заказ в базе данных
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: testOrder.PartUuids}, []model.Part{})
	s.Require().NoError(err)

	// Получаем заказ
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)

	// Проверяем данные
	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestGetOrderWithManyPartUUIDs() {
	// Создаем заказ с большим количеством частей
	orderUUID := uuid.New()
	userUUID := uuid.New()

	// Создаем 10 частей
	partUUIDs := make([]uuid.UUID, 10)
	for i := 0; i < 10; i++ {
		partUUIDs[i] = uuid.New()
	}

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      1000.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Создаем заказ в базе данных
	parts10 := make([]model.Part, len(partUUIDs))
	for i, id := range partUUIDs {
		parts10[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, parts10)
	s.Require().NoError(err)

	// Получаем заказ
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)

	// Проверяем данные
	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}
