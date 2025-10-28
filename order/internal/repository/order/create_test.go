package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/nkolesnikov999/micro2-OK/order/internal/model"
)

func (s *RepositorySuite) TestCreateOrderSuccess() {
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

	// Создаем заказ в базе данных
	parts := make([]model.Part, len(partUUIDs))
	for i, id := range partUUIDs {
		parts[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, parts)
	s.Require().NoError(err)

	// Проверяем, что заказ действительно создан
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	// Проверяем данные
	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderAlreadyExists() {
	// Создаем тестовый заказ
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      50.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Создаем заказ первый раз
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, []model.Part{{Uuid: partUUIDs[0]}})
	s.Require().NoError(err)

	// Пытаемся создать заказ с тем же UUID второй раз
	err = s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, []model.Part{{Uuid: partUUIDs[0]}})
	s.Require().Error(err)
	s.Require().Equal(model.ErrOrderAlreadyExists, err)
}

func (s *RepositorySuite) TestCreateOrderWithContextCancellation() {
	// Создаем отмененный контекст
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Отменяем контекст сразу

	// Пытаемся создать заказ с отмененным контекстом
	testOrder := model.Order{
		OrderUUID:       uuid.New(),
		UserUUID:        uuid.New(),
		PartUuids:       []uuid.UUID{uuid.New()},
		TotalPrice:      100.0,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	parts := make([]model.Part, len(testOrder.PartUuids))
	for i, id := range testOrder.PartUuids {
		parts[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(ctx, testOrder, model.PartsFilter{Uuids: testOrder.PartUuids}, parts)
	s.Require().Error(err)
	// Проверяем, что ошибка связана с отменой контекста
	s.Require().Contains(err.Error(), "context canceled")
}

func (s *RepositorySuite) TestCreateOrderWithPaidStatus() {
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
	parts3 := make([]model.Part, len(partUUIDs))
	for i, id := range partUUIDs {
		parts3[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, parts3)
	s.Require().NoError(err)

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.TransactionUUID, result.TransactionUUID)
	s.Equal(testOrder.PaymentMethod, result.PaymentMethod)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderWithCancelledStatus() {
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
	parts10 := make([]model.Part, len(partUUIDs))
	for i, id := range partUUIDs {
		parts10[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, parts10)
	s.Require().NoError(err)

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderWithEmptyPartUUIDs() {
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

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderWithManyPartUUIDs() {
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
	parts := make([]model.Part, len(partUUIDs))
	for i, id := range partUUIDs {
		parts[i] = model.Part{Uuid: id}
	}
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, parts)
	s.Require().NoError(err)

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderWithZeroTotalPrice() {
	// Создаем заказ с нулевой ценой
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      0.0, // Нулевая цена
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Создаем заказ в базе данных
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, []model.Part{{Uuid: partUUIDs[0]}})
	s.Require().NoError(err)

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderWithNegativeTotalPrice() {
	// Создаем заказ с отрицательной ценой
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      -100.0, // Отрицательная цена
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Создаем заказ в базе данных
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, []model.Part{{Uuid: partUUIDs[0]}})
	s.Require().NoError(err)

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderWithVeryLargeTotalPrice() {
	// Создаем заказ с большой ценой (но в пределах DECIMAL(10,2))
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      99999999.99, // Максимальная цена для DECIMAL(10,2)
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          "PENDING_PAYMENT",
	}

	// Создаем заказ в базе данных
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, []model.Part{{Uuid: partUUIDs[0]}})
	s.Require().NoError(err)

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.Status, result.Status)
}

func (s *RepositorySuite) TestCreateOrderWithLongTransactionUUID() {
	// Создаем заказ с длинным transaction UUID (но в пределах VARCHAR(36))
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUIDs := []uuid.UUID{uuid.New()}

	// Создаем строку длиной 36 символов для transaction UUID
	longTransactionUUID := uuid.New().String() // UUID всегда 36 символов

	testOrder := model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      100.0,
		TransactionUUID: longTransactionUUID,
		PaymentMethod:   "CARD",
		Status:          "PAID",
	}

	// Создаем заказ в базе данных
	err := s.repository.CreateOrder(s.ctx, testOrder, model.PartsFilter{Uuids: partUUIDs}, []model.Part{{Uuid: partUUIDs[0]}})
	s.Require().NoError(err)

	// Проверяем, что заказ создан с правильными данными
	result, err := s.repository.GetOrder(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().NotNil(result)

	s.Equal(testOrder.OrderUUID, result.OrderUUID)
	s.Equal(testOrder.UserUUID, result.UserUUID)
	s.Equal(testOrder.PartUuids, result.PartUuids)
	s.Equal(testOrder.TotalPrice, result.TotalPrice)
	s.Equal(testOrder.TransactionUUID, result.TransactionUUID)
	s.Equal(testOrder.PaymentMethod, result.PaymentMethod)
	s.Equal(testOrder.Status, result.Status)
}
